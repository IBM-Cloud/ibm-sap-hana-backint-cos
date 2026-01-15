// Copyright 2026 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

// Functions for interacting with the IBM Cloud Object Storage
package cos

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/config"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Downloading one object
*/
func Download(s3Client *s3.S3, element CosObject) Result {
	global.Logger.Info(fmt.Sprintf(
		"Start downloading object '%s'.",
		element.Key),
	)

	sourceSize := getCosObjectSize(s3Client, element.Key)
	downloadParts, numParts := generateDownloadParts(
		s3Client,
		sourceSize,
		element.Key,
	)

	startTime := time.Now()

	// Opening destination pipe for writing
	fifo := openPipeForWriting(element.Destination)

	defer func() {
		_ = fifo.Close()
	}()

	// Getting the buffer size of the pipe
	pipeBufferSize := getPipeBufferSize(fifo)

	// Initializing asynchronous processing
	var wgGetObject sync.WaitGroup
	downloadPartsResults := make(chan DownloadPartResult, numParts)

	// Make sure that not more than the maximum number run concurrently
	sem := make(chan struct{}, config.BackintConfig.MaxConcurrency())

	// Map containing the parts which are already downloaded
	// but could not yet be written to pipe.
	// As the download is asynchronously,
	// the parts may not have the correct order.
	partsNotYetWritten := make(ByteMap)

	// Downloading parts asynchronously
	// (limited by value of maxConcurrency)
	for _, downloadPart := range downloadParts {
		wgGetObject.Add(1)
		sem <- struct{}{} // block if maxConcurrency reached

		downloadSingle := DownloadSingePart{
			fifo:               fifo,
			partsNotYetWritten: partsNotYetWritten,
			downloadPart:       downloadPart,
			nextIndex:          element.NextIndex,
			pipeBufferSize:     pipeBufferSize,
		}

		global.Logger.Debug(fmt.Sprintf("Next index for '%s' is '%d'", fifo.Name(), *element.NextIndex))

		go runDownloadSinglePart(
			s3Client,
			&wgGetObject,
			sem,
			downloadPartsResults,
			downloadSingle,
		)
	}

	// Waiting for all processes to finish
	go func(key string) {
		global.Logger.Debug(fmt.Sprintf(
			"'%s': Waiting for processes to finish.", key,
		))
		wgGetObject.Wait()
		close(downloadPartsResults)
	}(element.Key)

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()

	downloadedSize := int64(0)

	for r := range downloadPartsResults {
		if r.err != nil {
			global.Logger.Info(fmt.Sprintf("'%s': Error %s", element.Key, r.err))
			return Result{
				Err:        r.err,
				Duration:   duration,
				Key:        element.Key,
				ETag:       element.ETag,
				SourcePath: element.Destination,
			}
		}
		downloadedSize += r.size
	}

	global.Logger.Info(
		fmt.Sprintf("Finished downloading object '%s'.",
			element.Destination),
	)
	return Result{
		Err:        nil,
		Duration:   duration,
		Key:        element.Key,
		ETag:       element.ETag,
		SourcePath: element.Destination,
		SourceSize: sourceSize,
		TargetSize: downloadedSize,
	}
}

/*
Downloading one single part
*/
func runDownloadSinglePart(
	s3Client *s3.S3,
	wgGetObject *sync.WaitGroup,
	sem chan struct{},
	results chan DownloadPartResult,
	downloadSingle DownloadSingePart,
) {
	defer wgGetObject.Done()
	defer func() { <-sem }()

	global.Logger.Debug(
		fmt.Sprintf("Downloading part number '%d' of '%d' for key '%s'.",
			downloadSingle.downloadPart.partNumber,
			downloadSingle.downloadPart.numParts,
			downloadSingle.downloadPart.Key),
	)
	input := s3.GetObjectInput{
		Bucket:     aws.String(config.BackintConfig.BucketName()),
		Key:        aws.String(downloadSingle.downloadPart.Key),
		PartNumber: aws.Int64(downloadSingle.downloadPart.partNumber),
		Range:      aws.String(downloadSingle.downloadPart.byteRange),
	}

	response, err := s3Client.GetObject(&input)

	global.Logger.Debug(
		fmt.Sprintf("Finished downloading part number '%d' of '%d' for key '%s'.",
			downloadSingle.downloadPart.partNumber,
			downloadSingle.downloadPart.numParts,
			downloadSingle.downloadPart.Key,
		),
	)

	if err != nil {
		global.Logger.Error(
			fmt.Sprintf("'%s': Error downloading part with number '%d'.",
				downloadSingle.downloadPart.Key,
				downloadSingle.downloadPart.partNumber,
			))
		results <- DownloadPartResult{
			partNumber: downloadSingle.downloadPart.partNumber,
			err:        err,
		}
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)

	if err != nil {
		global.Logger.Error(fmt.Sprintf(
			"'%s': Could not copy data to buffer. Error: %s",
			downloadSingle.downloadPart.Key,
			err,
		))
		results <- DownloadPartResult{
			partNumber: downloadSingle.downloadPart.partNumber,
			err:        err,
		}
		return
	}

	if *(downloadSingle.nextIndex) > downloadSingle.downloadPart.numParts {
		global.Logger.Debug(fmt.Sprintf(
			"'%s': nextIndex higher than numParts.",
			downloadSingle.downloadPart.Key,
		))
		results <- DownloadPartResult{
			partNumber: downloadSingle.downloadPart.partNumber,
			err:        err,
			size:       *response.ContentLength,
		}
		return
	}

	// Writing data to buffer or pipe
	written := sendDataToHANA(
		downloadSingle.fifo,
		downloadSingle.nextIndex,
		&(downloadSingle.partsNotYetWritten),
		downloadSingle.downloadPart.partNumber,
		downloadSingle.pipeBufferSize,
		buf,
	)

	if !written {
		var errCouldNotWriteToPipe = errors.New("could not write to pipe")

		results <- DownloadPartResult{
			partNumber: downloadSingle.downloadPart.partNumber,
			err:        errCouldNotWriteToPipe,
			size:       *response.ContentLength,
		}
	} else {
		results <- DownloadPartResult{
			partNumber: downloadSingle.downloadPart.partNumber,
			err:        nil,
			size:       *response.ContentLength,
		}
	}
}

/*
Getting the numbers of parts uploaded of an object from IBM Cloud Object Storage
*/
func getPartsCount(s3Client *s3.S3, Key string) int64 {
	global.Logger.Debug(fmt.Sprintf(
		"Getting the PartsCount for key '%s'.", Key))
	result := getHeadObject(s3Client, Key)

	var partsCount int64 = 1
	if result.PartsCount != nil {
		partsCount = *result.PartsCount
	}

	global.Logger.Debug(fmt.Sprintf(
		"PartsCount for key '%s' is '%d'.",
		Key,
		partsCount),
	)
	return partsCount
}

/*
Getting the size of an object from from IBM Cloud Object Storage
*/
func getCosObjectSize(s3Client *s3.S3, Key string) int64 {
	global.Logger.Debug(fmt.Sprintf(
		"Getting the COS Object size for key '%s'.",
		Key),
	)
	result := getHeadObject(s3Client, Key)
	total_length := aws.Int64Value(result.ContentLength)
	global.Logger.Debug(
		fmt.Sprintf(
			"COS Object size for key '%s' is '%d'.",
			Key,
			total_length,
		),
	)
	return total_length
}

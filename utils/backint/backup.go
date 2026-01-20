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

package backint

import (
	"fmt"
	"sync"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/cos"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/logging"

	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Saving data in IBM Cloud Object Storage
*/
func Backup(
	s3Session *session.Session,
	s3Client *s3.S3,
) bool {
	global.Logger.Debug("Function: backup")
	sourcePaths := getSourcePathsForBackup()
	if len(sourcePaths) == 0 {
		global.Logger.Info(
			fmt.Sprintf("No source paths specified in %s", global.Args.InputFile),
		)
		return true
	}

	// Initializing asynchronous processing
	var wgUpload sync.WaitGroup
	chanUpload := make(chan cos.Result, len(sourcePaths))

	// Running all uploads asynchronously
	for x, sourcePath := range sourcePaths {
		wgUpload.Add(1)
		global.Logger.Info(fmt.Sprintf(
			"Storing '%s' in process #%d.", sourcePath, x,
		))
		go runUpload(s3Session, s3Client, &wgUpload, sourcePath, chanUpload)
	}

	// Waiting for all processes to finish
	wgUpload.Wait()
	close(chanUpload)

	global.Logger.Debug("All processes done.")

	// Checking the results
	return backupResultHandler(chanUpload)
}

/*
Executing the upload of one object to IBM Cloud Object Storage asynchronously
*/
func runUpload(
	s3Session *session.Session,
	s3Client *s3.S3,
	wg *sync.WaitGroup,
	pipe string,
	chanUpload chan cos.Result,
) {
	key := generateCosObjectKeyname(pipe)
	defer wg.Done()
	storeResult := cos.Upload(s3Session, s3Client, pipe, key)
	chanUpload <- storeResult
}

/*
Handling the results of uploading one object to COS
*/
func backupResultHandler(chanUpload chan cos.Result) bool {
	success := true
	for result := range chanUpload {
		if result.Err == nil {
			logging.BackintResultMsgs.AddBackupSuccessMessage(
				result.ETag,
				result.SourcePath,
				result.SourceSize,
			)
			logging.BackintResultMsgs.AddBackupMetrics(
				result.SourceSize,
				result.TargetSize,
				result.Duration,
			)
		} else {
			logging.BackintResultMsgs.AddErrorMessage(
				result.SourcePath,
				result.Err,
			)
			success = false
		}
	}
	return success
}

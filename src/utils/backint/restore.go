package backint

import (
	"fmt"
	"hdbbackint/utils/cos"
	"hdbbackint/utils/global"
	"hdbbackint/utils/logging"
	"sync"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Restoring the objects from IBM Cloud Object Storage
*/
func Restore(
	s3Client *s3.S3,
) bool {
	global.Logger.Debug("Function: restore")
	cosObjects := getCosObjectsForRestore()

	if cosObjects == nil {
		global.Logger.Error("Wrong keyword(s) in input file.")
		return false
	}

	// Initializing asynchronous processing
	var wgDownload sync.WaitGroup
	chanDownload := make(chan cos.Result, len(cosObjects))

	// Running all downloads asynchronously
	for n, element := range cosObjects {
		if element.ETag == "" {
			etag := cos.GetETagOfLatestVersionForKey(s3Client, element.Key)
			if etag == "" {
				chanDownload <- setObjectNotFoundResult(element)
				continue
			}
			element.ETag = etag
			wgDownload.Add(1)
		}
		logMessage := fmt.Sprintf(
			"Restoring backup '%s' with '%s' in process #%d",
			element.Key, element.ETag, n,
		)
		global.Logger.Info(logMessage)

		go runDownload(s3Client, &wgDownload, element, chanDownload)
	}
	wgDownload.Wait()
	close(chanDownload)

	global.Logger.Info("Restore: All processes finished.")

	// Checking the results of the single object downloads and return
	return restoreResultHandler(chanDownload)
}

/*
Executing download of a single object from
IBM Cloud Object Storage asynchronously
*/
func runDownload(
	s3Client *s3.S3,
	wg *sync.WaitGroup,
	element cos.CosObject,
	chanDownload chan cos.Result,
) {
	defer wg.Done()
	restoreResult := cos.Download(s3Client, element)
	chanDownload <- restoreResult
}

/*
Checking the results of all downloads, setting the messages and the return code
*/
func restoreResultHandler(chanDownload chan cos.Result) bool {
	success := true
	for result := range chanDownload {
		if result.Err == nil {
			if result.ETag == "" {
				// backup not found
				logging.BackintResultMsgs.AddObjectNotFoundMessage(
					result.SourcePath,
				)
			} else {
				// backup successful
				logging.BackintResultMsgs.AddRestoreSuccessMessage(
					result.ETag,
					result.SourcePath,
				)
			}
		} else {
			success = false
			logging.BackintResultMsgs.AddErrorMessage(
				result.SourcePath,
				result.Err,
			)
		}
	}
	return success
}

/*
Setting empty result struct
*/
func setObjectNotFoundResult(element cos.CosObject) cos.Result {
	return cos.Result{
		Err:        nil,
		Duration:   0,
		SourceSize: 0,
		TargetSize: 0,
		SourcePath: element.Destination,
		Key:        element.Key,
		ETag:       "",
	}
}

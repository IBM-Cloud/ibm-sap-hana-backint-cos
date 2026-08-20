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

	// Running all downloads asynchronously.
	// wgDownload is incremented before launching each goroutine so that
	// wgDownload.Wait() correctly blocks until every download finishes,
	// regardless of whether the ETag was already known or had to be looked up.
	for n, element := range cosObjects {
		wgDownload.Add(1)
		go func(n int, element cos.CosObject) {
			defer wgDownload.Done()

			// Resolve ETag and VersionId so Download can pin the exact COS
			// version.  On a versioned bucket, issuing GetObject without a
			// VersionId returns whatever the current latest version is — which
			// may be a delete marker or a newer object — causing HANA to see
			// "received 0 bytes" errors on PIT recovery.
			if element.ETag == "" {
				// NULL path: ETag not supplied by HANA — look up both.
				etag, versionId := cos.GetETagOfLatestVersionForKey(s3Client, element.Key)
				if etag == "" {
					chanDownload <- setObjectNotFoundResult(element)
					return
				}
				element.ETag = etag
				element.VersionId = versionId
			} else if element.VersionId == "" {
				// EBID path: ETag known but VersionId not yet resolved.
				// Resolve so the download is pinned to the correct version.
				element.VersionId = cos.GetVersionIdForETag(s3Client, element.Key, element.ETag)
			}

			global.Logger.Info(fmt.Sprintf(
				"Restoring backup '%s' with '%s' in process #%d",
				element.Key, element.ETag, n,
			))

			restoreResult := cos.Download(s3Client, element)
			chanDownload <- restoreResult
		}(n, element)
	}
	wgDownload.Wait()
	close(chanDownload)

	global.Logger.Info("Restore: All processes finished.")

	// Checking the results of the single object downloads and return
	return restoreResultHandler(chanDownload)
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

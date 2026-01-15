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

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/cos"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/logging"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Deleting the cloud objects specified in the input file
*/
func DeleteCloudObjects(
	s3Client *s3.S3,
) bool {
	success := true
	global.Logger.Debug("Function: delete")

	cosObjects := getCosObjectsForDelete(s3Client)
	deleteResults := cos.DeleteMultiple(s3Client, cosObjects)

	for _, r := range deleteResults {
		logging.BackintResultMsgs.AddKeyword(
			r.Status,
			[]string{r.ETag, r.Key},
		)
		if r.Status == "ERROR" {
			global.Logger.Error(
				fmt.Sprintf("Failed to delete object '%s' with ETag '%s'.",
					r.Key,
					r.ETag,
				),
			)
			success = false
		}
	}
	return success
}

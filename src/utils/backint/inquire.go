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
	"sort"
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/cos"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/logging"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Getting the objects from IBM Cloud Object Storage
*/
func Inquire(
	s3Client *s3.S3,
) bool {
	global.Logger.Debug("Function: inquire")

	for _, i := range global.InputFileContent {
		var splitted []string
		if strings.Contains(i.Parameter, " ") {
			splitted = strings.Split(i.Parameter, " ")
		}
		switch i.Keyword {
		case "NULL":
			Key := i.Parameter
			cosObjectList := cos.ListObjectsOfBucket(s3Client)
			sort.Slice(cosObjectList, func(i, j int) bool {
				return *cosObjectList[i].Key < *cosObjectList[j].Key
			})
			found := false
			for _, element := range cosObjectList {
				if Key != "" {
					if *element.Key == Key {
						found = true
						logging.BackintResultMsgs.AddKeyword(
							"BACKUP",
							[]string{*element.ETag, *element.Key},
						)
					}
				} else {
					found = true
					logging.BackintResultMsgs.AddKeyword(
						"BACKUP",
						[]string{*element.ETag},
					)
				}
			}

			if !found {
				// Nothing found
				if Key == "" {
					logging.BackintResultMsgs.AddKeyword(
						"NOTFOUND",
						nil,
					)
				} else {
					logging.BackintResultMsgs.AddKeyword(
						"NOTFOUND",
						[]string{Key},
					)
				}
			}

		case "EBID":
			if len(splitted) == 2 {
				ETag := splitted[0]
				Key := splitted[1]
				if cos.BackupExists(s3Client, ETag) {
					logging.BackintResultMsgs.AddKeyword(
						"BACKUP",
						[]string{ETag, Key},
					)
				} else {
					logging.BackintResultMsgs.AddKeyword(
						"NOTFOUND",
						[]string{ETag, Key},
					)
				}
			}
		default:
			// TODO Error -> Issue #21
		}
	}
	return true
}

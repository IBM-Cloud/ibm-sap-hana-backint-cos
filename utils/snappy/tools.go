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

package snappy

import (
	"fmt"
	"os"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/cos"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Executing the functions specified as arguments
*/
func Execute(function string) bool {
	_, s3Client := cos.GenerateCOSSession()

	switch function {
	case global.BUCKET_GET_LIFECYCLE:
		return getBucketLifeCycle(
			s3Client,
			global.Args.Bucket,
			global.Args.ResultFile,
		)
	case global.BUCKET_GET_LIST:
		return getObjectList(
			s3Client,
			global.Args.Bucket,
			global.Args.ResultFile,
		)
	case global.BUCKET_VERIFY:
		return verifyBucket(
			s3Client,
			global.Args.Bucket,
		)
	case global.FILE_UPLOAD:
		return uploadFile(
			s3Client,
			global.Args.Bucket,
			global.Args.Source,
			global.Args.Key,
		)
	}
	return true
}

/*
Verifying the given bucket
*/
func verifyBucket(s3Client *s3.S3, bucket string) bool {
	success, err := cos.RunBucketExists(s3Client, bucket)
	if err != nil {
		fmt.Printf("Error discovering bucket information: %s\n", err)
		return false
	}
	if success {
		// verify bucket versioning
		status, err := cos.RunIsBucketVersioning(s3Client, bucket)
		if err != nil {
			fmt.Printf("Error discovering bucket versioning: %s\n", err)
			return false
		}
		return status == "Enabled"
	}
	return success
}

/*
Getting the bucket lifecycle settings of the given bucket
*/
func getBucketLifeCycle(s3Client *s3.S3, bucket string, fileName string) bool {
	response, err := cos.RunGetBucketLifecycleRules(s3Client, bucket)

	if err != nil {
		fmt.Printf("Error discovering bucket lifecycle information: %s\n", err)
		return false
	}

	var lines []string

	for _, rule := range response {
		if *rule.Status == "Enabled" && rule.Expiration != nil {
			set, days := getExpirationDays(rule.Expiration)
			if set {
				lmRule := fmt.Sprintf(
					"ID:%s;Expiration:%d",
					*rule.ID,
					days,
				)
				lines = append(lines, lmRule)
			}
		}
	}
	return writeLinesToFile(fileName, lines)
}

/*
Listing all objects of a given bucket
*/
func getObjectList(s3Client *s3.S3, bucket string, fileName string) bool {
	response, err := cos.RunListObjectsOfBucket(s3Client, bucket)

	if err != nil {
		fmt.Printf("Error discovering bucket content: %s\n", err)
		return false
	}

	var lines []string
	for _, object := range response {
		lines = append(lines, *object.Key)
	}
	return writeLinesToFile(fileName, lines)
}

/*
Uploading one file to the given bucket
*/
func uploadFile(s3Client *s3.S3, bucket string, source string, key string) bool {
	err := cos.UploadSingleFile(s3Client, bucket, source, key)
	if err != nil {
		fmt.Printf("Error uploading file: %s\n", err)
	}
	return err == nil
}

/*
Getting the Expiration days of one rule
*/
func getExpirationDays(ruleExpiration *s3.LifecycleExpiration) (bool, int64) {
	if ruleExpiration.ExpiredObjectDeleteMarker != nil {
		if !*ruleExpiration.ExpiredObjectDeleteMarker {
			if *ruleExpiration.Days != 0 {
				return true, *ruleExpiration.Days
			} else {
				return false, 0
			}
		} else {
			return false, 0
		}
	}
	return true, *ruleExpiration.Days
}

func writeLinesToFile(fileName string, lines []string) bool {

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", fileName, err)
		return false
	}

	defer func() {
		_ = f.Close()
	}()

	for _, line := range lines {
		_, err = fmt.Fprintln(f, line)
		if err != nil {
			fmt.Printf("Error writing to file %s: %s\n", fileName, err)
			return false
		}
	}

	_ = f.Sync()

	return true
}

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

package cos

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/config"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
)

// ---------------------------------------------------------------------------
// Bucket checks
// ---------------------------------------------------------------------------

/*
Checking if the bucket exists and if versioning is enabled (if required)
*/
func BucketExists(s3Client *s3.S3) bool {
	bucket := config.BackintConfig.BucketName()
	global.Logger.Debug(fmt.Sprintf("Checking if bucket '%s' exists.", bucket))

	success, err := RunBucketExists(s3Client, bucket)

	if success {
		return true
	}

	if err != nil {
		global.Logger.Error(
			fmt.Sprintf("Error during getting bucket information. %s", err),
		)
		os.Exit(global.FAILURE)
	}
	return false
}

/*
Executing the bucket existence check
*/
func RunBucketExists(s3Client *s3.S3, bucket string) (bool, error) {
	_, err := s3Client.HeadBucket(
		&s3.HeadBucketInput{Bucket: aws.String(bucket)},
	)

	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "NotFound") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

/*
Checking if versioning is enabled for a given bucket
*/
func IsBucketVersioning(s3Client *s3.S3, bucket string) bool {
	global.Logger.Debug(
		fmt.Sprintf("Checking if versioning is set for '%s'.", bucket),
	)

	status, err := RunIsBucketVersioning(s3Client, bucket)

	global.CheckForError(err,
		fmt.Sprintf("Error discovering versioning of bucket '%s'", bucket),
		global.FAILURE,
	)

	global.Logger.Info(
		fmt.Sprintf("Versioning status of bucket '%s' is '%s'.",
			bucket,
			status,
		),
	)
	return status == "Enabled"
}

/*
Executing the call for bucket versioning
*/
func RunIsBucketVersioning(s3Client *s3.S3, bucket string) (string, error) {
	bucketVersioningInput := s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}
	bucketVersioning, err := s3Client.GetBucketVersioning(&bucketVersioningInput)
	if err != nil {
		return "", err
	}

	if bucketVersioning.Status != nil {
		return *bucketVersioning.Status, nil
	}
	return "", nil
}

// ---------------------------------------------------------------------------
// Object listing
// ---------------------------------------------------------------------------

/*
Getting the list of all objects for a given bucket
*/
func ListObjectsOfBucket(s3Client *s3.S3) []*s3.Object {
	bucket := config.BackintConfig.BucketName()

	global.Logger.Info(
		fmt.Sprintf("Creating list of all objects for bucket '%s'.", bucket),
	)

	cosObjectList, err := RunListObjectsOfBucket(s3Client, bucket)
	global.CheckForError(err,
		fmt.Sprintf(
			"Could not discover objects from bucket '%s'. Error: %s",
			bucket,
			err),
		global.FAILURE,
	)
	return cosObjectList
}

/*
Executing the discovery of the objects
*/
func RunListObjectsOfBucket(s3Client *s3.S3, bucket string) ([]*s3.Object, error) {
	isTruncated := true

	var cosObjectList []*s3.Object

	listObjectsInput := s3.ListObjectsInput{Bucket: aws.String(bucket)}
	for isTruncated {
		objectsOutput, err := s3Client.ListObjects(&listObjectsInput)

		if err != nil {
			return nil, err
		}

		isTruncated = *objectsOutput.IsTruncated
		cosObjectList = append(cosObjectList, objectsOutput.Contents...)
		listObjectsInput = s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Marker: objectsOutput.NextMarker}
	}
	return cosObjectList, nil
}

/*
Getting the list of lifecycle rules for bucket
*/
func RunGetBucketLifecycleRules(s3Client *s3.S3, bucket string) ([]*s3.LifecycleRule, error) {
	input := s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}
	response, err := s3Client.GetBucketLifecycleConfiguration(&input)
	if err != nil {
		return nil, err
	}

	return response.Rules, nil
}

// ---------------------------------------------------------------------------
// Object operations
// ---------------------------------------------------------------------------

/*
UploadSingleFile uploads a small file without multipart upload.
*/
func UploadSingleFile(s3Client *s3.S3, bucket string, source string, key string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	input := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   sourceFile,
	}
	_, err = s3Client.PutObject(&input)

	return err
}

/*
Upload uploads one object from a named pipe to IBM Cloud Object Storage
using the SDK's UploadWithPipe function.
*/
func Upload(
	s3Session *session.Session,
	s3Client *s3.S3,
	sourcePath string,
	key string,
) Result {
	global.Logger.Info(
		fmt.Sprintf("Uploading data from '%s' to '%s'.", sourcePath, key),
	)
	startTime := time.Now()
	global.Logger.Debug(fmt.Sprintf("multipart chunksize: %d", config.BackintConfig.MultipartChunksize()))

	uploader := s3manager.NewUploader(s3Session, func(u *s3manager.Uploader) {
		u.PartSize = config.BackintConfig.MultipartChunksize()
		u.Concurrency = config.BackintConfig.MaxConcurrency()
	})

	tags := config.BackintConfig.Tags()
	var pLockMode *string
	var pLockDate *time.Time
	var pLockLegalHold *string

	objectLockMode := config.BackintConfig.ObjectLockRetentionMode()
	if objectLockMode != global.LOCKMODE_NONE {
		lockMode := global.OBJECTLOCKMODES[objectLockMode]
		lockDate := config.BackintConfig.ObjectLockRetentionDate()
		pLockMode = &lockMode
		pLockDate = &lockDate
	}
	if lockLegalHold := config.BackintConfig.ObjectLockLegalHoldStatus(); lockLegalHold != "" {
		pLockLegalHold = &lockLegalHold
	}

	var pTags *string
	if tags != "" {
		pTags = &tags
	}

	uploadResult, copyError := uploader.UploadWithPipe(
		aws.BackgroundContext(),
		&s3manager.UploadPipeInput{
			PipePath: aws.String(sourcePath),
			UploadInput: &s3manager.UploadInput{
				Bucket:                    aws.String(config.BackintConfig.BucketName()),
				Key:                       aws.String(key),
				ObjectLockLegalHoldStatus: pLockLegalHold,
				ObjectLockMode:            pLockMode,
				ObjectLockRetainUntilDate: pLockDate,
				Tagging:                   pTags,
			},
		},
	)

	duration := time.Since(startTime).Seconds()

	if copyError != nil {
		global.Logger.Error(fmt.Sprintf(
			"Error uploading from %s. Error: %s",
			sourcePath,
			copyError),
		)
		return Result{
			Err:        copyError,
			SourcePath: sourcePath,
			Key:        key,
		}
	}

	global.Logger.Info(fmt.Sprintf(
		"Successfully uploaded '%s' to '%s'.",
		sourcePath,
		key),
	)

	size := getCosObjectSize(s3Client, key)
	return Result{
		Err:        nil,
		Duration:   duration,
		SourceSize: size,
		TargetSize: size,
		SourcePath: sourcePath,
		Key:        key,
		ETag:       *uploadResult.ETag,
	}
}

/*
Download downloads one object from IBM Cloud Object Storage to a named pipe
using the SDK's DownloadWithPipe function, which handles parallel part
downloads with ordered delivery.
*/
func Download(s3Client *s3.S3, element CosObject) Result {
	global.Logger.Info(fmt.Sprintf(
		"Start downloading object '%s'.",
		element.Key),
	)

	startTime := time.Now()

	downloader := s3manager.NewDownloaderWithClient(s3Client, func(d *s3manager.Downloader) {
		d.PartSize = config.BackintConfig.MultipartChunksize()
		d.Concurrency = config.BackintConfig.RecoverMaxConcurrency()
	})

	// PipeTimeout: 300 seconds matches the previous writeDataToPipe implementation.
	// HANA can slow its pipe reads significantly when processing data from other
	// channels or doing internal I/O — the default 30 s SDK timeout is too short
	// for large restores and causes "pipe write timeout" errors near completion.
	pipeTimeoutSecs := int64(300)

	n, err := downloader.DownloadWithPipe(
		aws.BackgroundContext(),
		&s3manager.DownloadPipeInput{
			PipePath:    aws.String(element.Destination),
			PipeTimeout: &pipeTimeoutSecs,
			GetObjectInput: &s3.GetObjectInput{
				Bucket: aws.String(config.BackintConfig.BucketName()),
				Key:    aws.String(element.Key),
			},
		},
	)

	duration := time.Since(startTime).Seconds()

	if err != nil {
		global.Logger.Error(fmt.Sprintf(
			"'%s': Error downloading object: %s",
			element.Key, err),
		)
		return Result{
			Err:        err,
			Duration:   duration,
			Key:        element.Key,
			ETag:       element.ETag,
			SourcePath: element.SourcePath,
		}
	}

	global.Logger.Info(fmt.Sprintf(
		"Finished downloading object '%s'.",
		element.Destination),
	)

	return Result{
		Err:        nil,
		Duration:   duration,
		Key:        element.Key,
		ETag:       element.ETag,
		SourcePath: element.SourcePath,
		// SourceSize == TargetSize == n: on a successful download every byte of
		// the COS object has been written to the pipe, so n equals ContentLength.
		// A separate HeadObject call to retrieve ContentLength independently
		// would be correct but wasteful given restoreResultHandler never reads
		// these fields.
		SourceSize: n,
		TargetSize: n,
	}
}

/*
Deleting multiple objects
*/
func DeleteMultiple(s3Client *s3.S3, cosObjects []CosObject) []CosObject {
	var results []CosObject
	for _, element := range cosObjects {
		if !element.Found {
			element.Status = "NOTFOUND"
			results = append(results, element)
			continue
		}

		deleteObjectInput := &s3.DeleteObjectInput{
			Bucket: aws.String(config.BackintConfig.BucketName()),
			Key:    aws.String(element.Key),
		}
		_, err := s3Client.DeleteObject(deleteObjectInput)
		element.Status = "DELETED"
		if err != nil {
			element.Status = "ERROR"
		}
		results = append(results, element)
	}
	return results
}

// ---------------------------------------------------------------------------
// Versioning
// ---------------------------------------------------------------------------

/*
Checking if a specific object exists
*/
func BackupExists(s3Client *s3.S3, ETag string) bool {
	cleanETag := strings.ReplaceAll(ETag, "\"", "")
	cosObjectList := ListObjectsOfBucket(s3Client)
	for _, element := range cosObjectList {
		if element.ETag != nil {
			cleanElementETag := strings.ReplaceAll(*element.ETag, "\"", "")
			if cleanElementETag == cleanETag {
				return true
			}
		}
	}
	return false
}

/*
Getting the ETag of the latest version of a given object
*/
func GetETagOfLatestVersionForKey(s3Client *s3.S3, Key string) string {
	global.Logger.Info(fmt.Sprintf("Getting latest version for '%s'.", Key))

	objectVersions := listObjectVersions(s3Client, Key, "")

	for _, v := range objectVersions {
		if *v.Key == Key && *v.IsLatest {
			ETag := strings.ReplaceAll(*v.ETag, "\"", "")
			global.Logger.Info(
				fmt.Sprintf("Latest version for key '%s' has entity tag '%s'.",
					Key,
					ETag,
				),
			)
			return ETag
		}
	}
	global.Logger.Info(fmt.Sprintf("No version found for key '%s'.", Key))
	return ""
}

/*
Getting the list of versions for a given object
*/
func listObjectVersions(
	s3Client *s3.S3,
	keyPrefix string,
	keyMarker string,
) []*s3.ObjectVersion {

	global.Logger.Info(
		fmt.Sprintf("Getting the list of object versions for key with prefix '%s'.",
			keyPrefix),
	)
	bucket := config.BackintConfig.BucketName()

	listObjectVersionsInput := s3.ListObjectVersionsInput{}
	if keyMarker != "" {
		listObjectVersionsInput = s3.ListObjectVersionsInput{
			Bucket:    aws.String(bucket),
			Prefix:    aws.String(keyPrefix),
			KeyMarker: aws.String(keyMarker)}
	} else {
		listObjectVersionsInput = s3.ListObjectVersionsInput{
			Bucket: aws.String(bucket),
			Prefix: aws.String(keyPrefix)}
	}

	listObjectVersionsOut, err := s3Client.ListObjectVersions(&listObjectVersionsInput)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Error discovering versions for '%s': %s", keyPrefix, err))
		return nil
	}

	if len(listObjectVersionsOut.Versions) == 0 {
		global.Logger.Error("No versions found")
	}
	var versions []*s3.ObjectVersion
	versions = append(versions, listObjectVersionsOut.Versions...)
	for _, v := range versions {
		global.Logger.Debug("Version :" + *v.VersionId)
	}
	return versions
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

/*
Getting the HeadObject for a given object
*/
func getHeadObject(s3Client *s3.S3, Key string) *s3.HeadObjectOutput {
	headObj := s3.HeadObjectInput{
		Bucket: aws.String(config.BackintConfig.BucketName()),
		Key:    aws.String(Key),
	}

	result, err := s3Client.HeadObject(&headObj)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Error getting HeadObject for Key '%s': %s", Key, err))
		return &s3.HeadObjectOutput{}
	}
	return result
}

// getCosObjectSize returns the total byte size of an object in IBM Cloud Object Storage.
func getCosObjectSize(s3Client *s3.S3, key string) int64 {
	global.Logger.Debug(fmt.Sprintf("Getting the COS Object size for key '%s'.", key))
	result := getHeadObject(s3Client, key)
	totalLength := aws.Int64Value(result.ContentLength)
	global.Logger.Debug(fmt.Sprintf("COS Object size for key '%s' is '%d'.", key, totalLength))
	return totalLength
}

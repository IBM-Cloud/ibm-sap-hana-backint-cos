// Generating the internal backint configuration from the arguments and the hdbbackint.cfg parameter file
package config

import (
	"fmt"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
)

/*
Getting the value for a given key from the backint configuration
*/
func (b BackintConfigT) Get(key string) string {
	return b[key]
}

/*
Setting a new value of a given key
*/
func (b BackintConfigT) set(key string, value string) {
	b[key] = value
}

/*
Getting the additional key prefix
used for setting the object name in IBM Cloud object storage
*/
func (b BackintConfigT) AdditionalKeyPrefix() string {
	return b.Get("additional_key_prefix")
}

/*
Getting the logging level in uppercase
*/
func (b BackintConfigT) AgentLogLevelU() string {
	return strings.ToUpper(b.Get("agent_log_level"))
}

/*
Getting the API Key
*/
func (b BackintConfigT) Apikey() string {
	return b.Get("apikey")
}

/*
Getting the path to the apikey file
*/
func (b BackintConfigT) authKeypath() string {
	return b.Get("auth_keypath")
}

/*
Getting the authentication method
*/
func (b BackintConfigT) AuthMethod() string {
	return b.Get("auth_mode")
}

/*
Getting the bucket name
*/
func (b BackintConfigT) BucketName() string {
	return b.Get("bucket")
}

/*
Checking if compression is set
*/
func (b BackintConfigT) Compression() bool {
	return strings.ToUpper(b.CompressionString()) == "TRUE"
}

func (b BackintConfigT) CompressionString() string {
	return b.Get("compression")
}

/*
Getting the compression level
*/
func (b BackintConfigT) CompressionLevel() zstd.EncoderLevel {
	zstdCompressionLevels := [4]zstd.EncoderLevel{
		zstd.SpeedFastest,
		zstd.SpeedDefault,
		zstd.SpeedBetterCompression,
		zstd.SpeedBestCompression}
	index := global.ToInteger(b.Get("compression_level")) - 1
	return zstdCompressionLevels[index]
}

/*
Getting the endpoint url
*/
func (b BackintConfigT) EndpointUrl() string {
	return b.Get("endpoint_url")
}

/*
Getting the IBM Authorization endpoint
*/
func (b BackintConfigT) IBMAuthEndpoint() string {
	return b.Get("ibm_auth_endpoint")
}

/*
Getting the maximum concurrency
*/
func (b BackintConfigT) MaxConcurrency() int {
	return global.ToInteger(b.Get("max_concurrency"))
}

/*
Getting the multipart chunksize
*/
func (b BackintConfigT) MultipartChunksize() int64 {
	return int64(global.ToInteger(b.Get("multipart_chunksize")))
}

/*
Getting the status of the object lock legal hold
*/
func (b BackintConfigT) ObjectLockLegalHoldStatus() string {
	return b.Get("object_lock_legal_hold_status")
}

/*
Getting the object lock retention mode
*/
func (b BackintConfigT) ObjectLockRetentionMode() string {
	return b.Get("object_lock_retention_mode")
}

/*
Getting the object lock retention date
*/
func (b BackintConfigT) ObjectLockRetentionDate() time.Time {
	rp := b.ObjectLockRetentionPeriod()
	splitted := strings.Split(rp, ",")
	y := global.ToInteger(splitted[0])
	m := global.ToInteger(splitted[1])
	d := global.ToInteger(splitted[2])

	retentionDate := time.Now().AddDate(y, m, d)

	fmt.Printf("ObjectLockRetentionDate set to %s\n",
		retentionDate.String(),
	)
	return retentionDate
}

/*
Getting the object lock retention period
*/
func (b BackintConfigT) ObjectLockRetentionPeriod() string {
	return b.Get("object_lock_retention_period")
}

/*
Getting the region
*/
func (b BackintConfigT) Region() string {
	return b.Get("region")
}

/*
Getting the key prefix to be removed
*/
func (b BackintConfigT) RemoveKeyPrefix() string {
	return b.Get("remove_key_prefix")
}

/*
Getting the service Instance Id
*/
func (b BackintConfigT) ServiceInstanceId() string {
	return b.Get("service_instance_id")
}

/*
Getting the tags
*/
func (b BackintConfigT) Tags() string {
	tags := b.Get("object_tags")
	if tags != "" {
		tags = strings.ReplaceAll(tags, ",", "&")
	}
	return tags
}

/*
Getting the timeout
*/
func (b BackintConfigT) Timeout() int {
	return global.ToInteger(b.Get("timeout_microsecond"))
}

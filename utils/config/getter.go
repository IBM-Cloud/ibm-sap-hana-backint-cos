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

package config

import (
	"fmt"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"strings"
	"time"
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
	return b.Get(additional_key_prefix.key)
}

/*
Getting the logging level in uppercase
*/
func (b BackintConfigT) AgentLogLevelU() string {
	return strings.ToUpper(b.Get(agent_log_level.key))
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
func (b BackintConfigT) AuthKeypath() string {
	return b.Get(auth_keypath.key)
}

/*
Getting the authentication mode
*/
func (b BackintConfigT) AuthMode() string {
	return b.Get(auth_mode.key)
}

/*
Getting the bucket name
*/
func (b BackintConfigT) BucketName() string {
	return b.Get(bucket.key)
}

/*
Getting the endpoint url
*/
func (b BackintConfigT) EndpointUrl() string {
	return b.Get(endpoint_url.key)
}

/*
Getting the IAM Profile Id
*/
func (b BackintConfigT) IamProfileId() string {
	return b.Get(iam_profile_id.key)
}

/*
Getting the IAM Profile Name
*/
func (b BackintConfigT) IamProfileName() string {
	return b.Get(iam_profile_name.key)
}

/*
Getting the IBM Authorization endpoint
*/
func (b BackintConfigT) IBMAuthEndpoint() string {
	return b.Get(ibm_auth_endpoint.key)
}

/*
Getting the maximum concurrency
*/
func (b BackintConfigT) MaxConcurrency() int {
	return global.ToInteger(b.Get(max_concurrency.key))
}

/*
Getting the multipart chunksize
*/
func (b BackintConfigT) MultipartChunksize() int64 {
	return int64(global.ToInteger(b.Get(multipart_chunksize.key)))
}

/*
Getting the status of the object lock legal hold
*/
func (b BackintConfigT) ObjectLockLegalHoldStatus() string {
	return b.Get(object_lock_legal_hold_status.key)
}

/*
Getting the object lock retention mode
*/
func (b BackintConfigT) ObjectLockRetentionMode() string {
	return strings.ToLower(b.Get(object_lock_retention_mode.key))
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
	return b.Get(object_lock_retention_period.key)
}

/*
Getting the max concurrency for recovery
If set to 0, we take the value from max_concurrency
*/
func (b BackintConfigT) RecoverMaxConcurrency() int {
	return global.ToInteger(b.Get(recover_max_concurrency.key))
}

/*
Getting the region
*/
func (b BackintConfigT) Region() string {
	return b.Get(region.key)
}

/*
Getting the key prefix to be removed
*/
func (b BackintConfigT) RemoveKeyPrefix() string {
	return b.Get(remove_key_prefix.key)
}

/*
Getting the tags
*/
func (b BackintConfigT) Tags() string {
	tags := b.Get(object_tags.key)
	if tags != "" {
		tags = strings.ReplaceAll(tags, ",", "&")
	}
	return tags
}

/*
Getting the timeout
*/
func (b BackintConfigT) Timeout() int {
	return global.ToInteger(b.Get(timeout_microsecond.key))
}

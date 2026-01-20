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

/*
cloud_storage Section
*/
var auth_keypath = Default{
	key:            "auth_keypath",
	section:        SECTION_CLOUD_STORAGE,
	mandatory:      true,
	validationType: CONFIG_FILE}

var auth_mode = Default{
	key:            "auth_mode",
	section:        SECTION_CLOUD_STORAGE,
	mandatory:      true,
	defaultValue:   AUTH_APIKEY,
	possibleValues: []string{AUTH_APIKEY},
	validationType: CONFIG_LIST}

var bucket = Default{
	key:            "bucket",
	section:        SECTION_CLOUD_STORAGE,
	mandatory:      true,
	validationType: CONFIG_STRING}

var endpoint_url = Default{
	key:            "endpoint_url",
	section:        SECTION_CLOUD_STORAGE,
	mandatory:      true,
	validationType: CONFIG_URL}

var ibm_auth_endpoint = Default{
	key:            "ibm_auth_endpoint",
	section:        SECTION_CLOUD_STORAGE,
	defaultValue:   "https://private.iam.cloud.ibm.com/identity/token",
	mandatory:      false,
	validationType: CONFIG_URL}

var region = Default{
	key:     "region",
	section: SECTION_CLOUD_STORAGE,
	possibleValues: []string{
		"au-syd",
		"br-sao",
		"ca-tor",
		"eu-de",
		"eu-es",
		"eu-gb",
		"jp-osa",
		"jp-tok",
		"us-east",
		"us-south"},
	mandatory:      true,
	validationType: CONFIG_LIST}

/*
backint Section
*/
var max_concurrency = Default{
	key:            "max_concurrency",
	section:        SECTION_BACKINT,
	defaultValue:   "10",
	min:            1,
	max:            20,
	mandatory:      false,
	validationType: CONFIG_RANGE}

var multipart_chunksize = Default{
	key:            "multipart_chunksize",
	section:        SECTION_BACKINT,
	defaultValue:   "134000000",
	mandatory:      false,
	validationType: CONFIG_CHUNKSIZE}

// Not propagated to customer
var timeout_microsecond = Default{
	key:            "timeout_microsecond",
	section:        SECTION_BACKINT,
	defaultValue:   "1",
	mandatory:      false,
	validationType: CONFIG_INT,
}

/*
object Section
*/
var additional_key_prefix = Default{
	key:            "additional_key_prefix",
	section:        SECTION_OBJECTS,
	defaultValue:   "",
	mandatory:      false,
	validationType: CONFIG_STRING}

var remove_key_prefix = Default{
	key:            "remove_key_prefix",
	section:        SECTION_OBJECTS,
	defaultValue:   "",
	mandatory:      false,
	validationType: CONFIG_STRING}

var object_lock_legal_hold_status = Default{
	key:            "object_lock_legal_hold_status",
	section:        SECTION_OBJECTS,
	defaultValue:   "OFF",
	possibleValues: []string{"OFF", "ON"},
	mandatory:      false,
	validationType: CONFIG_STRING}

var object_lock_retention_mode = Default{
	key:            "object_lock_retention_mode",
	section:        SECTION_OBJECTS,
	defaultValue:   "None",
	possibleValues: []string{"None", "cmp"},
	mandatory:      false,
	validationType: CONFIG_LIST}

var object_lock_retention_period = Default{
	key:            "object_lock_retention_period",
	section:        SECTION_OBJECTS,
	defaultValue:   "0,0,0",
	mandatory:      false,
	validationType: CONFIG_PERIOD}

var object_tags = Default{
	key:            "object_tags",
	section:        SECTION_OBJECTS,
	defaultValue:   "",
	mandatory:      false,
	validationType: CONFIG_TAG}

var agent_log_level = Default{
	key:          "agent_log_level",
	section:      SECTION_TRACE,
	defaultValue: "info",
	possibleValues: []string{
		"debug",
		"info",
		"warning",
		"error",
		"critical",
		"http"},
	mandatory:      false,
	validationType: CONFIG_LIST}

var configDefaults = []Default{
	auth_mode,
	auth_keypath,
	bucket,
	region,
	endpoint_url,
	ibm_auth_endpoint,
	max_concurrency,
	multipart_chunksize,
	remove_key_prefix,
	additional_key_prefix,
	object_tags,
	object_lock_retention_mode,
	object_lock_retention_period,
	object_lock_legal_hold_status,
	agent_log_level,
	timeout_microsecond,
}

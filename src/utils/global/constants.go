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

package global

// Versions
const (
	BACKINT_VERSION = "backint 1.04"
	TOOL_VERSION    = "Backint for IBM Object Store version: '0.0.3'"
)

// Function names
const (
	BACKUP        = "BACKUP"
	DELETE        = "DELETE"
	INQUIRE       = "INQUIRE"
	RESTORE       = "RESTORE"
	INTERNAL_TEST = "TEST"

	// Functions used for calls from dbbackup tool
	BUCKET_VERIFY        = "BUCKET-VERIFY"
	BUCKET_GET_LIST      = "BUCKET-GET-LIST"
	BUCKET_GET_LIFECYCLE = "BUCKET-GET-LIFECYCLE"
	FILE_UPLOAD          = "FILE-UPLOAD"
)

var FUNCTIONLIST = []string{
	BACKUP,
	DELETE,
	INQUIRE,
	RESTORE,
	INTERNAL_TEST,
	BUCKET_VERIFY,
	BUCKET_GET_LIST,
	BUCKET_GET_LIFECYCLE,
	FILE_UPLOAD,
}

var DBBACKUP_FUNCTIONLIST = []string{
	BUCKET_VERIFY,
	BUCKET_GET_LIST,
	BUCKET_GET_LIFECYCLE,
	FILE_UPLOAD,
}

// Exit codes
const (
	SUCCESS         = 0
	FAILURE         = 1
	WRONG_PARAMETER = 2
)

// Values for METADATA
const METADATA_COMPRESSION_LABEL string = "compression"

// Values for comparison with parameter file settings
const OBJECTLOCKMODE string = "COMPLIANCE"

// Default pipe buffer size used for recovery in case
// the system call to get the buffer size produces an error
const PIPE_BUFFER_SIZE = 1024 * 1024 * 1024

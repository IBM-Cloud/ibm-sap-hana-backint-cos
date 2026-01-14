/*
Contains the constants
*/
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

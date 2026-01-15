//Constants for configuration parameter validation

package config

const (
	CONFIG_BOOL      = "bool"
	CONFIG_CHUNKSIZE = "chunksize"
	CONFIG_FILE      = "file"
	CONFIG_INT       = "int"
	CONFIG_LIST      = "list"
	CONFIG_PERIOD    = "period"
	CONFIG_RANGE     = "range"
	CONFIG_STRING    = "string"
	CONFIG_TAG       = "tag"
	CONFIG_URL       = "url"
)

const (
	UNIT_KB = "KB"
	UNIT_MB = "MB"
	UNIT_GB = "GB"
)

var validConfigTypes = []string{
	CONFIG_BOOL,
	CONFIG_CHUNKSIZE,
	CONFIG_FILE,
	CONFIG_INT,
	CONFIG_LIST,
	CONFIG_PERIOD,
	CONFIG_RANGE,
	CONFIG_STRING,
	CONFIG_TAG,
	CONFIG_URL,
}

var validSizeUnits = []string{
	UNIT_KB,
	UNIT_MB,
	UNIT_GB,
}

// Parameter configuration file sections
const (
	SECTION_CLOUD_STORAGE = "cloud_storage"
	SECTION_BACKINT       = "backint"
	SECTION_OBJECTS       = "objects"
	SECTION_TRACE         = "trace"
)

var validSections = []string{
	SECTION_CLOUD_STORAGE,
	SECTION_BACKINT,
	SECTION_OBJECTS,
	SECTION_TRACE,
}

// Maximum number of allowed tags
const MAX_NUMBER_OF_TAGS int = 10

// Modes for authentication method
const (
	AUTH_APIKEY string = "apikey"
)

// File validation values
const (
	FILEOK          = 0
	FILENOTFOUND    = 1
	FILENOTREADABLE = 2

	FILEMUSTEXIST = true
	FILENOTEXIST  = false
)

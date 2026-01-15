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

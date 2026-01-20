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
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/version"
)

/*
Reading and validating the command line arguments
*/
func GetCommandLineArguments() (global.CommandLineArguments, bool, bool) {
	global.Args = parseArguments()

	return global.Args,
		argsValid(),
		isDbBackupFunction(global.Args.Function)
}

/*
Parsing the command line arguments
*/
func parseArguments() global.CommandLineArguments {
	parameterFile := flag.String("p", "", "Parameter file")
	userId := flag.String("u", "", "user id")
	function := flag.String("f", "", "Function to be executed")
	inputFile := flag.String("i", "", "Input file")
	outputFile := flag.String("o", "", "Output file")
	backupId := flag.Int("s", -1, "Backup Id")
	numberOfObjects := flag.Int("c", -1, "Number of objects")
	backupLevel := flag.String("l", "", "backup level")
	authKeypath := flag.String("keypath", "", "path to the apikey file")
	authEndpoint := flag.String("authendpoint", "https://private.iam.cloud.ibm.com/identity/token", "IBM auth endpoint")
	region := flag.String("region", "", "region")
	endpointUrl := flag.String("endpoint", "", "endpoint url")
	bucket := flag.String("bucket", "", "bucket name")
	source := flag.String("source", "", "source file path")
	key := flag.String("key", "", "object name")
	resultFile := flag.String("r", "", "file containing values")

	// Version
	var version bool
	flag.BoolVar(&version, "V", false, "Print version string")
	flag.BoolVar(&version, "v", false, "Print version string")

	// Check parameter file settings
	var checkParms bool
	flag.BoolVar(&checkParms, "check", false, "check parameter file")

	flag.Parse()

	global.Args.ParameterFile = *parameterFile
	global.Args.UserId = *userId
	global.Args.Function = strings.ToUpper(*function)
	global.Args.InputFile = *inputFile
	global.Args.OutputFile = *outputFile
	global.Args.BackupId = *backupId
	global.Args.NumberOfObjects = *numberOfObjects
	global.Args.BackupLevel = *backupLevel
	global.Args.Version = version
	global.Args.CheckParms = checkParms

	// Used when called from snappy agent
	global.Args.AuthKeypath = *authKeypath
	global.Args.AuthEndpoint = *authEndpoint
	global.Args.EndpointUrl = *endpointUrl
	global.Args.Region = *region
	global.Args.Bucket = *bucket
	global.Args.Source = *source
	global.Args.Key = *key
	global.Args.ResultFile = *resultFile

	return global.Args
}

/*
Printing out the version
*/
func PrintVersion() {
	fmt.Printf("\"%s\" \"%s\" \n",
		version.BACKINT_VERSION,
		version.TOOL_VERSION,
	)
}

/*
Validating command line arguments
*/
func argsValid() bool {
	// check version flag
	if global.Args.Version {
		return true
	}

	// If --check specified, the -p must be specified too
	if global.Args.CheckParms {
		if global.Args.ParameterFile != "" {
			message := isFileValid(global.Args.ParameterFile, FILEMUSTEXIST)
			if message != "" {
				fmt.Println("Parameter", message)
				return false
			}
		} else {
			fmt.Println(
				"You specified --check but the parameter file option is missing.",
			)
			return false
		}
		return true
	}

	// Check function
	if !isFunctionValid(global.Args.Function) {
		fmt.Println("Invalid function specified.")
		return false
	}

	if isDbBackupFunction(global.Args.Function) {
		return dbBackupParametersValid(global.Args.Function)
	}

	// check parameter file
	if global.Args.ParameterFile != "" {
		message := isFileValid(global.Args.ParameterFile, FILEMUSTEXIST)
		if message != "" {
			fmt.Println("Parameter", message)
			return false
		}
	}

	// Check Userid
	if global.Args.UserId == "" {
		fmt.Println("Userid must be specified.")
		return false
	}

	// Check input file
	if global.Args.InputFile != "" {
		message := isFileValid(global.Args.InputFile, FILEMUSTEXIST)
		if message != "" {
			fmt.Println("Input", message)
			return false
		}
	} else {
		fmt.Println("Input file must be specified.")
		return false
	}

	// Check output file
	if global.Args.OutputFile != "" {
		message := isFileValid(global.Args.OutputFile, FILENOTEXIST)
		if message != "" {
			fmt.Println("Output", message)
			return false
		}
	} else {
		fmt.Println("Output file must be specified.")
		return false
	}

	if global.Args.Function == global.BACKUP {
		// Checking backup id
		if global.Args.BackupId == -1 {
			fmt.Println("Function 'backup' requires a backup id.")
			return false
		}
		// Checking backup level
		if !isBackupLevelValid(global.Args.BackupLevel) {
			return false
		}
	}

	return true
}

func dbBackupParametersValid(function string) bool {
	if global.Args.EndpointUrl == "" {
		fmt.Println(
			"You must specify an valid URL for argument '-endpoint'.")
		return false
	}

	if global.Args.Region == "" {
		fmt.Println(
			"You must specify an valid region for argument '-region'.")
		return false
	}

	if isFileValid(global.Args.AuthKeypath, true) != "" {
		fmt.Println("You must specify a valid path to the " +
			"file containing your apikey.")
		return false
	}

	if function == global.BUCKET_GET_LIFECYCLE ||
		function == global.BUCKET_GET_LIST ||
		function == global.BUCKET_VERIFY {
		if global.Args.Bucket == "" {
			fmt.Printf(
				"For function '%s'"+
					"the parameter '-bucket' must be specified.",
				function,
			)
			return false
		}
	}
	if function == global.BUCKET_GET_LIFECYCLE ||
		function == global.BUCKET_GET_LIST {
		if global.Args.ResultFile == "" {
			fmt.Printf(
				"For function '%s'"+
					"the parameter '-r' must be specified.",
				function,
			)
			return false
		}

		message := isFileValid(global.Args.ResultFile, FILEMUSTEXIST)
		if message != "" {
			fmt.Println("Parameter", message)
			return false
		}
	}
	if function == global.FILE_UPLOAD {
		if global.Args.Source == "" ||
			global.Args.Key == "" {
			fmt.Println(
				"For function 'file_upload' the parameters" +
					" '-source' and '-key' must be specified.",
			)
			return false
		}
	}
	return true
}

/*
Checking if the given function is valid
*/
func isFunctionValid(function string) bool {
	if function == "" {
		return false
	}
	return slices.Contains(global.FUNCTIONLIST, function)
}

/*
Checking if the given file is valid
*/
func isFileValid(filename string, filemustexist bool) string {
	// Check if parameter file exists
	switch checkFile(filename) {
	case FILENOTFOUND:
		if filemustexist {
			message := fmt.Sprintf("File '%s' does not exist.", filename)
			return message
		}
	case FILENOTREADABLE:
		message := fmt.Sprintf("File '%s' is not readable.", filename)
		return message
	}
	return ""
}

/*
Checking a given file
*/
func checkFile(filename string) int {
	_, err := os.Stat(filename)

	if errors.Is(err, os.ErrNotExist) {
		return FILENOTFOUND
	}

	_, err = os.Open(filename)

	if err != nil {
		return FILENOTREADABLE
	}
	return 0
}

/*
Checking if a given backup level is valid
*/
func isBackupLevelValid(backupLevel string) bool {
	if backupLevel == "" {
		return false
	}

	switch backupLevel {
	case "COMPLETE", "LOG", "INCREMENTAL", "DIFFERENTIAL":
		return true
	default:
		return false
	}
}

func isDbBackupFunction(function string) bool {
	return slices.Contains(global.DBBACKUP_FUNCTIONLIST, function)
}

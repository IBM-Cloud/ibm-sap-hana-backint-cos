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
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
)

/*
Checking the configuration parameters
in case of -check argument
*/
func CheckParameters() int {
	fmt.Printf(
		"Validating parameter configuration file %s...\n\n",
		global.Args.ParameterFile,
	)

	// Generating the configuration from the parameter file
	// configuration settings and defaults.
	_, success := GenerateConfiguration(
		global.Args.ParameterFile,
	)

	for _, m := range checkParmMessages {
		// In some cases (unknown key or key in wrong section) the validation during a
		// hdbbackint agent call does not produce an error. These keys and their values are ignored.
		// In case of argument -check the agent should return an error.
		// To prevent side effects regarding the backint functionality we
		// check here the output of the validation process.
		if strings.Contains(m, VALIDATION_ERROR) {
			success = false
		}
		fmt.Println(m)
	}

	if !success {
		// Error during generating/validating the configuration
		fmt.Println("Error(s) during validation of parameter configuration file.")
		return global.WRONG_PARAMETER
	}

	fmt.Println(
		"All configuration parameters are valid.",
	)
	return global.SUCCESS
}

/*
Validating the configuration
*/
func validateConfig(basicConfig []Default) []InvalidValue {
	if validateConfigTypes(basicConfig) != nil {
		os.Exit(global.WRONG_PARAMETER)
	}

	// Validating the mandatory parameters
	if global.Args.CheckParms {
		checkParmMessages = append(checkParmMessages,
			"\nValidating existence of mandatory parameters",
		)
	}
	for _, cp := range basicConfig {
		cp.validateMandatory()
	}

	// Validating the optional parameters
	if global.Args.CheckParms {
		checkParmMessages = append(checkParmMessages,
			"\nValidating values",
		)
	}
	for _, cp := range basicConfig {
		cp.validateParameter()
	}

	validateSpecial(basicConfig)
	return invalidValues
}

/*
Validating the sections
*/
func validateSections(sections []string) {
	if global.Args.CheckParms {
		checkParmMessages = append(checkParmMessages, "Validating sections")
	}

	for i, section := range sections {
		found := false
		for _, validSection := range validSections {
			if section == validSection {
				found = true
				if global.Args.CheckParms {
					checkParmMessages = append(checkParmMessages,
						fmt.Sprintf("\t%s: Section %s is valid.", VALIDATION_OK, section),
					)
				}
			}
		}

		if !found {
			errMsg := fmt.Sprintf(
				"%s: You specified the section '%s',"+
					" but it is not part of the hdbbackint configuration. "+
					"All parameters specified in this section are ignored.",
				VALIDATION_ERROR,
				section,
			)
			if global.Args.CheckParms {
				checkParmMessages = append(checkParmMessages, "\t"+errMsg)
			} else {
				fmt.Println(errMsg)
			}

			// Removing invalid section name
			sections = append(sections[:i], sections[i+1:]...)
		}
	}
}

/*
Validating if keys are located in the correct sections
*/
func validateKeysInSections(parms []ConfigParameter) []ConfigParameter {
	if global.Args.CheckParms {
		checkParmMessages = append(
			checkParmMessages,
			"\nValidating classification of parameters",
		)
	}

	for i := range parms {
		p := &parms[i]
		keyFromFile := p.key
		sectionFromFile := p.section
		found := false
		for _, d := range configDefaults {
			if keyFromFile == d.key {
				found = true
				if sectionFromFile != d.section {
					errMsg := fmt.Sprintf("%s: You specified '%s'"+
						" in section '%s', but key belongs to section %s."+
						" The value of '%s' will be ignored.",
						VALIDATION_ERROR,
						keyFromFile,
						sectionFromFile,
						d.section,
						keyFromFile,
					)
					if global.Args.CheckParms {
						checkParmMessages = append(checkParmMessages, "\t"+errMsg)
					} else {
						fmt.Println(errMsg)
					}
					// Removing invalid config parameter
					p.ignored = true
				} else {
					if global.Args.CheckParms {
						checkParmMessages = append(
							checkParmMessages,
							fmt.Sprintf(
								"\t%s: '%s' specified in correct section",
								VALIDATION_OK,
								keyFromFile),
						)
					}
				}
				break
			}
		}
		if !found {
			errMsg := fmt.Sprintf("%s: You specified '%s'"+
				" in section '%s', but the key is unknown."+
				" The value of '%s' will be ignored.",
				VALIDATION_ERROR,
				keyFromFile,
				sectionFromFile,
				keyFromFile,
			)
			if global.Args.CheckParms {
				checkParmMessages = append(checkParmMessages, "\t"+errMsg)
			} else {
				fmt.Println(errMsg)
			}
		}
	}
	return parms
}

/*
Validating the types of the config values
*/
func validateConfigTypes(basicConfig []Default) error {
	var err error
	for _, cp := range basicConfig {
		if !contains(validConfigTypes, cp.validationType) {
			fmt.Printf(
				"INTERNAL PROGRAMMING ERROR: WRONG DATATYPE FOR '%q' SPECIFIED!",
				cp.key,
			)
			return err
		}
	}
	return nil
}

/*
Validating if a mandatory parameter is set in config file
*/
func (cp Default) validateMandatory() {
	if cp.mandatory {
		if cp.configValue == "" {
			cp.addMissingMandatoryMsg()
		} else {
			if global.Args.CheckParms {
				checkParmMessages = append(checkParmMessages,
					fmt.Sprintf(
						"\t%s: Mandatory parameter '%s' exists.",
						VALIDATION_OK,
						cp.key,
					))
			}
		}
	}
}

/*
Validating single parameter if set
*/
func (cp Default) validateParameter() {
	if cp.configValue != "" {
		cp.validateConfigValue()
	}
}

/*
Validating special settings:
lock retention belong to more than one parameter
*/
func validateSpecial(basicConfig []Default) {
	validateLockRetention(basicConfig)
	validateAuthMode(basicConfig)
}

/*
Validating one single config value depending on datatype
*/
func (cp Default) validateConfigValue() {
	switch cp.validationType {
	case CONFIG_BOOL:
		cp.validateBool()
	case CONFIG_CHUNKSIZE:
		cp.validateChunksize()
	case CONFIG_FILE:
		cp.validateFile()
	case CONFIG_INT:
		cp.validateInt()
	case CONFIG_LIST:
		cp.validateList()
	case CONFIG_PERIOD:
		cp.validatePeriod()
	case CONFIG_RANGE:
		if cp.validateInt() {
			cp.validateRange()
		}
	case CONFIG_TAG:
		cp.validateTag()
	case CONFIG_URL:
		cp.validateUrl()
	}
}

/*
Validating boolean value
*/
func (cp Default) validateBool() {
	if !cp.isBool() {
		cp.addInvalidValueMsg(
			"The value you specified must be of type 'boolean'.",
		)
		return
	}
	addOkMessage(cp.key)
}

/*
Validating chunksize value
*/
func (cp Default) validateChunksize() {
	errMsg := "The value you specified does not have to correct format." +
		" It must be either an integer value or must have the format: " +
		" <size><unit> while <unit> must be either 'KB', 'MB', or 'GB' and " +
		" <size> must not be 0 or undefined."

	// Checking first if the value is an integer
	_, err := strconv.Atoi(cp.configValue)
	if err == nil {
		addOkMessage(cp.key)
		return
	}

	// Splitting value in size and unit
	// Format: <size><unit> while unit is
	// KB, MB or GB (case unsensitive)
	if len(cp.configValue) < 3 {
		cp.addInvalidValueMsg(errMsg)
		return
	}

	size, unitU := getChunksizeSizeAndUnit(cp.configValue)

	// Checking if a valid size is specified
	if size == "0" {
		cp.addInvalidValueMsg(errMsg)
		return
	}

	// Checking if a valid unit is specified
	if !contains(validSizeUnits, unitU) {
		cp.addInvalidValueMsg(errMsg)
		return
	}
	cp.configValue = calculateChunksizeInBytes(size, unitU)
	fmt.Println(cp.configValue)
	addOkMessage(cp.key)
}

/*
Validating file
*/
func (cp Default) validateFile() {
	info, err := os.Stat(cp.configValue)
	if errors.Is(err, os.ErrNotExist) {
		cp.addInvalidValueMsg(
			"The file you specified does not exist.",
		)
		return
	}

	if info != nil {
		if info.Mode().Perm()&0400 != 0400 {
			cp.addInvalidValueMsg(
				"The file you specified does not have read permissions.",
			)
			return
		}
	}
	// TODO?
	addOkMessage(cp.key)
}

/*
Validating integer value
Function is also used for datatype = "range"
*/
func (cp Default) validateInt() bool {
	_, err := strconv.Atoi(cp.configValue)
	if err != nil {
		cp.addInvalidValueMsg(
			"You did not specify an integer value.",
		)
		return false
	}
	addOkMessage(cp.key)
	return true
}

/*
Validating a value from a list
*/
func (cp Default) validateList() {
	if !contains(cp.possibleValues, cp.configValue) {
		message := "It must be one of the following:"
		for _, pv := range cp.possibleValues {
			if global.Args.CheckParms {
				message += fmt.Sprintf("\n\t\t%s", pv)
			} else {
				message += fmt.Sprintf("\n\t%s", pv)
			}
		}
		cp.addInvalidValueMsg(message)
		return
	}
	addOkMessage(cp.key)
}

/*
Validating object_lock_retention_period
*/
func (cp Default) validatePeriod() {
	message := "The value you specified for" +
		" 'object_lock_retention_period' does not" +
		" have the correct format." +
		" It must be a comma separated value string" +
		" while the first position represents the years," +
		" the second position the months and the" +
		" third position the days." +
		" All values must be integers."
	// period must have the format: year,month,day
	splitted := strings.Split(cp.configValue, ",")
	if len(splitted) != 3 {
		cp.addInvalidValueMsg(message)
		return
	} else {
		for _, p := range splitted {
			_, err := strconv.Atoi(p)
			if err != nil {
				cp.addInvalidValueMsg(message)
				break
			}
		}
	}
	// TODO?
	addOkMessage(cp.key)
}

/*
Validating a range
*/
func (cp Default) validateRange() {
	valueI := global.ToInteger(cp.configValue)
	if !isValueInRange(valueI, cp.min, cp.max) {
		message := fmt.Sprintf(
			"It must be"+
				" between '%d' and '%d'.",
			cp.min,
			cp.max,
		)
		cp.addInvalidValueMsg(message)
		return
	}
	// TODO?
	addOkMessage(cp.key)
}

/*
Validating tag
*/
func (cp Default) validateTag() {
	// tag has the format: "tag1=val1,tag2=val2"
	tags := strings.Split(cp.configValue, ",")
	// Validate max. number of tags
	if len(tags) > MAX_NUMBER_OF_TAGS {
		message := fmt.Sprintf(
			"You specified '%d' number of different tags, ",
			len(tags),
		)
		message += fmt.Sprintf(
			"it must not exceed '%d'.",
			MAX_NUMBER_OF_TAGS,
		)
		cp.addInvalidValueMsg(message)
		return
	} else {
		// Validate tag format
		for _, tag := range tags {
			t := strings.Split(tag, "=")
			if len(t) == 1 {
				// No "=" in tag -> Error
				message := fmt.Sprintf("You specified '%q' as a tag."+
					" The format of the tag is wrong."+
					" It must be: 'tag=val'",
					tag)
				cp.addInvalidValueMsg(message)
				return
			}
		}
	}
	addOkMessage(cp.key)
}

/*
Validating an url
*/
func (cp Default) validateUrl() {
	urlcmp := "https://s3."
	if !strings.HasPrefix(cp.configValue, urlcmp) {
		message := "You did not specify a valid url."
		message += "The value must start with "
		message += urlcmp
		cp.addInvalidValueMsg(message)
		return
	}
	addOkMessage(cp.key)
}

/*
Special validation:
Validating object lock retention
*/
func validateLockRetention(basicConfig []Default) {
	if isObjectLockRetentionMode(basicConfig) {
		if !isObjectLockRetentionPeriod(basicConfig) {
			message := fmt.Sprintf("%s: You specified 'cmp' or 'gov' for 'object_lock_retention_mode', "+
				"but no 'object_lock_retention_period' is specified.",
				VALIDATION_ERROR,
			)
			Default{}.addInvalidValueMsg(message)
		}
	} else {
		if isObjectLockRetentionPeriod(basicConfig) {
			message := fmt.Sprintf("%s: You did not specify 'object_lock_retention_mode' "+
				"or 'object_lock_retention_mode' is set to 'None', "+
				"but 'object_lock_retention_period' is specified.",
				VALIDATION_ERROR,
			)
			Default{}.addInvalidValueMsg(message)
		}
	}
}

/*
	Special validation:

Validation of authorization mode.
If auth_mode = apikey, auth_keypath must be set and valid
If iam_profile_id, or iam_profile_crn is set, auth_mode must be oauth
*/
func validateAuthMode(basicConfig []Default) {
	authMode := strings.ToLower(getObjForKey(basicConfig, auth_mode.key).configValue)

	switch authMode {
	case AUTH_TRUSTEDPROFILE:
		if isSetAuthKeypath(basicConfig) {
			message := fmt.Sprintf("%s: You did not specify '%s' for 'auth_mode', "+
				" but you did set 'auth_key_path'.",
				VALIDATION_ERROR,
				AUTH_APIKEY,
			)
			Default{}.addInvalidValueMsg(message)
		}
		return
	case AUTH_APIKEY:
		if isSetIamProfileId(basicConfig) {
			message := fmt.Sprintf("%s: You did specify '%s', "+
				" but 'auth_mode' is not set to '%s'.",
				VALIDATION_ERROR,
				iam_profile_id.key,
				AUTH_TRUSTEDPROFILE,
			)
			Default{}.addInvalidValueMsg(message)
		}

		if !isSetAuthKeypath(basicConfig) {
			message := fmt.Sprintf("%s: You did not specify 'auth_keypath', "+
				"but 'auth_mode' is set to '%s'.",
				VALIDATION_ERROR,
				AUTH_APIKEY,
			)
			Default{}.addInvalidValueMsg(message)
		}
		return
	default:
		return
	}
}

/*
returns true if config value is of type boolean
*/
func (cp Default) isBool() bool {
	_, err := strconv.ParseBool(cp.configValue)
	return err == nil
}

/*
Returns true if auth_keypath is set
*/
func isSetAuthKeypath(basicConfig []Default) bool {
	return isSet(basicConfig, auth_keypath.key)
}

/*
Returns true uf metadata_service_url is set
*/
func isSetIamProfileId(basicConfig []Default) bool {
	return isSet(basicConfig, iam_profile_id.key)
}

/*
Checks if a parameter is set in config file
*/
func isSet(basicConfig []Default, key string) bool {
	b := getObjForKey(basicConfig, key)
	return b.configValue != ""
}

/*
Returns true if object_lock_retention_mode is set
*/
func isObjectLockRetentionMode(basicConfig []Default) bool {
	b := getObjForKey(basicConfig, "object_lock_retention_mode")
	val := strings.ToLower(b.configValue)
	return (val == global.LOCKMODE_CMP || val == global.LOCKMODE_GOV)
}

/*
Returns true if object_lock_retention_period is set
*/
func isObjectLockRetentionPeriod(basicConfig []Default) bool {
	b := getObjForKey(basicConfig, "object_lock_retention_period")
	return b.configValue != ""
}

/*
Returns true if a value is object of a slice
*/
func contains(possibleValues []string, value string) bool {
	return slices.Contains(possibleValues, value)
}

/*
Returns the basicConfig object for a given key
*/
func getObjForKey(basicConfig []Default, key string) Default {
	for _, b := range basicConfig {
		if b.key == key {
			return b
		}
	}
	return Default{}
}

/*
Returns true if a given value is within a given range
*/
func isValueInRange(value int, min int, max int) bool {
	if value < min || value > max {
		return false
	}
	return true
}

/*
Adding error message for missing mandatory parameter
*/
func (cp Default) addMissingMandatoryMsg() {
	message := fmt.Sprintf("%s: You did not specify a value"+
		" for the mandatory parameter"+
		" '%s'.",
		VALIDATION_ERROR,
		cp.key,
	)
	if global.Args.CheckParms {
		checkParmMessages = append(checkParmMessages, "\t"+message)
	}
	invalid := InvalidValue{
		errorMessage: message,
		invalidParm:  cp,
	}
	invalidValues = append(invalidValues, invalid)
}

/*
Adding error message to the invalidValues slice
*/
func (cp Default) addInvalidValueMsg(msg string) {
	message := ""
	if cp.configValue != "" && cp.key != "" {
		message = fmt.Sprintf(
			"%s: '%s': the value '%s' you specified is invalid. ",
			VALIDATION_ERROR,
			cp.key,
			cp.configValue,
		)
	}
	message += msg

	if global.Args.CheckParms {
		checkParmMessages = append(checkParmMessages, "\t"+message)
	}

	invalid := InvalidValue{
		errorMessage: message,
		invalidParm:  cp,
	}
	invalidValues = append(invalidValues, invalid)
}

func addOkMessage(key string) {
	if global.Args.CheckParms {

		checkParmMessages = append(checkParmMessages,
			fmt.Sprintf(
				"\t%s: Parameter '%s' exists"+
					" and its value is valid.",
				VALIDATION_OK,
				key,
			))
	}
}

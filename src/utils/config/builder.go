// Generating the internal backint configuration from the arguments and the hdbbackint.cfg parameter file
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/bigkevmcd/go-configparser"
)

/*
Generating the hdbbackint configuration from config file,
validating the values and setting defaults
*/
func GenerateConfiguration(
	configFilename string,
) (BackintConfigT, bool) {
	basicConfig := buildBasicConfig(configFilename)
	invalidValues = validateConfig(basicConfig)

	if len(invalidValues) > 0 {
		if !global.Args.CheckParms {
			// Don't print messages in case of -check,
			// The messages will be printed in a different way
			for _, v := range invalidValues {
				fmt.Println(v.errorMessage)
			}
		}
		return nil, false
	}

	BackintConfig = updateConfigWithDefaults(basicConfig)

	BackintConfig = updateConfigWithApikey(BackintConfig)
	return BackintConfig, true
}

/*
Building the internal basic configuration
*/
func buildBasicConfig(
	configFilename string,
) []Default {
	configParms := readConfigfile(configFilename)
	configParms = addConfigParmsFromTooloption(configParms)
	configuration := updateWithConfigValues(configParms)
	return configuration
}

/*
Reading the config file
*/
func readConfigfile(filename string) []ConfigParameter {
	parser, err := configparser.Parse(filename)
	global.CheckForError(
		err,
		fmt.Sprintf("Error reading config file '%s'.", filename),
		global.WRONG_PARAMETER,
	)

	sections := parser.Sections()

	// Checking if unknown sections are found in config file
	// Parameters in these sections are ignored
	validateSections(sections)

	var parmsFromConfig []ConfigParameter

	for _, section := range sections {
		items, _ := parser.Items(section)
		keys := items.Keys()
		for _, key := range keys {
			configParm := ConfigParameter{
				section: section,
				key:     key,
				value:   items[key],
				ignored: false,
			}
			parmsFromConfig = append(parmsFromConfig, configParm)
		}
	}

	validateKeysInSections(parmsFromConfig)

	return parmsFromConfig
}

/*
Adding config parm from input file if TOOLOPTION is used as keyword
*/
func addConfigParmsFromTooloption(
	parmsFromConfig []ConfigParameter,
) []ConfigParameter {

	for _, inputKeyword := range global.InputFileContent {
		if inputKeyword.Keyword != "TOOLOPTION" {
			continue
		}
		// TODO -> issue #20
	}
	return parmsFromConfig
}

/*
Updating the internal configuration with the values from config file
*/
func updateWithConfigValues(parmsFromConfig []ConfigParameter) []Default {
	for _, cfgParm := range parmsFromConfig {
		if !contains(validSections, cfgParm.section) {
			// Ignoring all parameters set in invalid sections
			continue
		}
		if cfgParm.ignored {
			// Ignoring parameter if set in wrong section
			continue
		}
		found := cfgParm.updateMatchingObj()
		if !found && !global.Args.CheckParms {
			message := fmt.Sprintf(
				"ERROR: You specified '%s' in section '%s',"+
					" but the parameter does not exist in the definition."+
					" The parameter does not have any effect.",
				cfgParm.key,
				cfgParm.section,
			)

			invalidValues = append(invalidValues,
				InvalidValue{
					errorMessage: message,
					invalidParm:  Default{},
				})
		}
	}
	return configDefaults
}

/*
Updating the configuration parameter with value from config file
*/
func (cfgParm ConfigParameter) updateMatchingObj() bool {
	for i, obj := range configDefaults {
		if obj.section == cfgParm.section && obj.key == cfgParm.key {
			// Special case for multipart_chunksize:
			// value must be calculated if a size unit is specified
			if cfgParm.key == multipart_chunksize.key {
				size, unitU := getChunksizeSizeAndUnit(cfgParm.value)
				configDefaults[i].configValue = calculateChunksizeInBytes(size, unitU)
			} else {
				configDefaults[i].configValue = cfgParm.value
			}

			return true
		}
	}
	return false
}

/*
Updating all parameters not set by config file with defaults
*/
func updateConfigWithDefaults(basicConfig []Default) BackintConfigT {
	backintConfig := make(BackintConfigT)
	for _, cp := range basicConfig {
		value := cp.configValue
		if value == "" {
			value = cp.defaultValue
		}

		if value != "" {
			// Don't add parameters with a "None" default
			backintConfig.set(cp.key, value)
		}
	}
	return backintConfig
}

/*
Reading the apikey from file "auth_keypath" and storing the value in map
*/
func updateConfigWithApikey(backintConfig BackintConfigT) BackintConfigT {
	apikey, err := global.ReadApikeyFromFile(backintConfig.authKeypath())
	if err != nil {
		fmt.Printf("Could not discover the apikey."+
			" Check if file '%s' is available and contains the apikey.",
			backintConfig.authKeypath(),
		)
		os.Exit(global.WRONG_PARAMETER)
	}
	backintConfig.set("apikey", apikey)

	return backintConfig
}

/*
Calculating the chunksize in case a unit is specified
Returning the value as string
*/
func calculateChunksizeInBytes(size string, unit string) string {
	sizeI, _ := strconv.Atoi(size)
	switch unit {
	case UNIT_KB:
		sizeI = sizeI * 1024
	case UNIT_MB:
		sizeI = sizeI * 1024 * 1024
	case UNIT_GB:
		sizeI = sizeI * 1024 * 1024 * 1024
	}
	return fmt.Sprintf("%d", sizeI)
}

func getChunksizeSizeAndUnit(chunksize string) (string, string) {
	_, err := strconv.Atoi(chunksize)
	if err == nil {
		return chunksize, ""
	}
	size := chunksize[:len(chunksize)-2]
	unit := chunksize[len(chunksize)-2:]
	unitU := strings.ToUpper(unit)

	return size, unitU
}

/*
Contains all datatypes for configuration
*/
package config

// Datatype representing one parameter for backint configuration
type ConfigParameter struct {
	section string
	key     string
	value   string
	ignored bool
}

// Datatype holding an invalid value and its appropriate message
type InvalidValue struct {
	errorMessage string
	invalidParm  Default
}

// Datatype representing the the defaults of the backint configuration
type Default struct {
	key            string
	section        string
	validationType string
	mandatory      bool
	possibleValues []string
	defaultValue   string
	min            int
	max            int
	configValue    string
}

// Datatype representing one single backint configuration value
type BackintConfigT map[string]string

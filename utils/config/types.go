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

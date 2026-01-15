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

/*
Contains all variables for configuration
*/
package config

// backint configuration
var BackintConfig BackintConfigT

// Slice containing all invalid values and its messages
var invalidValues []InvalidValue

// Slice containing all messages generated during -check
var checkParmMessages []string

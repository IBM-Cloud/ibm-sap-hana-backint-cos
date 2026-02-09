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

package tests

import (
	"testing"
)

// Test functions

func TestCheckConfig(t *testing.T) {
	configTests := setupConfigTests(t)
	executablePath := getExecutablePath()

	for _, cfgTest := range configTests {
		t.Run(cfgTest.name, func(t *testing.T) {
			output, err := runCommand(t, executablePath, "-p", cfgTest.cfgPath, "-check")

			if cfgTest.shouldSucceed {
				if err != nil {
					t.Error("hdbbackint -check should succeed but failed")
				} else {
					t.Log("Config check passed")
				}
			} else {
				if err == nil {
					t.Error("hdbbackint -check should fail but succeeded")
				} else {
					if errMsgOk(string(output), cfgTest.msgToCheck) {
						t.Log("Config check passed")
					} else {
						t.Errorf("Wrong error message(s): \nOutput: %s", output)
					}
				}
			}
		})
	}
}

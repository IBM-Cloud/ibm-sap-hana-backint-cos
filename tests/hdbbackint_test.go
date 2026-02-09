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

func TestVersion(t *testing.T) {
	output, err := runCommand(t, getExecutablePath(), "-v")
	if err != nil {
		t.Errorf("hdbbackint -v failed: %v\nOutput: %s", err, output)
	}

	if len(output) == 0 {
		t.Error("expected version output, got empty output")
	}

	t.Logf("Version output: %s", string(output))
}

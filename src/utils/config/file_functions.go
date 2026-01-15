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
	"bufio"
	"os"
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
)

/*
Reading the input file content
*/
func ReadInputFile(filePath string) []global.InputFileContentT {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}

	defer func() {
		_ = f.Close()
	}()

	var inputFileContentList []global.InputFileContentT

	fScanner := bufio.NewScanner(f)
	for fScanner.Scan() {
		line := fScanner.Text()
		if !strings.HasPrefix(line, "#") {
			// Ignore line
			continue
		}

		if strings.HasPrefix(strings.ToUpper(line), "#SOFTWAREID") {
			// Ignore line
			continue
		}
		inputFileContent := global.InputFileContentT{}
		keyword := line
		if strings.Contains(line, " ") {
			splitted := strings.SplitN(line, " ", 2)
			keyword = splitted[0]
			inputFileContent.Parameter = strings.ReplaceAll(splitted[1], "\"", "")
		}
		inputFileContent.Keyword = strings.ToUpper(
			strings.ReplaceAll(keyword, "#", ""),
		)
		inputFileContentList = append(inputFileContentList, inputFileContent)
	}

	return inputFileContentList
}

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

package global

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*
Getting the integer representation of a string
*/
func ToInteger(value string) int {
	valInt, err := strconv.Atoi(value)
	if err != nil {
		return 0
	} else {
		return valInt
	}
}

/*
Getting the string representation of an int64 value
*/
func ToString(val int64) string {
	return fmt.Sprintf("%d", val)
}

/*
Compress TODO -> Issue #7
*/
// func Compress(src []byte) []byte {
// 	return encoder.EncodeAll(src, make([]byte, 0, len(src)))
// }

/*
Check error and set OS Exit code
*/
func CheckForError(err error, message string, rc int) {
	if err != nil {
		if Logger == nil {
			fmt.Printf("%s Error: %s\n", message, err)
		} else {
			Logger.Error(fmt.Sprintf("%s Error: %s\n", message, err))
		}
		os.Exit(rc)
	}
}

func ReadApikeyFromFile(authKeypath string) (string, error) {
	authFile, err := os.Open(authKeypath)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = authFile.Close()
	}()

	scanner := bufio.NewScanner(authFile)
	var fileContent []string
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}
	if len(fileContent) > 1 || len(fileContent) == 0 {
		return "", nil
	}

	return fileContent[0], nil
}

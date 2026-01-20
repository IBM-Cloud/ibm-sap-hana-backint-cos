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

package logging

import (
	"fmt"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/version"
)

/*
Initializing the BackintResultMessages
*/
func InitializeBackintResultMessages() BackintResultMessages {
	b := make(BackintResultMessages, 0)
	message := fmt.Sprintf(
		"#%s \"%s\" \"%s\"",
		"SOFTWAREID",
		version.BACKINT_VERSION,
		version.TOOL_VERSION,
	)
	b = append(b, message)
	return b
}

/*
Writing the messages to the outfile as defined in the arguments list
*/
func (b BackintResultMessages) Dump() {
	for _, message := range b {
		_, _ = fmt.Fprintln(GetLogFile(), message)
	}
}

/*
Adding a message with a keyword
*/
func (b *BackintResultMessages) AddKeyword(keyword string, parms []string) {
	message := fmt.Sprintf("#%s ", keyword)
	for _, a := range parms {
		message = message + fmt.Sprintf("\"%s\" ", a)
	}
	*b = append(*b, message)
}

/*
Adding one or more comments
*/
func (b *BackintResultMessages) addComments(comments []string) {
	for _, m := range comments {
		*b = append(*b, m)
	}
}

/*
Adding the success message for BACKUP
*/
func (b *BackintResultMessages) AddBackupSuccessMessage(
	ETag string,
	sourcePath string,
	sourceSize int64,
) {
	keyword := "SAVED"
	parms := []string{ETag, sourcePath, global.ToString(sourceSize)}
	b.AddKeyword(keyword, parms)
}

/*
Adding the metrics comment for BACKUP
*/
func (b *BackintResultMessages) AddBackupMetrics(
	sourceSize int64,
	targetSize int64,
	duration float64,
) {
	comment := fmt.Sprintf("metrics: source: %d, destination: %d, seconds: %f",
		sourceSize,
		targetSize,
		duration,
	)
	b.addComments([]string{comment})
}

/*
Adding an error message
*/
func (b *BackintResultMessages) AddErrorMessage(
	sourcePath string,
	err error,
) {
	keyword := "ERROR"
	parms := []string{sourcePath, fmt.Sprintf("%s", err)}
	b.AddKeyword(keyword, parms)
}

/*
Adding the success message for RESTORE
*/
func (b *BackintResultMessages) AddRestoreSuccessMessage(
	ETag string,
	sourcePath string,
) {
	keyword := "RESTORED"
	parms := []string{ETag, sourcePath}
	b.AddKeyword(keyword, parms)
}

/*
Adding the OBJECT NOT FOUND message
*/
func (b *BackintResultMessages) AddObjectNotFoundMessage(
	sourcePath string,
) {
	keyword := "NOTFOUND"
	parms := []string{sourcePath}
	b.AddKeyword(keyword, parms)
}

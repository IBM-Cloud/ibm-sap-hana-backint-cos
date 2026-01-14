// Functions to generate and print messages in the format required by HANA backup/restore process
// The HANA backup/restore requires a specific format of messages in the trace file as indicator if a function
// was executed successfully or with error.
package logging

import (
	"fmt"
	"hdbbackint/utils/global"
)

/*
Initializing the BackintResultMessages
*/
func InitializeBackintResultMessages() BackintResultMessages {
	b := make(BackintResultMessages, 0)
	message := fmt.Sprintf(
		"#%s \"%s\" \"%s\"",
		"SOFTWAREID",
		global.BACKINT_VERSION,
		global.TOOL_VERSION,
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
Adding the compressed info comment for BACKUP
*/
func (b *BackintResultMessages) AddBackupCompressedInfo(
	sourceSize int64,
	targetSize int64,
) {
	comment := fmt.Sprintf(
		"compressed backup: original size: %s, compressed size: %s, factor: %.2f",
		printableSize(sourceSize),
		printableSize(targetSize),
		float32(sourceSize/targetSize),
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

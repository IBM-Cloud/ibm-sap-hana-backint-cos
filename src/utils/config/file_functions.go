/*
Functions for handling the input file content
*/
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

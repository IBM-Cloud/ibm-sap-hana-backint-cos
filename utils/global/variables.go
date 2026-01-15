/*
Contains all global variables
*/
package global

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Command line arguments
var Args CommandLineArguments

// Input file contents
var InputFileContent []InputFileContentT

// Logfile information
var LogFile *os.File
var Logger *logrus.Logger

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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/config"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/sirupsen/logrus"
)

/*
Setting the format of a log message
*/
func (f *backintFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	time := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())
	function := getFunctionName(entry.Caller.Function)
	line := entry.Caller.Line
	filename := getGoFileName(entry.Caller.File)
	exec := fmt.Sprintf("%s - %s:%d", filename, function, line)
	return fmt.Appendf(
			nil,
			"[%s] - %s - %s - %s\n",
			time,
			exec,
			level,
			entry.Message),
		nil
}

/*
Getting the file representation of the logfile
*/
func GetLogFile() *os.File {
	var err error
	global.LogFile, err = os.OpenFile(
		global.Args.OutputFile,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0666,
	)
	if err != nil {
		fmt.Printf(
			"Could not open logfile '%s' for writing.",
			global.Args.OutputFile,
		)
	}
	return global.LogFile
}

/*
Setting up logging
*/
func SetupLogging() *logrus.Logger {
	logger := generateLogger()
	logger.Info(fmt.Sprintf(
		"Running hdbbackint with %s.",
		global.TOOL_VERSION),
	)
	return logger
}

/*
Writing the backint configuration and the input file content
*/
func WriteBackintInfo(logger *logrus.Logger) {
	writeBackintConfiguration(logger)
	logInputFile(logger)
}

/*
Generating the logger with formatting
*/
func generateLogger() *logrus.Logger {
	log := &logrus.Logger{
		Out:   GetLogFile(),
		Level: getLogLevel(),
		Formatter: &backintFormatter{logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05,123",
			DisableLevelTruncation: true,
		}},
	}

	// Setting the caller's information in logger.entry.Caller
	log.SetReportCaller(true)
	return log
}

/*
Getting the loglevel from the given loglevel config parameter
*/
func getLogLevel() logrus.Level {
	fmt.Printf("Loglevel: %s", config.BackintConfig.AgentLogLevelU())
	switch config.BackintConfig.AgentLogLevelU() {
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "CRITICAL":
		return logrus.PanicLevel
	case "WARNING":
		return logrus.WarnLevel
	case "HTTP":
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}

/*
Getting the go file name where the log information is written
*/
func getGoFileName(filename string) string {
	splitted := strings.Split(filename, "/")
	return splitted[len(splitted)-1]
}

/*
Getting the go function name where the log information is written
*/
func getFunctionName(function string) string {
	splitted := strings.Split(function, ".")
	return splitted[len(splitted)-1]
}

/*
Writing the backint configuration
*/
func writeBackintConfiguration(logger *logrus.Logger) {
	logger.Info("Using backint configuration settings: ")
	logger.Info("=================================================================")
	for key, value := range config.BackintConfig {
		if key == "timeout_microsecond" {
			// Don't print the timeout to log file
			continue
		}
		if key == "apikey" {
			// Don't print the timeout to log file
			logger.Info(key + " = ****")
			continue
		}
		logger.Info(key + " = " + value)
	}
	logger.Info("=================================================================")
}

/*
Writing the input file content
*/
func logInputFile(logger *logrus.Logger) {
	f, _ := os.Open(global.Args.InputFile)
	defer func() {
		_ = f.Close()
	}()

	logger.Debug("Content of input file:")
	logger.Debug("=================================================================")
	fScanner := bufio.NewScanner(f)
	for fScanner.Scan() {
		logger.Debug(fScanner.Text())
	}
	logger.Debug("=================================================================")
}

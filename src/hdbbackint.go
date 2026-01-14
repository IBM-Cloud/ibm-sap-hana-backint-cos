package main

import (
	"fmt"
	"hdbbackint/utils/backint"
	"hdbbackint/utils/config"
	"hdbbackint/utils/cos"
	"hdbbackint/utils/global"
	"hdbbackint/utils/logging"
	"hdbbackint/utils/snappy"
	"os"
)

func main() {

	// Reading and validating the command line arguments
	argsValid := true //nolint:ineffassign
	dbBackupFunction := false
	global.Args, argsValid, dbBackupFunction = config.GetCommandLineArguments()

	if !argsValid {
		fmt.Println("Error reading command line arguments.")
		os.Exit(global.WRONG_PARAMETER)
	}

	//Printing hdbbackint version and exit
	if global.Args.Version {
		config.PrintVersion()
		os.Exit(global.SUCCESS)
	}

	// Printing info in case of -check argument
	if global.Args.CheckParms {
		exitCode := config.CheckParameters()
		os.Exit(exitCode)
	}

	// Executing functions called by snappy agent and exit
	if dbBackupFunction {
		if !snappy.Execute(global.Args.Function) {
			os.Exit(global.FAILURE)
		}
		os.Exit(global.SUCCESS)
	}

	// Reading the input file
	// The input file contains information of the objects to be
	// backed up / restored.
	// The format of the input file contents depends on the function to be executed.
	global.InputFileContent = config.ReadInputFile(global.Args.InputFile)
	if global.InputFileContent == nil {
		fmt.Println("Error: the input file is empty or could not be read.")
		os.Exit(global.WRONG_PARAMETER)
	}

	// Generating the configuration from the parameter file
	// configuration settings and defaults.
	config.BackintConfig, _ = config.GenerateConfiguration(
		global.Args.ParameterFile,
	)
	if config.BackintConfig == nil {
		// Error during generating/validating the configuration
		fmt.Println("Error generating the configuration.")
		os.Exit(global.WRONG_PARAMETER)
	}

	// Setting up the logger
	global.Logger = logging.SetupLogging()
	logging.WriteBackintInfo(global.Logger)

	// Initializing the variable which holds the messages to print out
	// These messages must have a pre-defined format for the HANA system
	// to recognize the results of the functions.
	logging.BackintResultMsgs = logging.InitializeBackintResultMessages()

	// Setting up the connection to IBM Cloud Object Storage
	s3Session, s3Client := cos.GenerateCOSSession()

	// Checking the existence of the given bucket
	if !cos.BucketExists(s3Client) {
		global.Logger.Error(fmt.Sprintf(
			"Bucket '%s' does not exist.",
			config.BackintConfig.BucketName()),
		)
		os.Exit(global.FAILURE)
	}

	// Checking if versioning is enabled for given bucket
	if !cos.IsBucketVersioning(s3Client, config.BackintConfig.BucketName()) {
		global.Logger.Error(fmt.Sprintf(
			"Versioning must be enabled for bucket '%s'.",
			config.BackintConfig.BucketName()),
		)
		os.Exit(global.FAILURE)
	}

	// Executing the given function
	success := true
	switch global.Args.Function {
	case global.BACKUP:
		success = backint.Backup(s3Session, s3Client)
	case global.DELETE:
		success = backint.DeleteCloudObjects(s3Client)
	case global.INQUIRE:
		success = backint.Inquire(s3Client)
	case global.RESTORE:
		success = backint.Restore(s3Client)
	}

	// Dumping the backint result messages to the log file
	logging.BackintResultMsgs.Dump()

	// Closing the logfile
	global.Logger.Writer().Close()

	if success {
		os.Exit(global.SUCCESS)
	}
	os.Exit(global.FAILURE)
}

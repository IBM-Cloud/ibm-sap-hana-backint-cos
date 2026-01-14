/*
Contains all global datatypes
*/
package global

// Datatype representing all command line arguments
type CommandLineArguments struct {
	ParameterFile   string
	UserId          string
	Function        string
	InputFile       string
	OutputFile      string
	BackupId        int
	NumberOfObjects int
	BackupLevel     string
	Version         bool
	CheckParms      bool

	// Arguments used in case hdbbackint is called by snappy agent
	AuthKeypath  string
	AuthEndpoint string
	Region       string
	EndpointUrl  string
	Bucket       string
	Source       string
	Key          string
	ResultFile   string
}

// Datatype representing the keyword and its parameter from the input file
type InputFileContentT struct {
	Keyword   string
	Parameter string
}

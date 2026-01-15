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

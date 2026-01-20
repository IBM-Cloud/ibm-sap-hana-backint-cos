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

package cos

import (
	"io"
	"os"
	"time"
)

// Datatype representing the HTTP settings
type HTTPClientSettings struct {
	Connect          time.Duration
	ConnKeepAlive    time.Duration
	ExpectContinue   time.Duration
	IdleConn         time.Duration
	MaxAllIdleConns  int
	MaxHostIdleConns int
	ResponseHeader   time.Duration
	TLSHandshake     time.Duration
	MaxConnsPerHost  int
}

// Datatype representing the result for one Cloud Object Storage action (Upload/Download/Delete)
type Result struct {
	Err        error
	Duration   float64
	SourceSize int64
	TargetSize int64
	SourcePath string
	Key        string
	ETag       string
}

// Datatype representing information of one IBM Cloud Object Storage Object
type CosObject struct {
	ETag        string
	Key         string
	Destination string
	Found       bool
	Status      string
	NextIndex   *int64
}

// Datatype representing the information of one part for downloading an object
type DownloadPart struct {
	Key        string
	numParts   int64
	partNumber int64
	byteRange  string
}

// Datatype representing the result of downloading one part of an object
type DownloadPartResult struct {
	partNumber int64
	size       int64
	err        error
}

// Datatype representing the parameters for the runDownloadSinglePart call
type DownloadSingePart struct {
	fifo               *os.File
	partsNotYetWritten ByteMap
	downloadPart       DownloadPart
	nextIndex          *int64
	pipeBufferSize     int
}

// Representation for easier handling of a Bytemap
type ByteMap map[int64][]byte

// Type for reading from pipe directly to upload data
type backintReader struct {
	r         io.Reader
	noOfbytes int64
}

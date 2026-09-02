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

import "time"

// HTTPClientSettings holds the timeouts and pool sizes for the HTTP client.
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

// Result represents the outcome of one Cloud Object Storage action (Upload/Download/Delete).
type Result struct {
	Err        error
	Duration   float64
	SourceSize int64
	TargetSize int64
	SourcePath string
	Key        string
	ETag       string
}

// CosObject holds the information about one IBM Cloud Object Storage object.
type CosObject struct {
	ETag        string
	Key         string
	SourcePath  string
	Destination string
	Found       bool
	Status      string
	NextIndex   *int64
}

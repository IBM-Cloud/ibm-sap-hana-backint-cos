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
	"net"
	"net/http"
	"time"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/config"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/logging"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
)

/*
Generating the session and the client to access the IBM Cloud Object Storage
*/
func GenerateCOSSession() (*session.Session, *s3.S3) {
	cfg := setupCosConfig()
	s3Session := session.Must(session.NewSession(cfg))
	s3Client := s3.New(s3Session)
	return s3Session, s3Client
}

/*
Setting up the Cloud Object Storage Configuration
*/
func setupCosConfig() *aws.Config {
	httpClient, err := NewHTTPClientWithSettings(HTTPClientSettings{
		Connect:          60 * time.Second,
		ExpectContinue:   1 * time.Second,
		IdleConn:         60 * time.Second,
		ConnKeepAlive:    60 * time.Second,
		MaxHostIdleConns: 10,
		ResponseHeader:   50 * time.Second,
		TLSHandshake:     50 * time.Second,
	})

	global.CheckForError(
		err,
		"Error creating the customized HTTP client",
		global.FAILURE,
	)

	var apikey string
	var region string
	var endpoint string
	var authEndpoint string
	var authMethod string

	if config.BackintConfig != nil {
		apikey = config.BackintConfig.Apikey()
		region = config.BackintConfig.Region()
		endpoint = config.BackintConfig.EndpointUrl()
		authEndpoint = config.BackintConfig.IBMAuthEndpoint()
		authMethod = config.BackintConfig.AuthMethod()

	} else {
		apikey, _ = global.ReadApikeyFromFile(global.Args.AuthKeypath)
		region = global.Args.Region
		endpoint = global.Args.EndpointUrl
		authEndpoint = global.Args.AuthEndpoint
		authMethod = config.AUTH_APIKEY
	}

	var creds *credentials.Credentials

	switch authMethod {
	case config.AUTH_APIKEY:
		creds = ibmiam.NewStaticCredentials(aws.NewConfig(),
			authEndpoint,
			apikey,
			"",
		)
	default:
		break
	}

	cfg := aws.NewConfig()
	cfg = cfg.WithRegion(region)
	cfg = cfg.WithEndpoint(endpoint)
	cfg = cfg.WithCredentials(creds)
	cfg = cfg.WithMaxRetries(5)
	cfg = cfg.WithS3ForcePathStyle(true)
	cfg = cfg.WithHTTPClient(httpClient)
	cfg = cfg.WithDisableRestProtocolURICleaning(true) // do not delete first '/'

	cfg = setupCosLogging(cfg)
	return cfg
}

/*
Setting up the HTTP Client
*/
func NewHTTPClientWithSettings(httpSettings HTTPClientSettings) (*http.Client, error) {
	var httpClient http.Client
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			DualStack: true,
			Timeout:   httpSettings.Connect,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
		MaxConnsPerHost:       httpSettings.MaxConnsPerHost,
	}

	err := http2.ConfigureTransport(tr)
	if err != nil {
		return &httpClient, err
	}

	return &http.Client{
		Transport: tr,
	}, nil
}

/*
Setting the logging of HTTP requests in case of loglevel = DEBUG
*/
func setupCosLogging(cfg *aws.Config) *aws.Config {
	if config.BackintConfig.AgentLogLevelU() == "HTTP" {
		awsLogger := aws.LoggerFunc(func(args ...any) {
			logrus.WithField("time", time.Now().Format(time.RFC850)).Info(args...)
			logrus.SetOutput(logging.GetLogFile())
		})
		cfg = cfg.WithLogger(awsLogger)
		cfg = cfg.WithLogLevel(
			aws.LogDebugWithRequestErrors | aws.LogDebugWithRequestRetries,
		)
	}
	return cfg
}

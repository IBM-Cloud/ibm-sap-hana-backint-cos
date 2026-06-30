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
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/config"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/logging"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/client"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

// cosRetryer extends the SDK's DefaultRetryer to also retry TCP write timeouts.
// The default retryer treats "write: connection timed out" as non-retryable
// because net.OpError.Temporary() returns false for write errors. On Power VS
// a stateful firewall can silently drop a long-running PUT TCP connection,
// causing exactly this error mid-part. Retrying re-dials a fresh socket.
type cosRetryer struct {
	client.DefaultRetryer
}

func (r cosRetryer) ShouldRetry(req *request.Request) bool {
	if req.Error != nil && strings.Contains(req.Error.Error(), "write: connection timed out") {
		return true
	}
	return r.DefaultRetryer.ShouldRetry(req)
}

/*
Generating the session and the client to access the IBM Cloud Object Storage
*/
func GenerateCOSSession() (*session.Session, *s3.S3) {
	cfg := setupCosConfig()
	s3Session := session.Must(session.NewSession(cfg))
	s3Session.Config.Retryer = cosRetryer{
		DefaultRetryer: client.DefaultRetryer{NumMaxRetries: 5},
	}
	s3Client := s3.New(s3Session)
	return s3Session, s3Client
}

/*
Setting up the Cloud Object Storage Configuration
*/
func setupCosConfig() *aws.Config {
	httpClient, err := NewHTTPClientWithSettings(HTTPClientSettings{
		Connect:          60 * time.Second,
		ExpectContinue:   10 * time.Second, // metadata service can take up to 5s for token fetch
		IdleConn:         120 * time.Second,
		ConnKeepAlive:    30 * time.Second, // frequent probes to survive firewall during token refresh stall
		MaxHostIdleConns: 100,              // up to 12 HANA channels × 4 concurrency threads + headroom
		MaxAllIdleConns:  100,
		ResponseHeader:   0, // no timeout: response arrives only after full part body is sent
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
		authMethod = config.BackintConfig.AuthMode()

	} else {
		region = global.Args.Region
		endpoint = global.Args.EndpointUrl
		authEndpoint = global.Args.AuthEndpoint
		authMethod = global.Args.AuthMode
		if authMethod == config.AUTH_APIKEY {
			apikey, _ = global.ReadApikeyFromFile(global.Args.AuthKeypath)
		}
	}

	var creds *credentials.Credentials

	switch authMethod {
	case config.AUTH_APIKEY:
		creds = ibmiam.NewStaticCredentials(aws.NewConfig(),
			authEndpoint,
			apikey,
			"",
		)
	case config.AUTH_TRUSTEDPROFILE:
		creds = newPowerVSCredentials()
	default:
		break
	}

	cfg := aws.NewConfig()
	cfg = cfg.WithRegion(region)
	cfg = cfg.WithEndpoint(endpoint)
	cfg = cfg.WithCredentials(creds)
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
	dialer := &net.Dialer{
		Timeout:   httpSettings.Connect,
		KeepAlive: httpSettings.ConnKeepAlive,
		DualStack: true,
	}

	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		// Use a custom DialContext that forces TCP keepalive parameters at the
		// socket level. This ensures the OS sends keepalive probes during stalls
		// (e.g. waiting for the remote TCP window to open on large PUTs), which
		// prevents stateful firewalls on the IBM Cloud private network from
		// dropping the connection due to perceived inactivity.
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			if tc, ok := conn.(*net.TCPConn); ok {
				_ = tc.SetKeepAlive(true)
				_ = tc.SetKeepAlivePeriod(httpSettings.ConnKeepAlive)
			}
			return conn, nil
		},
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
		MaxConnsPerHost:       httpSettings.MaxConnsPerHost,
		// Disable HTTP/2: S3 uses HTTP/1.1. HTTP/2 would multiplex all parallel
		// part uploads over a single TCP connection — if that socket dies, all
		// concurrent uploads fail together.
		ForceAttemptHTTP2: false,
		TLSNextProto:      make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
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

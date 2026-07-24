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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	iamtoken "github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
)

// powerVSProvider implements credentials.Provider for PowerVS.
// It calls initFunc (→ metadata service) on every refresh and never
// attempts the IAM refresh_token flow, which would always fail because
// PowerVS-sourced IAM tokens carry no refresh token.
// IsExpired triggers a Retrieve 5 minutes before the token's real expiry
// so the COS SDK gets a fresh token before any S3 request can see a 401.
type powerVSProvider struct {
	mu         sync.Mutex
	initFunc   func() (*iamtoken.Token, error)
	expiration time.Time
}

func (p *powerVSProvider) Retrieve() (credentials.Value, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	tok, err := p.initFunc()
	if err != nil {
		return credentials.Value{}, err
	}

	if tok.Expiration > 0 {
		p.expiration = time.Unix(tok.Expiration, 0)
	} else {
		p.expiration = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	}

	return credentials.Value{
		Token: iamtoken.Token{
			AccessToken: tok.AccessToken,
			TokenType:   tok.TokenType,
			ExpiresIn:   tok.ExpiresIn,
			Expiration:  tok.Expiration,
		},
		ProviderName: "PowerVSProvider",
		ProviderType: "oauth",
	}, nil
}

// IsExpired returns true 5 minutes before the token's real expiry so the
// COS SDK's Credentials.Get() calls Retrieve() and fetches a new token
// from the metadata service before any in-flight S3 request can 401.
func (p *powerVSProvider) IsExpired() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return time.Now().After(p.expiration.Add(-5 * time.Minute))
}

/*
Build COS credentials for PowerVS authentication using the VPC Instance Metadata Service.
The VpcInstanceAuthenticator exchanges the instance identity token for an IAM access token.
powerVSProvider wraps the initFunc and refreshes directly from the metadata service on
every expiry — it never attempts the IAM refresh_token POST, which would always fail
because PowerVS-sourced tokens carry no refresh token.
*/
func newPowerVSCredentials(iamProfileId string, iamProfileName string) *credentials.Credentials {
	authenticatorBuilder := core.NewVpcInstanceAuthenticatorBuilder().
		SetURL(global.METADATA_SERVICE_URL).
		SetServiceVersion(global.SERVICE_VERSION)

	if iamProfileId != "" {
		authenticatorBuilder.SetIAMProfileID(iamProfileId)
	}

	if iamProfileName != "" {
		authenticatorBuilder.SetIAMProfileName(iamProfileName)
	}
	authenticator, err := authenticatorBuilder.Build()

	global.CheckForError(
		err,
		"Error creating PowerVS authenticator",
		global.FAILURE,
	)

	return credentials.NewCredentials(&powerVSProvider{
		initFunc: newPowerVSInitFunc(authenticator),
	})
}

// newPowerVSInitFunc returns an ibmiam token init function backed by the
// VpcInstanceAuthenticator. The COS SDK calls this on each credential refresh
// and the result is fed into the ibmiam oauth signer which sets the
// Authorization header correctly.
func newPowerVSInitFunc(authenticator *core.VpcInstanceAuthenticator) func() (*iamtoken.Token, error) {
	return func() (*iamtoken.Token, error) {
		tokenValue, err := authenticator.GetToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get PowerVS IAM token: %w", err)
		}

		var expiresIn int64 = 3600 // safe fallback: 1 hour
		var expiration int64

		// Decode the JWT payload to extract exp — no signature verification needed
		// since the token came directly from the trusted metadata service.
		parts := strings.Split(tokenValue, ".")
		if len(parts) == 3 {
			payload, decodeErr := base64.RawURLEncoding.DecodeString(parts[1])
			if decodeErr == nil {
				var claims struct {
					Iat int64 `json:"iat"`
					Exp int64 `json:"exp"`
				}
				if json.Unmarshal(payload, &claims) == nil && claims.Exp > 0 {
					expiration = claims.Exp
					if claims.Iat > 0 {
						expiresIn = claims.Exp - claims.Iat
					}
				}
			}
		}

		return &iamtoken.Token{
			AccessToken: tokenValue,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
			Expiration:  expiration,
		}, nil
	}
}

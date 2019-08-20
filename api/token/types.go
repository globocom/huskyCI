// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"time"

	"github.com/globocom/huskyCI/api/auth"
	"github.com/globocom/huskyCI/api/types"
)

// ExternalCalls defines a group of functions
// used for external calls and validate some
// necessary information about TokenHandler.
type ExternalCalls interface {
	ValidateURL(url string) (string, error)
	GenerateToken() (string, error)
	GetTimeNow() time.Time
	StoreAccessToken(accessToken types.DBToken) error
	FindAccessToken(id string) (types.DBToken, error)
	UpdateAccessToken(id string, accesstoken types.DBToken) error
	FindRepoURL(repositoryURL string) error
	GenerateUUID() string
	EncodeBase64(m string) string
	DecodeToStringBase64(encodedVal string) (string, error)
}

// THandler is a struct used to handle with
// token generation, validation and deactivation.
// It implements TokenInterface interface.
type THandler struct {
	External ExternalCalls
	HashGen  auth.Pbkdf2Generator
}

// TCaller implements ExternalCalls interface.
type TCaller struct{}

// TInterface is used to define functions that
// handle with access token management.
type TInterface interface {
	GenerateAccessToken(repo types.TokenRequest) (string, error)
	ValidateToken(token, repositoryURL string) error
	VerifyRepo(repositoryURL string) error
}

// TValidator is used to validate an access token
// using the defined functions TokenInterface
type TValidator struct {
	TokenVerifier TInterface
}

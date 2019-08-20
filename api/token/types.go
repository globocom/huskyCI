// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"time"

	"github.com/globocom/huskyCI/api/auth"
	"github.com/globocom/huskyCI/api/types"
)

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

type TokenHandler struct {
	External ExternalCalls
	HashGen  auth.Pbkdf2Generator
}

type TokenCaller struct{}

type TokenInterface interface {
	GenerateAccessToken(repo types.TokenRequest) (string, error)
	ValidateToken(token, repositoryURL string) error
	VerifyRepo(repositoryURL string) error
}

type TokenValidator struct {
	TokenVerifier TokenInterface
}

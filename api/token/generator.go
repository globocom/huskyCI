// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/google/uuid"
)

// ValidateURL validates if an URL is malicious or not.
func (tC *TokenCaller) ValidateURL(url string) (string, error) {
	return util.CheckMaliciousRepoURL(url)
}

func generateRandomBytes() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return b, err
}

// GenerateToken generates a new token to be used in auth.
func (tC *TokenCaller) GenerateToken() (string, error) {
	b, err := generateRandomBytes()
	return base64.URLEncoding.EncodeToString(b), err
}

// GetTimeNow returns the time now.
func (tC *TokenCaller) GetTimeNow() time.Time {
	return time.Now()
}

// StoreAccessToken stores a new access token into MongoDB.
func (tC *TokenCaller) StoreAccessToken(accessToken types.DBToken) error {
	return db.InsertDBAccessToken(accessToken)
}

// FindAccessToken gets an AccessToken based on an given ID.
func (tC *TokenCaller) FindAccessToken(ID string) (types.DBToken, error) {
	aTokenQuery := map[string]interface{}{"uuid": ID}
	return db.FindOneDBAccessToken(aTokenQuery)
}

// FindRepoURL checks if a Access TOken is present based on a given URL.
func (tC *TokenCaller) FindRepoURL(repositoryURL string) error {
	repoQuery := map[string]interface{}{"repositoryURL": repositoryURL}
	_, err := db.FindOneDBAccessToken(repoQuery)
	return err
}

// GenerateUUID returns a new UUID.
func (tC *TokenCaller) GenerateUUID() string {
	return uuid.New().String()
}

// EncodeBase64 retunrs a string in base64.
func (tC *TokenCaller) EncodeBase64(m string) string {
	return base64.URLEncoding.EncodeToString([]byte(m))
}

// DecodeToStringBase64 decodes a base64 string.
func (tC *TokenCaller) DecodeToStringBase64(encodedVal string) (string, error) {
	decodedVal, err := base64.URLEncoding.DecodeString(encodedVal)
	return string(decodedVal), err
}

// UpdateAccessToken updates an access Token in MongoDB based on its UUID.
func (tC *TokenCaller) UpdateAccessToken(ID string, accesstoken types.DBToken) error {
	aTokenQuery := map[string]interface{}{"uuid": ID}
	return db.UpdateOneDBAccessToken(aTokenQuery, accesstoken)
}

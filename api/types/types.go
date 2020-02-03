// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"time"
)

// TokenRequest defines the JSON struct for an access token request
type TokenRequest struct {
	RepositoryURL string `json:"repositoryURL"`
}

// AccessToken defines the struct generated when a new token
// is requested for specific repository
type AccessToken struct {
	HuskyToken string `bson:"huskytoken" json:"huskytoken"`
}

// DBToken defines the struct that stores husky access token
// for a repository URL
type DBToken struct {
	HuskyToken string    `bson:"huskytoken" json:"huskytoken"`
	URL        string    `bson:"repositoryURL" json:"repositoryURL"`
	IsValid    bool      `bson:"isValid" json:"isValid"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	Salt       string    `bson:"salt" json:"salt"`
	UUID       string    `bson:"uuid" json:"uuid"`
}

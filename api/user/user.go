// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"io"
	"os"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/globocom/huskyCI/api/auth"
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/types"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// DefaultAPIUser is the default API user from huskyCI
	DefaultAPIUser = os.Getenv("HUSKYCI_API_DEFAULT_USERNAME")
	// DefaultAPIPassword is the default API password from huskyCI
	DefaultAPIPassword = os.Getenv("HUSKYCI_API_DEFAULT_PASSWORD")
)

// Create generates a new user
func Create() types.User {
	newUser := types.User{}
	return newUser
}

// InsertDefaultUser insert default user into MongoDB
func InsertDefaultUser() error {

	var pbkdf2Caller auth.Pbkdf2Caller
	defaultHashFunction := pbkdf2Caller.GetHashName()
	keyLength := pbkdf2Caller.GetKeyLength()
	iterations := pbkdf2Caller.GetIterations()

	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}
	newUser := types.User{}
	newUser.Username = DefaultAPIUser
	newUser.HashFunction = defaultHashFunction
	newUser.Iterations = iterations
	newUser.KeyLen = keyLength
	newUser.Salt = base64.StdEncoding.EncodeToString(salt)
	hashedPass := pbkdf2.Key([]byte(DefaultAPIPassword), salt, iterations, keyLength, sha256.New)
	newUser.Password = base64.StdEncoding.EncodeToString(hashedPass)
	return apiContext.APIConfiguration.DBInstance.InsertDBUser(newUser)
}

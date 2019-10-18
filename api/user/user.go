// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"os"

	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash"
	"io"
	"strconv"

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
	// DefaultIterations is the default number of iterations used for auth
	DefaultIterations = os.Getenv("HUSKYCI_API_DEFAULT_ITERATIONS")
	// DefaultKeyLength is the default key length used for auth
	DefaultKeyLength = os.Getenv("HUSKYCI_API_DEFAULT_KEY_LENGTH")
	// DefaultHashFunction is the default hash function name of iterations used for auth
	DefaultHashFunction = os.Getenv("HUSKYCI_API_DEFAULT_HASH_FUNCTION")
)

// Create generates a new user
func Create() types.User {
	newUser := types.User{}
	return newUser
}

// InsertDefaultUser insert default user into MongoDB
func InsertDefaultUser() error {
	hashFunction, isValid := auth.GetValidHashFunction(DefaultHashFunction)
	if !isValid {
		return errors.New("Invalid hash function")
	}
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}
	keyLength, err := strconv.Atoi(DefaultKeyLength)
	if err != nil {
		return err
	}
	iterations, err := strconv.Atoi(DefaultIterations)
	if err != nil {
		return err
	}
	newUser := types.User{}
	newUser.Username = DefaultAPIUser
	newUser.HashFunction = DefaultHashFunction
	newUser.Iterations = iterations
	newUser.KeyLen = keyLength
	newUser.Salt = base64.StdEncoding.EncodeToString(salt)
	hashedPass := pbkdf2.Key([]byte(DefaultAPIPassword), salt, iterations, keyLength, func() hash.Hash {
		return hashFunction
	})
	newUser.Password = base64.StdEncoding.EncodeToString(hashedPass)
	return apiContext.APIConfiguration.DBInstance.InsertDBUser(newUser)
}

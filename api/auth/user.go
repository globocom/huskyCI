// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"os"

	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

var (
	// DefaultAPIUser is the default API user from huskyCI
	DefaultAPIUser = os.Getenv("HUSKYCI_API_DEFAULT_USERNAME")
	// DefaultAPIPassword is the default API password from huskyCI
	DefaultAPIPassword = os.Getenv("HUSKYCI_API_DEFAULT_PASSWORD")
)

// User is the struct that holds all data from a huskyCI API user
// old
type User struct {
	Username           string `bson:"username" json:"username"`
	Password           string `bson:"password" json:"password"`
	Salt               string `bson:"salt,omitempty" json:"salt"`
	Iterations         int    `bson:"iterations,omitempty" json:"iterations"`
	KeyLen             int    `bson:"keylen,omitempty" json:"keylen"`
	HashFunction       string `bson:"hashfunction,omitempty" json:"hashfunction"`
	NewPassword        string `bson:"newPassword,omitempty" json:"newPassword"`
	ConfirmNewPassword string `bson:"confirmNewPassword,omitempty" json:"confirmNewPassword"`
}

// Create generates a new user
func Create() User {
	newUser := User{}
	return newUser
}

// InsertDefaultUser insert default user into MongoDB
func InsertDefaultUser() error {

	var pbkdf2Caller Pbkdf2Caller
	defaultHashFunction := pbkdf2Caller.GetHashName()
	keyLength := pbkdf2Caller.GetKeyLength()
	iterations := pbkdf2Caller.GetIterations()

	hashFunction, isValid := GetValidHashFunction(defaultHashFunction)
	if !isValid {
		return errors.New("Invalid hash function")
	}
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}
	newUser := User{}
	newUser.Username = DefaultAPIUser
	newUser.HashFunction = defaultHashFunction
	newUser.Iterations = iterations
	newUser.KeyLen = keyLength
	newUser.Salt = base64.StdEncoding.EncodeToString(salt)
	hashedPass := pbkdf2.Key([]byte(DefaultAPIPassword), salt, iterations, keyLength, func() hash.Hash {
		return hashFunction
	})
	newUser.Password = base64.StdEncoding.EncodeToString(hashedPass)
	return nil
	// return apiContext.APIConfiguration.DBInstance.InsertDBUser(newUser)
}

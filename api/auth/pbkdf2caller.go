// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"hash"
	"io"
	"os"
	"strconv"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"golang.org/x/crypto/pbkdf2"
)

// GetCredsFromDB returns an user info given an username.
func (pC *Pbkdf2Caller) GetCredsFromDB(username string) (types.User, error) {
	searchParam := map[string]interface{}{"username": username}
	return db.FindOneDBUser(searchParam)
}

// DecodeSaltValue decodes a salt and returns a string and an error.
func (pC *Pbkdf2Caller) DecodeSaltValue(salt string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(salt)
}

// GenHashValue returns the hash value given all pbkdf2 parameters.
func (pC *Pbkdf2Caller) GenHashValue(value, salt []byte, iter, keyLen int, h hash.Hash) string {
	return base64.StdEncoding.EncodeToString(pbkdf2.Key(value, salt, iter, keyLen, func() hash.Hash {
		return h
	}))
}

// GenerateSalt returns a random salt and en error.
func (pC *Pbkdf2Caller) GenerateSalt() (string, error) {
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	return base64.StdEncoding.EncodeToString(salt), err
}

// GetHashName returns the default hash name that is stored in an env var.
func (pC *Pbkdf2Caller) GetHashName() string {
	if value, ok := os.LookupEnv("HUSKYCI_API_DEFAULT_HASH_FUNCTION"); ok {
		return value
	}
	return "SHA512"
}

// GetIterations returns the default number of iteration that is stored in an env var.
func (pC *Pbkdf2Caller) GetIterations() (int, error) {
	if value, ok := os.LookupEnv("HUSKYCI_API_DEFAULT_ITERATIONS"); ok {
		return strconv.Atoi(value)
	}
	return 100000, nil
}

// GetKeyLength returns the default key lenght that is stored in an env var.
func (pC *Pbkdf2Caller) GetKeyLength() (int, error) {
	if value, ok := os.LookupEnv("HUSKYCI_API_DEFAULT_KEY_LENGTH"); ok {
		return strconv.Atoi(value)
	}
	return 512, nil
}

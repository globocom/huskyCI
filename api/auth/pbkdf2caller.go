package auth

import (
	"encoding/base64"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"golang.org/x/crypto/pbkdf2"
	"hash"
)

type Pbkdf2Caller struct{}

func (pC *Pbkdf2Caller) GetCredsFromDB(username string) (types.User, error) {
	searchParam := map[string]interface{}{"username": username}
	return db.FindOneDBUser(searchParam)
}

func (pC *Pbkdf2Caller) DecodeSaltValue(salt string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(salt)
}

func (pC *Pbkdf2Caller) GenHashValue(value, salt []byte, iter, keyLen int, h hash.Hash) string {
	return base64.StdEncoding.EncodeToString(pbkdf2.Key(value, salt, iter, keyLen, func() hash.Hash {
		return h
	}))
}

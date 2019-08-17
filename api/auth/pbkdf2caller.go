package auth

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"io"
	"os"
	"strconv"
)

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

func (pC *Pbkdf2Caller) GenerateSalt() (string, error) {
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	return base64.StdEncoding.EncodeToString(salt), err
}

func (pC *Pbkdf2Caller) GetHashName() string {
	return os.Getenv("HUSKYCI_API_DEFAULT_HASH_FUNCTION")
}

func (pC *Pbkdf2Caller) GetIterations() (int, error) {
	return strconv.Atoi(os.Getenv("HUSKYCI_API_DEFAULT_ITERATIONS"))
}

func (pC *Pbkdf2Caller) GetKeyLength() (int, error) {
	return strconv.Atoi(os.Getenv("HUSKYCI_API_DEFAULT_KEY_LENGTH"))
}

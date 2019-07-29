package user

import (
	"os"

	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/globocom/huskyCI/api/auth"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"io"
	"strconv"
)

var (
	// DefaultAPIUser is the default API user from huskyCI
	DefaultAPIUser = os.Getenv("HUSKYCI_API_DEFAULT_USERNAME")
	// DefaultAPIPassword is the default API password from huskyCI
	DefaultAPIPassword  = os.Getenv("HUSKYCI_API_DEFAULT_PASSWORD")
	DefaultIterations   = os.Getenv("HUSKYCI_API_DEFAULT_ITERATIONS")
	DefaultKeyLength    = os.Getenv("HUSKYCI_API_DEFAULT_KEY_LENGTH")
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
	return db.InsertDBUser(newUser)
}

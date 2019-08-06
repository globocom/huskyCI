package auth

import (
	"github.com/globocom/huskyCI/api/types"
	"hash"
)

type UserCredsHandler interface {
	GetPassFromDB(username string) (string, error)
	GetHashedPass(password string) (string, error)
}

type Pbkdf2Generator interface {
	GetCredsFromDB(username string) (types.User, error)
	DecodeSaltValue(salt string) ([]byte, error)
	GenHashValue(value, salt []byte, iter, keyLen int, h hash.Hash) string
	GenerateSalt() (string, error)
	GetHashName() string
	GetIterations() (int, error)
	GetKeyLength() (int, error)
}

type Pbkdf2Caller struct{}

type MongoBasic struct {
	ClientHandler UserCredsHandler
}

type ClientPbkdf2 struct {
	HashGen      Pbkdf2Generator
	Salt         string
	Iterations   int
	KeyLen       int
	HashFunction string
}

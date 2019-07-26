package auth

import (
	"hash"
)

type UserCredsHandler interface {
	GetPassFromDB(username string) (string, error)
	GetHashedPass(password string) (string, error)
}

type Pbkdf2Generator interface {
	GenHashValue(value, salt []byte, iter, keyLen int, h hash.Hash) []byte
}

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

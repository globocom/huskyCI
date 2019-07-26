package auth

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"golang.org/x/crypto/sha3"
	"hash"
	"os"
	"strings"
)

// GetPassFromDB will search for a valid user entry on DB through the
// received username. It will set all parameters required for PBKDF's
// hash generation and return the hash password stored.
func (cM *ClientPbkdf2) GetPassFromDB(username string) (string, error) {
	// TODO
	return "", nil
}

// GetHashedPass will return the hash value of given password based
// on the parameters set by GetPassFromDB. It will verify first if
// all parameters required are valid.
func (cM *ClientPbkdf2) GetHashedPass(password string) (string, error) {
	validHashes := os.Getenv("HUSKY_PBKDF2_VALID_HASHES")
	hashes := strings.SplitN(validHashes, ",", -1)
	isValid := false
	for _, hashAlg := range hashes {
		if strings.EqualFold(hashAlg, cM.HashFunction) {
			isValid = true
			break
		}
	}
	if cM.Salt == "" || cM.Iterations == 0 || cM.KeyLen == 0 || !isValid {
		return "", errors.New("Failed to generate a hash! It doesn't meet all criteria")
	}
	var hashFunction hash.Hash
	switch cM.HashFunction {
	case "sha224":
		hashFunction = sha256.New224()
	case "sha384":
		hashFunction = sha512.New384()
	case "sha512":
		hashFunction = sha512.New()
	case "sha3_224":
		hashFunction = sha3.New224()
	case "sha3_256":
		hashFunction = sha3.New256()
	case "sha3_384":
		hashFunction = sha3.New384()
	case "sha3_512":
		hashFunction = sha3.New512()
	default:
		hashFunction = sha256.New()
	}
	return string(cM.HashGen.GenHashValue([]byte(password), []byte(cM.Salt), cM.Iterations, cM.KeyLen, hashFunction)), nil
}

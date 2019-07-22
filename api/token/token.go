package token

import (
	"fmt"

	"github.com/globocom/huskyCI/api/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateHuskyCIToken creates a new huskyCI token and returns it
func GenerateHuskyCIToken() (types.HuskyCIToken, error) {

	huskyCIToken := types.HuskyCIToken{}

	newToken := generateID()
	hashedNewToken, err := bcrpytToken(newToken)
	if err != nil {
		return huskyCIToken, err
	}

	huskyCIToken.ID = newToken
	huskyCIToken.HashedToken = hashedNewToken

	return huskyCIToken, nil
}

// bcrpytToken returns a hashed token using bcrypt and an error.
func bcrpytToken(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	return string(bytes), err
}

// checkTokenHash returns a bool if a token matches its brcrypt hash.
func checkTokenHash(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

// generateID returns a UUID.
func generateID() string {
	u1 := uuid.New()
	userID := fmt.Sprintf("%s", u1)
	return userID
}

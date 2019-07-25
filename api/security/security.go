package security

import "golang.org/x/crypto/bcrypt"

// BcryptPassword returns a hashed token using bcrypt and an error.
func BcryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash returns a bool if a token matches its brcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

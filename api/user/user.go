package user

import (
	"os"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/security"
	"github.com/globocom/huskyCI/api/types"
)

var (
	// DefaultAPIUser is the default API user from huskyCI
	DefaultAPIUser = os.Getenv("HUSKYCI_API_DEFAULT_USERNAME")
	// DefaultAPIPassword is the default API password from huskyCI
	DefaultAPIPassword = os.Getenv("HUSKYCI_API_DEFAULT_PASSWORD")
)

// Create generates a new user
func Create() types.User {
	newUser := types.User{}
	return newUser
}

// InserDefaultUser insert default user into MongoDB
func InserDefaultUser() error {
	newUser := types.User{}

	hashedPassword, err := security.BcryptPassword(DefaultAPIPassword)
	if err != nil {
		return err
	}

	newUser.Name = DefaultAPIUser
	newUser.HashedPassword = hashedPassword

	return db.InsertDBUser(newUser)
}

package auth

import (
	"github.com/labstack/echo"
)

// ValidateUser is called by the echo's middleware for
// basic auth validation
func ValidateUser(username, password string, c echo.Context) (bool, error) {
	clientMongo := ClientPbkdf2{
		HashGen: &Pbkdf2Caller{},
	}
	basicClient := MongoBasic{
		ClientHandler: &clientMongo,
	}
	return basicClient.IsValidUser(username, password)
}

// IsValidUser will verify if it has a valid user for the username passed
// and validate password through its hash value compared with the
// hash value in stored
func (mB MongoBasic) IsValidUser(username, password string) (bool, error) {
	passDB, err := mB.ClientHandler.GetPassFromDB(username)
	if err != nil {
		return false, nil
	}
	hashedPass, err := mB.ClientHandler.GetHashedPass(password)
	if err != nil {
		return false, err
	}
	if passDB != hashedPass {
		return false, nil
	}
	return true, nil
}

package auth

import (
	"github.com/labstack/echo"
)

func ValidateUser(username, password string, c echo.Context) (bool, error) {
	basicClient := MongoBasic{}
	return basicClient.IsValidUser(username, password)
}

func (mB MongoBasic) IsValidUser(username, password string) (bool, error) {
	passDB, err := mB.ClientHandler.GetPassFromDB(username)
	if err != nil {
		return false, err
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

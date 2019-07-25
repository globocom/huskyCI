package auth

import (
	"github.com/labstack/echo"
)

func (bA BasicAuthentication) ValidateUser(username, password string, c echo.Context) (bool, error) {
	isValid, err := bA.BasicClient.IsValidUser(username, password)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}
	return true, nil
}

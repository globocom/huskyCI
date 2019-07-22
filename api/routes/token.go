package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/token"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/labstack/echo"
)

// CreateNewToken generates a new token
func CreateNewToken(c echo.Context) error {

	// valid JSON?
	attemptNewToken := types.HuskyCIToken{}
	err := c.Bind(&attemptNewToken)
	if err != nil {
		log.Warning("CreateNewToken", "TOKEN", 1024, err)
		reply := map[string]interface{}{"success": false, "error": "invalid token JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// valid Repositories?
	for _, repository := range attemptNewToken.Repositories {
		if _, err := util.CheckMaliciousRepoURL(repository, c); err != nil {
			return err
		}
	}

	// auhtorized? will check this later

	// generate new token
	newHuskyCIToken, err := token.GenerateHuskyCIToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newHuskyCIToken)
	}
	newHuskyCIToken.Repositories = attemptNewToken.Repositories

	// register in MongoDB

	// return to user
	userOutputToken := types.HuskyCIToken{
		ID: newHuskyCIToken.ID,
	}

	return c.JSON(http.StatusCreated, userOutputToken)
}

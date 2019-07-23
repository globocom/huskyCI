package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/token"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/labstack/echo"
)

// CreateNewToken generates a new token
func CreateNewToken(c echo.Context) error {

	// step 1: valid JSON?
	attemptNewToken := types.HuskyCIToken{}
	err := c.Bind(&attemptNewToken)
	if err != nil {
		log.Warning("CreateNewToken", "TOKEN", 1024, err)
		reply := map[string]interface{}{"success": false, "error": "invalid token JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// step 2: valid Repositories?
	if attemptNewToken.Repositories == nil {
		reply := map[string]interface{}{"success": false, "error": "invalid input"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	for _, repository := range attemptNewToken.Repositories {
		if _, err := util.CheckMaliciousRepoURL(repository, c); err != nil {
			reply := map[string]interface{}{"success": false, "error": "invalid input"}
			return c.JSON(http.StatusBadRequest, reply)
		}
	}

	// step 3: auhtorized? will check this later

	// step 4: generate new token
	newHuskyCIToken, err := token.GenerateHuskyCIToken()
	if err != nil {
		reply := map[string]interface{}{"success": false, "error": "internal error generating token"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// step 5: register in MongoDB
	newHuskyCIToken.Repositories = attemptNewToken.Repositories
	if err := db.InsertDBToken(newHuskyCIToken); err != nil {
		log.Error("CreateNewToken", "TOKEN", 1010, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// step 6: return to user only huskyCIToken.ID
	userOutputToken := types.HuskyCIToken{
		ID:           newHuskyCIToken.ID,
		Repositories: newHuskyCIToken.Repositories,
	}
	return c.JSON(http.StatusCreated, userOutputToken)
}

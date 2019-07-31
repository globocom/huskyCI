package routes

import (
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/token"
	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
	"net/http"
)

// HandleToken generate an access token for a specific repository
func HandleToken(c echo.Context) error {
	repoRequest := types.TokenRequest{}
	if err := c.Bind(&repoRequest); err != nil {
		log.Error("HandleToken", "TOKEN", 1025, err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "invalid token JSON"})
	}
	tokenCaller := token.TokenCaller{}
	tokenHandler := token.TokenHandler{
		External: &tokenCaller,
	}
	accessToken, err := tokenHandler.GenerateAccessToken(repoRequest)
	if err != nil {
		log.Error("HandleToken ", "TOKEN", 10216, err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "token generation failure"})
	}
	return c.JSON(http.StatusCreated, accessToken)
}

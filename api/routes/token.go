// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/auth"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/token"
	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
)

var (
	tokenHandler token.THandler
)

func init() {
	tokenCaller := token.TCaller{}
	hashGen := auth.Pbkdf2Caller{}
	tokenHandler = token.THandler{
		External: &tokenCaller,
		HashGen:  &hashGen,
	}
}

// HandleToken generate an access token for a specific repository
func HandleToken(c echo.Context) error {
	repoRequest := types.TokenRequest{}
	if err := c.Bind(&repoRequest); err != nil {
		log.Error("HandleToken", "TOKEN", 1025, err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "invalid token JSON"})
	}
	log.Info("HandleToken", "TOKEN", 24, repoRequest.RepositoryURL)
	accessToken, err := tokenHandler.GenerateAccessToken(repoRequest)
	if err != nil {
		log.Error("HandleToken ", "TOKEN", 1026, err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "token generation failure"})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"huskytoken": accessToken})
}

// HandleDeactivation will deactivate an access token passed in the body
// of the request
func HandleDeactivation(c echo.Context) error {
	tokenRequest := types.AccessToken{}
	if err := c.Bind(&tokenRequest); err != nil {
		log.Error("HandleInvalidate", "TOKEN", 1025, err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "invalid token JSON"})
	}
	if err := tokenHandler.InvalidateToken(tokenRequest.HuskyToken); err != nil {
		log.Error("HandleInvalidate ", "TOKEN", 1028, err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "token deactivation failure"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "error": ""})
}

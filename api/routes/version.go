// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/labstack/echo"
)

//GetAPIVersion returns the API version
func GetAPIVersion(c echo.Context) error {
	configAPI := apiContext.APIConfiguration
	return c.JSON(http.StatusOK, GetRequestResult(configAPI))
}

// GetRequestResult returns a map containing API's version and release date
func GetRequestResult(configAPI *apiContext.APIConfig) map[string]string {
	requestResult := map[string]string{
		"version": configAPI.Version,
		"date":    configAPI.ReleaseDate,
	}
	return requestResult
}

package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
)

// Version holds the API version to be returned in /version route.
var Version types.VersionAPI

//GetAPIVersion returns the API version
func GetAPIVersion(c echo.Context) error {
	requestResult := map[string]string{"version": Version.Version, "date": Version.Date}
	return c.JSON(http.StatusOK, requestResult)
}

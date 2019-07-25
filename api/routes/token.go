package routes

import (
	"github.com/labstack/echo"
	"net/http"
)

// Generate an access token for a specific repository
func HandleToken(c echo.Context) error {
	// TODO
	return c.String(http.StatusOK, "")
}

package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// HealthCheck is the heath check function.
func (es *EchoServer) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING\n")
}

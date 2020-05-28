package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// HealthCheck is the heath check function.
func (es *EchoServer) HealthCheck(c echo.Context) error {

	cOutput, err := es.RunnerSession.Run("huskyci/enry:latest", "whoami")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, cOutput)
}

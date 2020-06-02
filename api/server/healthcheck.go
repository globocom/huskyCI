package server

import (
	"net/http"

	"github.com/globocom/huskyCI/api/analysis"
	"github.com/labstack/echo"
)

// HealthCheck is the heath check function.
func (es *EchoServer) HealthCheck(c echo.Context) error {

	newAnalysis := analysis.New("https://github.com/globocom/huskyCI.git", "master")
	err := es.DatabaseSession.InsertAnalysis(newAnalysis)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "FOI?!")
}

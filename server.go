package main

import (
	"github.com/globocom/husky/analysis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	echoInstance := echo.New()

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", analysis.HealthCheck)
	echoInstance.GET("/husky/:id", analysis.StatusAnalysis)
	echoInstance.POST("/husky", analysis.StartAnalysis)
	echoInstance.POST("/securitytest", analysis.CreateNewSecurityTest)

	echoInstance.Logger.Fatal(echoInstance.Start(":9999"))

}

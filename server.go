package main

import (
	"github.com/globocom/husky/analysis"
	"github.com/globocom/husky/config"
	configMiddleware "github.com/globocom/husky/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	confEnv := new(config.Config)

	echoInstance := echo.New()

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())
	echoInstance.Use(configMiddleware.RequestConfigMiddleware(confEnv))

	echoInstance.GET("/healthcheck", analysis.HealthCheck)
	//echoInstance.GET("/analyze/:id", analysis.StatusAnalysis)
	echoInstance.POST("/analyze", analysis.StartAnalysis)

	echoInstance.Logger.Fatal(echoInstance.Start(":9999"))

}

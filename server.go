package main

import (
	"fmt"

	"github.com/globocom/husky/analysis"
	"github.com/globocom/husky/types"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	err := checkAndInitMongo()
	if err != nil {
		fmt.Println("Check MongoDB. Something went wrong:", err)
	}

	echoInstance := echo.New()

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", analysis.HealthCheck)
	//echoInstance.GET("/analyze/:id", analysis.StatusAnalysis)
	echoInstance.POST("/analyze", analysis.StartAnalysis)

	echoInstance.Logger.Fatal(echoInstance.Start(":9999"))

}

// checkAndInitMongo will check and initiate SecurityTestCollecion
func checkAndInitMongo() error {
	s := types.SecurityTest{Name: "enry"}
	_, err := analysis.CheckSecurityTest(s)
	if err != nil {
		fmt.Println("First time running Husky? Error:", err)
		err = analysis.InitSecurityTestCollection()
		if err != nil {
			fmt.Println("Could not initiate SecurityTestCollection. Is MongoDB running? Error:", err)
		}
	}
	return err
}

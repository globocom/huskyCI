// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/globocom/huskyCI/api/auth"
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/routes"
	"github.com/globocom/huskyCI/api/util"
	apiUtil "github.com/globocom/huskyCI/api/util/api"
)

func main() {

	configAPI, err := apiContext.DefaultConf.GetAPIConfig()

	if err != nil {
		fmt.Println("Error in configuration file: ", err)
		os.Exit(1)
	}

	log.InitLog(
		configAPI.GraylogConfig.DevelopmentEnv,
		configAPI.GraylogConfig.Address,
		configAPI.GraylogConfig.Protocol,
		configAPI.GraylogConfig.AppName,
		configAPI.GraylogConfig.Tag)
	log.Info("main", "SERVER", 11)

	checkHandler := &apiUtil.CheckUtils{}

	huskyUtils := apiUtil.HuskyUtils{
		CheckHandler: checkHandler,
	}

	if err := huskyUtils.CheckHuskyRequirements(configAPI); err != nil {
		log.Error("main", "SERVER", 1001, err)
		os.Exit(1)
	}

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{configAPI.AllowOriginValue},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// set new object for /api/1.0 route
	g := echoInstance.Group("/api/1.0")

	// use basic auth middleware
	g.Use(middleware.BasicAuth(auth.ValidateUser))

	// /token route with basic auth
	g.POST("/token", routes.HandleToken)
	g.POST("/token/deactivate", routes.HandleDeactivation)

	// generic routes
	echoInstance.GET("/healthcheck", routes.HealthCheck)
	echoInstance.GET("/version", routes.GetAPIVersion)

	// analysis routes
	echoInstance.POST("/analysis", routes.ReceiveRequest)
	echoInstance.GET("/analysis/:id", routes.GetAnalysis)
	// echoInstance.PUT("/analysis/:id", routes.UpdateAnalysis)
	// echoInstance.DELETE("/analysis/:id", routes.DeleteAnalysis)

	// stats routes
	echoInstance.GET("/stats/:metric_type", routes.GetMetric)

	// securityTest routes
	// echoInstance.GET("securityTest/:securityTestName", routes.GetSecurityTest)
	// echoInstance.POST("/securitytest", routes.CreateNewSecurityTest)
	// echoInstance.PUT("/securityTest/:securityTestName", routes.UpdateSecurityTest)
	// echoInstance.DELETE("/securityTest/:securityTestName", routes.DeleteSecurityTest)

	// repository routes
	// echoInstance.GET("/repository/:repoID", routes.GetRepository)
	// echoInstance.POST("/repository", routes.CreateNewRepository)
	// echoInstance.PUT("/repository/:repoID)
	// echoInstance.DELETE("/repository/:repoID)

	// user routes
	// echoInstance.GET("/user", routes.GetUser)
	// echoInstance.POST("/user", routes.CreateNewUser)
	echoInstance.PUT("/user", routes.UpdateUser)
	// echoInstance.DELETE("/user)

	huskyAPIport := fmt.Sprintf(":%d", configAPI.Port)

	if !configAPI.UseTLS {
		echoInstance.Logger.Fatal(echoInstance.Start(huskyAPIport))
	} else {
		echoInstance.Logger.Fatal(echoInstance.StartTLS(huskyAPIport, util.CertFile, util.KeyFile))
	}
}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/routes"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	apiUtil "github.com/globocom/huskyCI/api/util/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	version string
	commit  string
	date    string
)

const projectName = "HuskyCI"

func main() {

	log.InitLog()
	log.Info("main", "SERVER", 11)

	configAndPrintVersion(version, commit, date)
	configAPI := apiContext.GetAPIConfig()

	if err := apiUtil.CheckHuskyRequirements(configAPI); err != nil {
		log.Error("main", "SERVER", 1001, err)
		os.Exit(1)
	}

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	// generic routes
	echoInstance.GET("/healthcheck", routes.HealthCheck)
	echoInstance.GET("/version", routes.GetAPIVersion)

	// analysis routes
	echoInstance.GET("/analysis/:id", routes.GetAnalysis)
	echoInstance.POST("/analysis", routes.ReceiveRequest)
	// echoInstance.PUT("/analysis/:id", routes.UpdateAnalysis)
	// echoInstance.DELETE("/analysis/:id", routes.DeleteAnalysis)

	// securityTest routes
	// echoInstance.GET("securityTest/:securityTestName", routes.GetSecurityTest)
	echoInstance.POST("/securitytest", routes.CreateNewSecurityTest)
	// echoInstance.PUT("/securityTest/:securityTestName", routes.UpdateSecurityTest)
	// echoInstance.DELETE("/securityTest/:securityTestName", routes.DeleteSecurityTest)

	// repository routes
	// echoInstance.GET("/repository/:repoID", routes.GetRepository)
	echoInstance.POST("/repository", routes.CreateNewRepository)
	// echoInstance.PUT("/repository/:repoID)
	// echoInstance.DELETE("/repository/:repoID)

	huskyAPIport := fmt.Sprintf(":%d", configAPI.HuskyAPIPort)

	if !configAPI.UseTLS {
		echoInstance.Logger.Fatal(echoInstance.Start(huskyAPIport))
	} else {
		echoInstance.Logger.Fatal(echoInstance.StartTLS(huskyAPIport, util.CertFile, util.KeyFile))
	}
}

func configAndPrintVersion(version, commit, date string) {

	routes.Version.Project = projectName
	routes.Version.Version = version
	routes.Version.Commit = commit
	routes.Version.Date = date

	printVersion(routes.Version)
}

func printVersion(versionAPI types.VersionAPI) {
	vFlag := flag.Bool("v", false, "print current version")
	versionFlag := flag.Bool("version", false, "print current version")
	flag.Parse()

	if *vFlag || *versionFlag {
		versionAPI.Print()
		os.Exit(0)
	}

	if versionAPI.Version != "" {
		versionAPI.Print()
	}
}

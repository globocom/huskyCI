package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/globocom/glbgelf"
	"github.com/globocom/huskyci/analysis"
	apiContext "github.com/globocom/huskyci/context"
	db "github.com/globocom/huskyci/db/mongo"
	docker "github.com/globocom/huskyci/dockers"
	"github.com/globocom/huskyci/types"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	mgo "gopkg.in/mgo.v2"
)

var (
	version string
	commit  string
	date    string
)

const projectName = "HuskyCI"

func main() {

	configAndPrintVersion(version, commit, date)

	isDev := true
	if strings.EqualFold(os.Getenv("HUSKYCI_DEV"), "false") {
		isDev = false
	}

	glbgelf.InitLogger(os.Getenv("HUSKYCI_GRAYLOG_ADDR"), os.Getenv("HUSKYCI_APP_NAME"), os.Getenv("HUSKYCI_TAGS"), isDev, os.Getenv("HUSKYCI_GRAYLOG_PROTO"))

	if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "main",
		"info":   "SERVER"}, "INFO", "Starting Husky..."); errLog != nil {
		fmt.Println("glbgelf error: ", errLog)
	}

	configAPI := apiContext.GetAPIConfig()

	if err := checkHuskyRequirements(configAPI); err != nil {
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "main",
			"info":   "SERVER"}, "ERROR", "Error starting Husky:", err); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		os.Exit(1)
	}

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", analysis.HealthCheck)
	echoInstance.GET("/version", analysis.VersionHandler)
	echoInstance.GET("/husky/:id", analysis.StatusAnalysis)
	echoInstance.POST("/husky", analysis.ReceiveRequest)
	echoInstance.POST("/securitytest", analysis.CreateNewSecurityTest)
	echoInstance.POST("/repository", analysis.CreateNewRepository)

	huskyAPIport := fmt.Sprintf(":%d", configAPI.HuskyAPIPort)
	echoInstance.Logger.Fatal(echoInstance.Start(huskyAPIport))
}

func checkHuskyRequirements(configAPI *apiContext.APIConfig) error {

	// check if all environment variables are properly set.
	if err := checkEnvVars(); err != nil {
		return err
	}

	if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "checkHuskyRequirements",
		"info":   "SERVER"}, "INFO", "Environment Variables: OK!"); errLog != nil {
		fmt.Println("glbgelf error: ", errLog)
	}

	// check if all docker hosts are up and running docker API.
	if err := checkDockerHosts(configAPI); err != nil {
		return err
	}

	if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "checkHuskyRequirements",
		"info":   "SERVER"}, "INFO", "Docker API Hosts: OK!"); errLog != nil {
		fmt.Println("glbgelf error: ", errLog)
	}

	// check if MongoDB is acessible and credentials received are working.
	if err := checkMongoDB(); err != nil {
		return err
	}

	if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "checkHuskyRequirements",
		"info":   "SERVER"}, "INFO", "MongoDB: OK!"); errLog != nil {
		fmt.Println("glbgelf error: ", errLog)
	}

	// check if default securityTests are set into MongoDB.
	if err := checkDefaultSecurityTests(configAPI); err != nil {
		return err
	}

	if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "checkHuskyRequirements",
		"info":   "SERVER"}, "INFO", "Default security tests set: OK!"); errLog != nil {
		fmt.Println("glbgelf error: ", errLog)
	}

	return nil
}

func checkEnvVars() error {

	envVars := []string{
		"DOCKER_HOSTS_LIST",
		"MONGO_HOST",
		"MONGO_DATABASE_NAME",
		"MONGO_DATABASE_USERNAME",
		"MONGO_DATABASE_PASSWORD",
		"DOCKER_HOSTS_CERT",
		"DOCKER_HOSTS_KEY",
		// "GIT_PRIVATE_SSH_KEY", optional
		// "DOCKER_API_PORT", optional -> default value (2376)
		// "MONGO_PORT", optional -> default value (27017)
		// "HUSKY_API_PORT", optional -> default value (9999)
		// "MONGO_TIMEOUT", optional -> default value (60s)
	}

	var envIsSet bool
	var allEnvIsSet bool
	var errorString string

	env := make(map[string]string)
	allEnvIsSet = true
	for i := 0; i < len(envVars); i++ {
		env[envVars[i]], envIsSet = os.LookupEnv(envVars[i])
		if !envIsSet {
			errorString = errorString + envVars[i] + " "
			allEnvIsSet = false
		}
	}

	if allEnvIsSet == false {
		finalError := fmt.Sprintf("check environment variables: %s", errorString)
		return errors.New(finalError)
	}

	return nil
}

func checkDockerHosts(configAPI *apiContext.APIConfig) error {
	return docker.HealthCheckDockerAPI()
}

func checkMongoDB() error {

	_, err := db.Connect()

	if err != nil {
		mongoError := fmt.Sprintf("check mongoDB: %s", err)
		return errors.New(mongoError)
	}

	return nil
}

func checkDefaultSecurityTests(configAPI *apiContext.APIConfig) error {
	enryQuery := map[string]interface{}{"name": "enry"}
	enry, err := analysis.FindOneDBSecurityTest(enryQuery)
	if err == mgo.ErrNotFound {
		// As Enry securityTest is not set into MongoDB, HuskyCI will insert it.
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "checkDefaultSecurityTests",
			"info":   "SERVER"}, "ERROR", "Enry securityTest not found!"); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		enry = *configAPI.EnrySecurityTest
		if err := analysis.InsertDBSecurityTest(enry); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	gosecQuery := map[string]interface{}{"name": "gosec"}
	gosec, err := analysis.FindOneDBSecurityTest(gosecQuery)
	if err == mgo.ErrNotFound {
		// As Gosec securityTest is not set into MongoDB, HuskyCI will insert it.
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "checkDefaultSecurityTests",
			"info":   "SERVER"}, "ERROR", "Gosec securityTest not found!"); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		gosec = *configAPI.GosecSecurityTest
		if err := analysis.InsertDBSecurityTest(gosec); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	brakemanQuery := map[string]interface{}{"name": "brakeman"}
	brakeman, err := analysis.FindOneDBSecurityTest(brakemanQuery)
	if err == mgo.ErrNotFound {
		// As Brakeman securityTest is not set into MongoDB, HuskyCI will insert it.
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "checkDefaultSecurityTests",
			"info":   "SERVER"}, "ERROR", "Brakeman securityTest not found!"); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		brakeman = *configAPI.BrakemanSecurityTest
		if err := analysis.InsertDBSecurityTest(brakeman); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	banditQuery := map[string]interface{}{"name": "bandit"}
	bandit, err := analysis.FindOneDBSecurityTest(banditQuery)
	if err == mgo.ErrNotFound {
		// As Bandit securityTest is not set into MongoDB, HuskyCI will insert it.
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "checkDefaultSecurityTests",
			"info":   "SERVER"}, "ERROR", "Bandit securityTest not found!"); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		bandit = *configAPI.BanditSecurityTest
		if err := analysis.InsertDBSecurityTest(bandit); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	retirejsQuery := map[string]interface{}{"name": "retirejs"}
	retirejs, err := analysis.FindOneDBSecurityTest(retirejsQuery)
	if err == mgo.ErrNotFound {
		// As RetireJS securityTest is not set into MongoDB, HuskyCI will insert it.
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "checkDefaultSecurityTests",
			"info":   "SERVER"}, "ERROR", "RetireJS securityTest not found!"); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		retirejs = *configAPI.RetirejsSecurityTest
		if err := analysis.InsertDBSecurityTest(retirejs); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func configAndPrintVersion(version, commit, date string) {

	analysis.Version.Project = projectName
	analysis.Version.Version = version
	analysis.Version.Commit = commit
	analysis.Version.Date = date

	printVersion(analysis.Version)
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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/globocom/husky/analysis"
	apiContext "github.com/globocom/husky/context"
	db "github.com/globocom/husky/db/mongo"
	docker "github.com/globocom/husky/dockers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	mgo "gopkg.in/mgo.v2"
)

var (
	version string
	commit  string
	author  string
	date    string
)

const ProjectName = "HuskyCI"

func main() {

	configVersion(version, commit, author, date)

	fmt.Println("[*] Starting Husky...")

	configAPI := apiContext.GetAPIConfig()

	if err := checkHuskyRequirements(configAPI); err != nil {
		fmt.Println("[x] Error starting Husky:")
		fmt.Println("[x]", err)
		os.Exit(1)
	}

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", analysis.HealthCheck)
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

	fmt.Println("[*] Environment Variables: OK!")

	// check if all docker hosts are up and running docker API.
	if err := checkDockerHosts(configAPI); err != nil {
		return err
	}

	fmt.Println("[*] Docker API Hosts: OK!")

	// check if MongoDB is acessible and credentials received are working.
	if err := checkMongoDB(); err != nil {
		return err
	}

	fmt.Println("[*] MongoDB: OK!")

	// check if default securityTests are set into MongoDB.
	if err := checkDefaultSecurityTests(configAPI); err != nil {
		return err
	}

	fmt.Println("[*] Default security tests set: OK!")

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
		fmt.Println("[!] Enry securityTest not found!")
		enry = *configAPI.EnrySecurityTest
		if err := analysis.InsertDBSecurityTest(enry); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	gasQuery := map[string]interface{}{"name": "gas"}
	gas, err := analysis.FindOneDBSecurityTest(gasQuery)
	if err == mgo.ErrNotFound {
		// As Gas securityTest is not set into MongoDB, HuskyCI will insert it.
		fmt.Println("[!] Gas securityTest not found!")
		gas = *configAPI.GasSecurityTest
		if err := analysis.InsertDBSecurityTest(gas); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	brakemanQuery := map[string]interface{}{"name": "brakeman"}
	brakeman, err := analysis.FindOneDBSecurityTest(brakemanQuery)
	if err == mgo.ErrNotFound {
		// As Brakeman securityTest is not set into MongoDB, HuskyCI will insert it.
		fmt.Println("[!] Brakeman securityTest not found!")
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
		fmt.Println("[!] Bandit securityTest not found!")
		bandit = *configAPI.BanditSecurityTest
		if err := analysis.InsertDBSecurityTest(bandit); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func configVersion(version, commit, author, date string) {
	analysis.Version.Version = version
	analysis.Version.Commit = commit
	analysis.Version.Author = author
	analysis.Version.Date = date

	printVersion(version, commit, author, date)
}

func printVersion(version, commit, author, date string) {
	vFlag := flag.Bool("v", false, "print current version")
	versionFlag := flag.Bool("version", false, "print current version")
	flag.Parse()

	printVersion := fmt.Sprintf(`
	************************************
	project: %s
	version: %s
	commit: %s
	author: %s
	data build: %s
	************************************
	`, ProjectName, version, commit, author, date)

	if *vFlag || *versionFlag {
		fmt.Println(printVersion)
		os.Exit(0)
	}

	if version != "" && date != "" && author != "" {
		fmt.Println(printVersion)
	}
}

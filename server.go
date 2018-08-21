package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/globocom/husky/analysis"
	apiContext "github.com/globocom/husky/context"
	db "github.com/globocom/husky/db/mongo"
	docker "github.com/globocom/husky/dockers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	fmt.Println("[*] Starting Husky...")

	apiConfig := apiContext.GetAPIConfig()

	if err := checkHuskyRequirements(apiConfig); err != nil {
		fmt.Println("[x] Error starting Husky:")
		fmt.Println("[x]", err)
		os.Exit(1)
	}

	echoInstance := echo.New()

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", analysis.HealthCheck)
	echoInstance.GET("/husky/:id", analysis.StatusAnalysis)
	echoInstance.POST("/husky", analysis.ReceiveRequest)
	echoInstance.POST("/securitytest", analysis.CreateNewSecurityTest)
	echoInstance.POST("/repository", analysis.CreateNewRepository)

	huskyAPIport := fmt.Sprintf(":%d", apiConfig.HuskyAPIPort)
	echoInstance.Logger.Fatal(echoInstance.Start(huskyAPIport))
}

func checkHuskyRequirements(apiConfig *apiContext.APIConfig) error {

	// check if all environment variables are properly set.
	if err := checkEnvVars(apiConfig); err != nil {
		return err
	} else {
		fmt.Println("[*] Environment Variables: OK!")
	}

	// check if all docker hosts are up and running docker API.
	if err := checkDockerHosts(apiConfig); err != nil {
		return err
	} else {
		fmt.Println("[*] Docker API Hosts: OK!")
	}

	// check if MongoDB is acessible and credentials received are working.
	if err := checkMongoDB(apiConfig); err != nil {
		return err
	} else {
		fmt.Println("[*] MongoDB: OK!")
	}

	return nil
}

func checkEnvVars(apiConfig *apiContext.APIConfig) error {

	envVars := []string{
		"DOCKER_HOSTS_LIST",
		"MONGO_HOST",
		"MONGO_DATABASE_NAME",
		"MONGO_DATABASE_USERNAME",
		"MONGO_DATABASE_PASSWORD",
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

func checkDockerHosts(apiConfig *apiContext.APIConfig) error {

	dockerAPIPort := apiConfig.DockerVMsConfig.DockerAPIPort
	dockerHostsList := apiConfig.DockerVMsConfig.Addresses

	for _, dockerHost := range dockerHostsList {
		dockerAddress := fmt.Sprintf("%s:%d", dockerHost, dockerAPIPort)
		if err := docker.HealthCheckAPI(dockerAddress); err != nil {
			return err
		}
	}

	return nil
}

func checkMongoDB(apiConfig *apiContext.APIConfig) error {

	_, err := db.Connect()

	if err != nil {
		mongoError := fmt.Sprintf("check mongoDB: %s", err)
		return errors.New(mongoError)
	}

	return nil
}

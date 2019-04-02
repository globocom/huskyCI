package util

import (
	"errors"
	"fmt"
	"os"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/db"
	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	docker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	mgo "gopkg.in/mgo.v2"
)

// CheckHuskyRequirements checks for all requirements needed before starting huskyCI.
func CheckHuskyRequirements(configAPI *apiContext.APIConfig) error {

	// check if all environment variables are properly set.
	if err := checkEnvVars(); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 12)

	// check if all docker hosts are up and running docker API.
	if err := checkDockerHosts(configAPI); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 13)

	// check if MongoDB is acessible and credentials received are working.
	if err := checkMongoDB(); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 14)

	// check if default securityTests are set into MongoDB.
	if err := checkEachSecurityTest(configAPI); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 15)

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
		finalError := fmt.Sprintf("Check environment variables: %s", errorString)
		return errors.New(finalError)
	}

	return nil
}

func checkDockerHosts(configAPI *apiContext.APIConfig) error {
	return docker.HealthCheckDockerAPI()
}

func checkMongoDB() error {
	if err := mongoHuskyCI.Connect(); err != nil {
		mongoError := fmt.Sprintf("Check MongoDB: %s", err)
		return errors.New(mongoError)
	}
	return nil
}

func checkEachSecurityTest(configAPI *apiContext.APIConfig) error {
	securityTests := []string{"enry", "gosec", "brakeman", "bandit", "retirejs", "safety"}
	for _, securityTest := range securityTests {
		if err := checkSecurityTest(securityTest, configAPI); err != nil {
			return err
		}
	}
	return nil
}

func checkSecurityTest(securityTestName string, configAPI *apiContext.APIConfig) error {

	securityTestConfig := types.SecurityTest{}

	switch securityTestName {
	case "enry":
		securityTestConfig = *configAPI.EnrySecurityTest
	case "gosec":
		securityTestConfig = *configAPI.GosecSecurityTest
	case "brakeman":
		securityTestConfig = *configAPI.BrakemanSecurityTest
	case "bandit":
		securityTestConfig = *configAPI.BanditSecurityTest
	case "retirejs":
		securityTestConfig = *configAPI.RetirejsSecurityTest
	case "safety":
		securityTestConfig = *configAPI.SafetySecurityTest
	default:
		return errors.New("securityTest name not defined")
	}

	securityTestQuery := map[string]interface{}{"name": securityTestName}
	_, err := db.FindOneDBSecurityTest(securityTestQuery)
	if err == mgo.ErrNotFound {
		// As securityTest is not set into MongoDB, huskyCI will insert it.
		if err := db.InsertDBSecurityTest(securityTestConfig); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

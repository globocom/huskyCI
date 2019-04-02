package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/globocom/huskyCI/api/analysis"
	apiContext "github.com/globocom/huskyCI/api/context"
	db "github.com/globocom/huskyCI/api/db/mongo"
	docker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/log"
	mgo "gopkg.in/mgo.v2"
)

// CheckHuskyRequirements checks for all requirements needed before starting huskyCI.
func CheckHuskyRequirements(configAPI *apiContext.APIConfig) error {

	// check if all environment variables are properly set.
	if err := checkEnvVars(); err != nil {
		return err
	}
	log.Info("checkHuskyRequirements", "SERVER", 12)

	// check if all docker hosts are up and running docker API.
	if err := checkDockerHosts(configAPI); err != nil {
		return err
	}
	log.Info("checkHuskyRequirements", "SERVER", 13)

	// check if MongoDB is acessible and credentials received are working.
	if err := checkMongoDB(); err != nil {
		return err
	}
	log.Info("checkHuskyRequirements", "SERVER", 14)

	// check if default securityTests are set into MongoDB.
	if err := checkDefaultSecurityTests(configAPI); err != nil {
		return err
	}
	log.Info("checkHuskyRequirements", "SERVER", 15)

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
	if err := db.Connect(); err != nil {
		mongoError := fmt.Sprintf("Check MongoDB: %s", err)
		return errors.New(mongoError)
	}
	return nil
}

func checkDefaultSecurityTests(configAPI *apiContext.APIConfig) error {
	enryQuery := map[string]interface{}{"name": "enry"}
	enry, err := analysis.FindOneDBSecurityTest(enryQuery)
	if err == mgo.ErrNotFound {
		// As Enry securityTest is not set into MongoDB, HuskyCI will insert it.
		log.Warning("checkDefaultSecurityTests", "SERVER", 201)
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
		log.Warning("checkDefaultSecurityTests", "SERVER", 202)
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
		log.Warning("checkDefaultSecurityTests", "SERVER", 203)
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
		log.Warning("checkDefaultSecurityTests", "SERVER", 204)
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
		log.Warning("checkDefaultSecurityTests", "SERVER", 205)
		retirejs = *configAPI.RetirejsSecurityTest
		if err := analysis.InsertDBSecurityTest(retirejs); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	safetyQuery := map[string]interface{}{"name": "safety"}
	safety, err := analysis.FindOneDBSecurityTest(safetyQuery)
	if err == mgo.ErrNotFound {
		// As Safety securityTest is not set into MongoDB, HuskyCI will insert it.
		log.Warning("checkDefaultSecurityTests", "SERVER", 206)
		safety = *configAPI.SafetySecurityTest
		if err := analysis.InsertDBSecurityTest(safety); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

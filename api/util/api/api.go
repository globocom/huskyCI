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
	"github.com/globocom/huskyCI/api/user"
	mgo "gopkg.in/mgo.v2"
)

// CheckHuskyRequirements checks for all requirements needed before starting huskyCI.
func (hU HuskyUtils) CheckHuskyRequirements(configAPI *apiContext.APIConfig) error {

	// check if all environment variables are properly set.
	if err := hU.CheckHandler.checkEnvVars(); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 12)

	// check if all docker hosts are up and running docker API.
	if err := hU.CheckHandler.checkDockerHosts(configAPI); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 13)

	// check if MongoDB is acessible and credentials received are working.
	if err := hU.CheckHandler.checkMongoDB(); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 14)

	// check if default securityTests are set into MongoDB.
	if err := hU.CheckHandler.checkEachSecurityTest(configAPI); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 15)

	// check if default user is set into MongoDB.
	if err := hU.CheckHandler.checkDefaultUser(configAPI); err != nil {
		return err
	}
	log.Info("CheckHuskyRequirements", "API-UTIL", 20)

	return nil
}

func (cH *CheckUtils) checkEnvVars() error {

	envVars := []string{
		// Logging
		// "HUSKYCI_LOGGING_GRAYLOG_ADDR", (optional)
		// "HUSKYCI_LOGGING_GRAYLOG_PROTO", (optional)
		// "HUSKYCI_LOGGING_GRAYLOG_APP_NAME", (optional)
		// "HUSKYCI_LOGGING_GRAYLOG_TAG", (optional)
		// "HUSKYCI_LOGGING_GRAYLOG_DEV", (optional)

		// Database:
		"HUSKYCI_DATABASE_MONGO_ADDR",
		"HUSKYCI_DATABASE_MONGO_DBNAME",
		"HUSKYCI_DATABASE_MONGO_DBUSERNAME",
		"HUSKYCI_DATABASE_MONGO_DBPASSWORD",
		// "HUSKYCI_DATABASE_MONGO_PORT", (optional)
		// "HUSKYCI_DATABASE_MONGO_TIMEOUT", (optional)
		// "HUSKYCI_DATABASE_MONGO_POOL_LIMIT", (optional)

		// Docker API:
		"HUSKYCI_DOCKERAPI_ADDR",
		"HUSKYCI_DOCKERAPI_CERT_PATH",
		"HUSKYCI_DOCKERAPI_CERT_FILE",
		"HUSKYCI_DOCKERAPI_CERT_KEY",
		"HUSKYCI_API_DEFAULT_USERNAME",
		"HUSKYCI_API_DEFAULT_PASSWORD",
		// "HUSKYCI_API_DEFAULT_ITERATIONS", (optional)
		// "HUSKYCI_API_DEFAULT_KEY_LENGTH", (optional)
		// "HUSKYCI_API_DEFAULT_HASH_FUNCTION", (optional)
		// "HUSKYCI_DOCKERAPI_CERT_FILE_VALUE", (optional)
		// "HUSKYCI_DOCKERAPI_CERT_KEY_VALUE", (optional)
		// "HUSKYCI_DOCKERAPI_API_TLS_CERT_VALUE", (optional)
		// "HUSKYCI_DOCKERAPI_API_TLS_KEY_VALUE", (optional)
		// "HUSKYCI_DOCKERAPI_CERT_CA", (optional)
		// "HUSKYCI_DOCKERAPI_CERT_CA_VALUE", (optional)
		// "HUSKYCI_DOCKERAPI_PORT", (optional)
		// "HUSKYCI_DOCKERAPI_TLS_VERIFY", (optional)
		// "HUSKYCI_DOCKERAPI_MAX_CONTAINERS_BEFORE_CLEANING", (optional)

		// huskyCI API:
		// "HUSKYCI_API_PORT", (optional)
		// "HUSKYCI_API_ENABLE_HTTPS", (optional)
		// "HUSKYCI_API_GIT_PRIVATE_SSH_KEY", (optional)
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

func (cH *CheckUtils) checkDockerHosts(configAPI *apiContext.APIConfig) error {
	// writes necessary keys for TLS to respective files
	if err := createAPIKeys(); err != nil {
		return err
	}

	return docker.HealthCheckDockerAPI()
}

func (cH *CheckUtils) checkMongoDB() error {
	if err := mongoHuskyCI.Connect(); err != nil {
		mongoError := fmt.Sprintf("Check MongoDB: %s", err)
		return errors.New(mongoError)
	}
	return nil
}

func (cH *CheckUtils) checkEachSecurityTest(configAPI *apiContext.APIConfig) error {
	securityTests := []string{"enry", "gitauthors", "gosec", "brakeman", "bandit", "npmaudit", "yarnaudit", "safety"}
	for _, securityTest := range securityTests {
		if err := checkSecurityTest(securityTest, configAPI); err != nil {
			errMsg := fmt.Sprintf("%s %s", securityTest, err)
			log.Error("checkEachSecurityTest", "API-UTIL", 1023, errMsg)
			return err
		}
		log.Info("checkEachSecurityTest", "API-UTIL", 19, securityTest)
	}
	return nil
}

func (cH *CheckUtils) checkDefaultUser(configAPI *apiContext.APIConfig) error {

	defaultUserQuery := map[string]interface{}{"username": user.DefaultAPIUser}
	_, err := db.FindOneDBUser(defaultUserQuery)
	if err != nil {
		if err == mgo.ErrNotFound {
			// user not found, add default user
			if err := user.InsertDefaultUser(); err != nil {
				return err
			}
		} else {
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
	case "gitauthors":
		securityTestConfig = *configAPI.GitAuthorsSecurityTest
	case "gosec":
		securityTestConfig = *configAPI.GosecSecurityTest
	case "brakeman":
		securityTestConfig = *configAPI.BrakemanSecurityTest
	case "bandit":
		securityTestConfig = *configAPI.BanditSecurityTest
	case "npmaudit":
		securityTestConfig = *configAPI.NpmAuditSecurityTest
	case "yarnaudit":
		securityTestConfig = *configAPI.YarnAuditSecurityTest
	case "safety":
		securityTestConfig = *configAPI.SafetySecurityTest
	default:
		return errors.New("securityTest name not defined")
	}

	securityTestQuery := map[string]interface{}{"name": securityTestName}
	_, err := db.UpsertOneDBSecurityTest(securityTestQuery, securityTestConfig)
	if err != nil {
		return err
	}
	return nil
}

func createAPIKeys() error {
	certValue, check := os.LookupEnv("HUSKYCI_DOCKERAPI_CERT_FILE_VALUE")
	if check {
		f, err := os.OpenFile("/home/application/current/api/cert.pem", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		_, err = f.WriteString(certValue)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	certKeyValue, check := os.LookupEnv("HUSKYCI_DOCKERAPI_CERT_KEY_VALUE")
	if check {
		f, err := os.OpenFile("/home/application/current/api/key.pem", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		_, err = f.WriteString(certKeyValue)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	apiCertValue, check := os.LookupEnv("HUSKYCI_DOCKERAPI_API_TLS_CERT_VALUE")
	if check {
		f, err := os.OpenFile("/home/application/current/api/api-tls-cert.pem", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		_, err = f.WriteString(apiCertValue)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	apiKeyValue, check := os.LookupEnv("HUSKYCI_DOCKERAPI_API_TLS_KEY_VALUE")
	if check {
		f, err := os.OpenFile("/home/application/current/api/api-tls-key.pem", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		_, err = f.WriteString(apiKeyValue)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	caValue, check := os.LookupEnv("HUSKYCI_DOCKERAPI_CERT_CA_VALUE")
	if check {
		f, err := os.OpenFile("/home/application/current/api/ca.pem", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		_, err = f.WriteString(caValue)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	return nil
}

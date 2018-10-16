package config

import (
	"errors"
	"fmt"
	"os"
)

// RepositoryURL stores the repository URL of the project to be analyzed.
var RepositoryURL string

// HuskyAPI stores the address of Husky's API.
var HuskyAPI string

// RepositoryBranch stores the repository branch of the project to be analyzed.
var RepositoryBranch string

// SetConfigs sets all configuration needed to start the client.
func SetConfigs() {
	RepositoryURL = os.Getenv(`HUSKYCI_REPO_URL`)
	RepositoryBranch = os.Getenv(`HUSKYCI_REPO_BRANCH`)
	HuskyAPI = os.Getenv(`HUSKYCI_API`)
}

// CheckEnvVars checks if all environment vars are set.
func CheckEnvVars() error {

	envVars := []string{
		"HUSKYCI_REPO_URL",
		"HUSKYCI_REPO_BRANCH",
		"HUSKYCI_API",
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
		finalError := fmt.Sprintf("%s", errorString)
		return errors.New(finalError)
	}
	return nil
}

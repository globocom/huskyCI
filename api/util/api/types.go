package util

import (
	apiContext "github.com/globocom/huskyCI/api/context"
)

// CheckInterface is the interface that stores all check functions.
type CheckInterface interface {
	checkEnvVars() error
	checkDockerHosts(configAPI *apiContext.APIConfig) error
	checkDB(configAPI *apiContext.APIConfig) error
	checkEachSecurityTest(configAPI *apiContext.APIConfig) error
	checkDefaultUser(configAPI *apiContext.APIConfig) error
}

// CheckUtils is the struct used for testing utils.
type CheckUtils struct{}

// HuskyUtils is the struct that stores the check handler used for testing.
type HuskyUtils struct {
	CheckHandler CheckInterface
}

// FakeCheck is the struct used for testing checks functions.
type FakeCheck struct {
	EnvVarsError          error
	DockerHostsError      error
	MongoDBError          error
	EachSecurityTestError error
	DefaultUserError      error
}

func (fC *FakeCheck) checkEnvVars() error {
	return fC.EnvVarsError
}

func (fC *FakeCheck) checkDockerHosts(configAPI *apiContext.APIConfig) error {
	return fC.DockerHostsError
}

func (fC *FakeCheck) checkDB(configAPI *apiContext.APIConfig) error {
	return fC.MongoDBError
}

func (fC *FakeCheck) checkEachSecurityTest(configAPI *apiContext.APIConfig) error {
	return fC.EachSecurityTestError
}

func (fC *FakeCheck) checkDefaultUser(configAPI *apiContext.APIConfig) error {
	return fC.DefaultUserError
}

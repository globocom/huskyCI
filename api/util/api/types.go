package util

import (
	apiContext "github.com/globocom/huskyCI/api/context"
)

type CheckInterface interface {
	checkEnvVars() error
	checkDockerHosts(configAPI *apiContext.APIConfig) error
	checkMongoDB() error
	checkEachSecurityTest(configAPI *apiContext.APIConfig) error
}

type CheckUtils struct{}

type HuskyUtils struct {
	CheckHandler CheckInterface
}

type FakeCheck struct {
	EnvVarsError          error
	DockerHostsError      error
	MongoDBError          error
	EachSecurityTestError error
}

func (fC *FakeCheck) checkEnvVars() error {
	return fC.EnvVarsError
}

func (fC *FakeCheck) checkDockerHosts(configAPI *apiContext.APIConfig) error {
	return fC.DockerHostsError
}

func (fC *FakeCheck) checkMongoDB() error {
	return fC.MongoDBError
}

func (fC *FakeCheck) checkEachSecurityTest(configAPI *apiContext.APIConfig) error {
	return fC.EachSecurityTestError
}

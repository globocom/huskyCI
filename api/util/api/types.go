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

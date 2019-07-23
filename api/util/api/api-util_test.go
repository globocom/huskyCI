package util

import (
	"errors"
	"github.com/globocom/glbgelf"
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeCheck struct {
	envVarsError          error
	dockerHostsError      error
	mongoDBError          error
	eachSecurityTestError error
}

func (fC *FakeCheck) checkEnvVars() error {
	return fC.envVarsError
}

func (fC *FakeCheck) checkDockerHosts(configAPI *apiContext.APIConfig) error {
	return fC.dockerHostsError
}

func (fC *FakeCheck) checkMongoDB() error {
	return fC.mongoDBError
}

func (fC *FakeCheck) checkEachSecurityTest(configAPI *apiContext.APIConfig) error {
	return fC.eachSecurityTestError
}

type CheckHuskyData struct {
	configApi             *apiContext.APIConfig
	envVarsError          error
	dockerHostsError      error
	mongoDBError          error
	eachSecurityTestError error
	expectedError         error
}

var checkHuskyTests = []CheckHuskyData{
	{&apiContext.APIConfig{}, errors.New("Falha ao verificar variáveis de ambiente"), nil, nil, nil, errors.New("Falha ao verificar variáveis de ambiente")},
	{&apiContext.APIConfig{}, nil, errors.New("Falha ao verificar docker host"), nil, nil, errors.New("Falha ao verificar docker host")},
	{&apiContext.APIConfig{}, nil, nil, errors.New("Erro ao verificar banco do mongo"), nil, errors.New("Erro ao verificar banco do mongo")},
	{&apiContext.APIConfig{}, nil, nil, nil, errors.New("Erro ao verificar testes de segurança"), errors.New("Erro ao verificar testes de segurança")},
	{&apiContext.APIConfig{}, nil, nil, nil, nil, nil},
}

func TestCheckHuskyRequirements(t *testing.T) {
	glbgelf.InitLogger("", "huskytest", "", true, "UDP")
	for _, checkHuskyTest := range checkHuskyTests {
		fakeCheck := &FakeCheck{
			envVarsError:          checkHuskyTest.envVarsError,
			dockerHostsError:      checkHuskyTest.dockerHostsError,
			mongoDBError:          checkHuskyTest.mongoDBError,
			eachSecurityTestError: checkHuskyTest.eachSecurityTestError,
		}
		huskyCheck := HuskyUtils{
			CheckHandler: fakeCheck,
		}
		err := huskyCheck.CheckHuskyRequirements(checkHuskyTest.configApi)
		assert.Equal(t, checkHuskyTest.expectedError, err)
	}
}

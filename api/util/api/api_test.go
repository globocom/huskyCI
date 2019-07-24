package util_test

import (
	"errors"
	"github.com/globocom/glbgelf"
	apiContext "github.com/globocom/huskyCI/api/context"
	apiUtil "github.com/globocom/huskyCI/api/util/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CheckHuskyData struct {
	configApi             *apiContext.APIConfig
	envVarsError          error
	dockerHostsError      error
	mongoDBError          error
	eachSecurityTestError error
	expectedError         error
}

var checkHuskyTests = []CheckHuskyData{
	// Scenario #1: checkEnvVars returns an error
	{&apiContext.APIConfig{}, errors.New("Failed verifying environment variables"), nil, nil, nil, errors.New("Failed verifying environment variables")},
	// Scenario #2: checkDockerHosts returns an error
	{&apiContext.APIConfig{}, nil, errors.New("Failed verifying docker API"), nil, nil, errors.New("Failed verifying docker API")},
	// Scenario #3: checkMongoDB returns an error
	{&apiContext.APIConfig{}, nil, nil, errors.New("Error verifying mongoDB"), nil, errors.New("Error verifying mongoDB")},
	// Scenario #4: checkEachSecurityTest returns an error
	{&apiContext.APIConfig{}, nil, nil, nil, errors.New("Error verifying security tests"), errors.New("Error verifying security tests")},
	// Scenario #5: Checks are successful!
	{&apiContext.APIConfig{}, nil, nil, nil, nil, nil},
}

var _ = Describe("Util API", func() {
	glbgelf.InitLogger("", "huskytest", "", true, "UDP")
	Describe("CheckHuskyRequirements", func() {
		Context("When checkEnvVars returns an error", func() {
			fakeCheck := &apiUtil.FakeCheck{
				EnvVarsError:          checkHuskyTests[0].envVarsError,
				DockerHostsError:      checkHuskyTests[0].dockerHostsError,
				MongoDBError:          checkHuskyTests[0].mongoDBError,
				EachSecurityTestError: checkHuskyTests[0].eachSecurityTestError,
			}
			huskyCheck := apiUtil.HuskyUtils{
				CheckHandler: fakeCheck,
			}
			It("Should return the same error", func() {
				Expect(huskyCheck.CheckHuskyRequirements(checkHuskyTests[0].configApi)).To(Equal(checkHuskyTests[0].expectedError))
			})
		})
		Context("When checkDockerHosts returns an error", func() {
			fakeCheck := &apiUtil.FakeCheck{
				EnvVarsError:          checkHuskyTests[1].envVarsError,
				DockerHostsError:      checkHuskyTests[1].dockerHostsError,
				MongoDBError:          checkHuskyTests[1].mongoDBError,
				EachSecurityTestError: checkHuskyTests[1].eachSecurityTestError,
			}
			huskyCheck := apiUtil.HuskyUtils{
				CheckHandler: fakeCheck,
			}
			It("Should return the same error", func() {
				Expect(huskyCheck.CheckHuskyRequirements(checkHuskyTests[1].configApi)).To(Equal(checkHuskyTests[1].expectedError))
			})
		})
		Context("When checkMongoDB returns an error", func() {
			fakeCheck := &apiUtil.FakeCheck{
				EnvVarsError:          checkHuskyTests[2].envVarsError,
				DockerHostsError:      checkHuskyTests[2].dockerHostsError,
				MongoDBError:          checkHuskyTests[2].mongoDBError,
				EachSecurityTestError: checkHuskyTests[2].eachSecurityTestError,
			}
			huskyCheck := apiUtil.HuskyUtils{
				CheckHandler: fakeCheck,
			}
			It("Should return the same error", func() {
				Expect(huskyCheck.CheckHuskyRequirements(checkHuskyTests[2].configApi)).To(Equal(checkHuskyTests[2].expectedError))
			})
		})
		Context("When checkEachSecurityTest returns an error", func() {
			fakeCheck := &apiUtil.FakeCheck{
				EnvVarsError:          checkHuskyTests[3].envVarsError,
				DockerHostsError:      checkHuskyTests[3].dockerHostsError,
				MongoDBError:          checkHuskyTests[3].mongoDBError,
				EachSecurityTestError: checkHuskyTests[3].eachSecurityTestError,
			}
			huskyCheck := apiUtil.HuskyUtils{
				CheckHandler: fakeCheck,
			}
			It("Should return the same error", func() {
				Expect(huskyCheck.CheckHuskyRequirements(checkHuskyTests[3].configApi)).To(Equal(checkHuskyTests[3].expectedError))
			})
		})
		Context("When all aux functions return a nil error", func() {
			fakeCheck := &apiUtil.FakeCheck{
				EnvVarsError:          checkHuskyTests[4].envVarsError,
				DockerHostsError:      checkHuskyTests[4].dockerHostsError,
				MongoDBError:          checkHuskyTests[4].mongoDBError,
				EachSecurityTestError: checkHuskyTests[4].eachSecurityTestError,
			}
			huskyCheck := apiUtil.HuskyUtils{
				CheckHandler: fakeCheck,
			}
			It("Should return nil", func() {
				Expect(huskyCheck.CheckHuskyRequirements(checkHuskyTests[4].configApi)).To(BeNil())
			})
		})
	})
})

package context_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
)

type FakeCaller struct {
	expectedIntegerValue         int
	expectedConvertStrToIntError error
	expectedEnvVar               string
	expectedSetConfigFileError   error
	expectedStringFromConfig     string
	expectedBoolFromConfig       bool
	expectedIntFromConfig        int
}

func (fC *FakeCaller) ConvertStrToInt(str string) (int, error) {
	return fC.expectedIntegerValue, fC.expectedConvertStrToIntError
}

func (fC *FakeCaller) GetEnvironmentVariable(envName string) string {
	return fC.expectedEnvVar
}

func (fC *FakeCaller) SetConfigFile(configName, configPath string) error {
	return fC.expectedSetConfigFileError
}

func (fC *FakeCaller) GetStringFromConfigFile(value string) string {
	return fC.expectedStringFromConfig
}

func (fC *FakeCaller) GetBoolFromConfigFile(value string) bool {
	return fC.expectedBoolFromConfig
}

func (fC *FakeCaller) GetIntFromConfigFile(value string) int {
	return fC.expectedIntFromConfig
}

func (fC *FakeCaller) GetTimeDurationInSeconds(duration int) time.Duration {
	return time.Duration(duration) * time.Second
}

var _ = Describe("Context", func() {
	Describe("GetAPIPort", func() {
		Context("When ConvertStrToInt returns an error", func() {
			It("Should return the expected 8888 port", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         0,
					expectedConvertStrToIntError: errors.New("Failed converting string to integer"),
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetAPIPort()).To(Equal(8888))
			})
		})
		Context("When ConvertStrToInt returns a valid port", func() {
			It("Should return the expected port", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         1234,
					expectedConvertStrToIntError: nil,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetAPIPort()).To(Equal(fakeCaller.expectedIntegerValue))
			})
		})
	})
	Describe("GetAPIUseTLS", func() {
		Context("When GetEnvironmentVariable returns a valid option", func() {
			It("Should return a true boolean", func() {
				fakeCaller := FakeCaller{
					expectedEnvVar: "True",
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetAPIUseTLS()).To(BeTrue())
			})
		})
		Context("When GetEnvironmentVariable returns a not valid option", func() {
			It("Should return a true boolean", func() {
				fakeCaller := FakeCaller{
					expectedEnvVar: "Invalid",
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetAPIUseTLS()).To(BeFalse())
			})
		})
	})
	Describe("GetGrayLogIsDev", func() {
		Context("When GetEnvironmentVariable returns valid option", func() {
			It("Should return a false boolean", func() {
				fakeCaller := FakeCaller{
					expectedEnvVar: "False",
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetGraylogIsDev()).To(BeFalse())
			})
		})
		Context("When GetEnvironmentVariable returns invalid option", func() {
			It("Should return a false boolean", func() {
				fakeCaller := FakeCaller{
					expectedEnvVar: "",
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetGraylogIsDev()).To(BeTrue())
			})
		})
	})
	Describe("GetDBPort", func() {
		Context("When ConvertStrToInt returns an error", func() {
			It("Should return 27017 port", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         0,
					expectedConvertStrToIntError: errors.New("Failed converting string to integer"),
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBPort()).To(Equal(27017))
			})
		})
		Context("When ConvertStrToInt returns an error", func() {
			It("Should return 27017 port", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         2222,
					expectedConvertStrToIntError: nil,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBPort()).To(Equal(fakeCaller.expectedIntegerValue))
			})
		})
	})
	Describe("GetDBTimeout", func() {
		Context("When ConvertStrToInt returns an error", func() {
			It("Should return 60s of duration", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         0,
					expectedConvertStrToIntError: errors.New("Failed converting string to integer"),
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBTimeout()).To(Equal(time.Duration(60) * time.Second))
			})
		})
		Context("When ConvertStrToInt returns a valid timeout", func() {
			It("Should return the expected duration", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         200,
					expectedConvertStrToIntError: nil,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBTimeout()).To(Equal(time.Duration(fakeCaller.expectedIntegerValue) * time.Second))
			})
		})
	})
	Describe("GetDBPoolLimit", func() {
		Context("When ConvertStr returns an error", func() {
			It("Should return a pool of connections with a size 1000", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         124,
					expectedConvertStrToIntError: errors.New("Error during the convertion from string to integer"),
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBPoolLimit()).To(Equal(1000))
			})
		})
		Context("When ConvertStr returns an negative value", func() {
			It("Should return a pool of connections with a size 1000", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         -24,
					expectedConvertStrToIntError: nil,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBPoolLimit()).To(Equal(1000))
			})
		})
		Context("When ConvertStr returns a valid value", func() {
			It("Should return a pool of connections with the expected size", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         50,
					expectedConvertStrToIntError: nil,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDBPoolLimit()).To(Equal(fakeCaller.expectedIntegerValue))
			})
		})
	})
	Describe("GetDockerAPIPort", func() {
		Context("When ConvertStrToInt returns an error", func() {
			It("Should return 2376 port", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         0,
					expectedConvertStrToIntError: errors.New("Error during the convertion from string to integer"),
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDockerAPIPort()).To(Equal(2376))
			})
		})
		Context("When ConvertStrToInt returns an error", func() {
			It("Should return 2376 port", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         25000,
					expectedConvertStrToIntError: nil,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDockerAPIPort()).To(Equal(fakeCaller.expectedIntegerValue))
			})
		})
	})
	Describe("GetDockerAPITLSVerify", func() {
		Context("When GetEnvironmentVariable returns a valid value", func() {
			It("Should return 0", func() {
				fakeCaller := FakeCaller{
					expectedEnvVar: "False",
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDockerAPITLSVerify()).To(Equal(0))
			})
		})
		Context("When GetEnvironmentVariable returns an invalid value", func() {
			It("Should return 1", func() {
				fakeCaller := FakeCaller{
					expectedEnvVar: "Invalid",
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				Expect(config.GetDockerAPITLSVerify()).To(Equal(1))
			})
		})
	})
	Describe("GetAPIConfig", func() {
		Context("When SetConfigFile returns an error", func() {
			It("Should return the expected error", func() {
				fakeCaller := FakeCaller{
					expectedSetConfigFileError: errors.New("Could not load configuration file"),
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				apiConfig, err := config.GetAPIConfig()
				Expect(apiConfig).To(BeNil())
				Expect(err).To(Equal(fakeCaller.expectedSetConfigFileError))
			})
		})
		Context("When SetConfigFile returns a nil error", func() {
			It("Should return the expected struct and a nil error", func() {
				fakeCaller := FakeCaller{
					expectedIntegerValue:         1234,
					expectedEnvVar:               "1",
					expectedConvertStrToIntError: nil,
					expectedSetConfigFileError:   nil,
					expectedStringFromConfig:     "teste",
					expectedBoolFromConfig:       true,
					expectedIntFromConfig:        1234,
				}
				config := DefaultConfig{
					Caller: &fakeCaller,
				}
				apiConfig, err := config.GetAPIConfig()
				expectedConfig := &APIConfig{
					Port:             fakeCaller.expectedIntegerValue,
					Version:          "0.14.0",
					ReleaseDate:      "2020-06-24",
					AllowOriginValue: fakeCaller.expectedEnvVar,
					UseTLS:           true,
					GitPrivateSSHKey: fakeCaller.expectedEnvVar,
					GraylogConfig: &GraylogConfig{
						Address:        fakeCaller.expectedEnvVar,
						Protocol:       fakeCaller.expectedEnvVar,
						AppName:        fakeCaller.expectedEnvVar,
						Tag:            fakeCaller.expectedEnvVar,
						DevelopmentEnv: true,
					},
					DBConfig: &DBConfig{
						Address:         fakeCaller.expectedEnvVar,
						DatabaseName:    fakeCaller.expectedEnvVar,
						Username:        fakeCaller.expectedEnvVar,
						Password:        fakeCaller.expectedEnvVar,
						Port:            fakeCaller.expectedIntegerValue,
						Timeout:         time.Duration(fakeCaller.expectedIntegerValue) * time.Second,
						PoolLimit:       fakeCaller.expectedIntegerValue,
						MaxOpenConns:    fakeCaller.expectedIntegerValue,
						MaxIdleConns:    fakeCaller.expectedIntegerValue,
						ConnMaxLifetime: time.Duration(fakeCaller.expectedIntegerValue) * time.Hour,
					},
					DockerHostsConfig: &DockerHostsConfig{
						Address:         "1",
						DockerAPIPort:   fakeCaller.expectedIntegerValue,
						PathCertificate: fakeCaller.expectedEnvVar,
						Host:            "1:1234",
						TLSVerify:       1,
					},
					EnrySecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					GitAuthorsSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					GosecSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					BanditSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					BrakemanSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					NpmAuditSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					YarnAuditSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					SafetySecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					GitleaksSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					SpotBugsSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					TFSecSecurityTest: &types.SecurityTest{
						Name:             fakeCaller.expectedStringFromConfig,
						Image:            fakeCaller.expectedStringFromConfig,
						ImageTag:         fakeCaller.expectedStringFromConfig,
						Cmd:              fakeCaller.expectedStringFromConfig,
						Type:             fakeCaller.expectedStringFromConfig,
						Language:         fakeCaller.expectedStringFromConfig,
						Default:          fakeCaller.expectedBoolFromConfig,
						TimeOutInSeconds: fakeCaller.expectedIntFromConfig,
					},
					DBInstance: &db.MongoRequests{},
					Cache:      apiConfig.Cache, // cannot be compared due to channels inside the structure
				}
				Expect(apiConfig).To(Equal(expectedConfig))
				Expect(err).To(BeNil())
			})
		})
	})
})

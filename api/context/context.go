// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/globocom/huskyCI/api/types"
	"github.com/spf13/viper"
)

// APIConfiguration holds all API configuration.
var APIConfiguration *APIConfig
var onceConfig sync.Once

// MongoConfig represents MongoDB configuration.
type MongoConfig struct {
	Address      string
	DatabaseName string
	Username     string
	Password     string
	Port         int
	Timeout      time.Duration
	PoolLimit    int
}

// DockerHostsConfig represents Docker Hosts configuration.
type DockerHostsConfig struct {
	Address              string
	DockerAPIPort        int
	Certificate          string
	PathCertificate      string
	Key                  string
	Host                 string
	TLSVerify            int
	MaxContainersAllowed int
}

// GraylogConfig represents Graylog configuration.
type GraylogConfig struct {
	Address        string
	Protocol       string
	AppName        string
	Tag            string
	DevelopmentEnv bool
}

// APIConfig represents API configuration.
type APIConfig struct {
	Port                 int
	Version              string
	ReleaseDate          string
	UseTLS               bool
	GitPrivateSSHKey     string
	GraylogConfig        *GraylogConfig
	MongoDBConfig        *MongoConfig
	DockerHostsConfig    *DockerHostsConfig
	EnrySecurityTest     *types.SecurityTest
	GosecSecurityTest    *types.SecurityTest
	BanditSecurityTest   *types.SecurityTest
	BrakemanSecurityTest *types.SecurityTest
	RetirejsSecurityTest *types.SecurityTest
	NpmAuditSecurityTest *types.SecurityTest
	SafetySecurityTest   *types.SecurityTest
}

type ExternalCalls struct{}

func (eC *ExternalCalls) SetConfigFile(configName, configPath string) error {
	viper.SetConfigFile(configName)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig()
}

func (eC *ExternalCalls) GetEnvironmentVariable(envName string) string {
	return os.Getenv(envName)
}

func (eC *ExternalCalls) ConvertStrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func (eC *ExternalCalls) GetTimeDurationInSeconds(duration int) time.Duration {
	return time.Duration(duration) * time.Second
}

func (eC *ExternalCalls) GetStringFromConfigFile(value string) string {
	return viper.GetString(value)
}

func (eC *ExternalCalls) GetBoolFromConfigFile(value string) bool {
	return viper.GetBool(value)
}

func (eC *ExternalCalls) GetIntFromConfigFile(value string) int {
	return viper.GetInt(value)
}

type CallerInterface interface {
	SetConfigFile(configName, configPath string) error
	GetStringFromConfigFile(value string) string
	GetBoolFromConfigFile(value string) bool
	GetIntFromConfigFile(value string) int
	GetEnvironmentVariable(envName string) string
	ConvertStrToInt(str string) (int, error)
	GetTimeDurationInSeconds(duration int) time.Duration
}

type DefaultConfig struct {
	Caller CallerInterface
}

// GetAPIConfig returns the instance of an APIConfig.
func (dF DefaultConfig) GetAPIConfig() (*APIConfig, error) {

	// load Viper using api/config.yml
	if err := dF.Caller.SetConfigFile("config", "api/"); err != nil {
		fmt.Println("Error reading Viper config: ", err)
		return nil, err
	}
	SetOnceConfig()
	return APIConfiguration
}

// SetOnceConfig sets APIConfiguration once
func (dF DefaultConfig) SetOnceConfig() {
	onceConfig.Do(func() {
		APIConfiguration = &APIConfig{
			Port:                 dF.GetAPIPort(),
			Version:              dF.GetAPIVersion(),
			ReleaseDate:          dF.GetAPIReleaseDate(),
			UseTLS:               dF.GetAPIUseTLS(),
			GitPrivateSSHKey:     dF.getGitPrivateSSHKey(),
			GraylogConfig:        dF.getGraylogConfig(),
			MongoDBConfig:        dF.getMongoConfig(),
			DockerHostsConfig:    dF.getDockerHostsConfig(),
			EnrySecurityTest:     dF.getSecurityTestConfig("enry"),
			GosecSecurityTest:    dF.getSecurityTestConfig("gosec"),
			BanditSecurityTest:   dF.getSecurityTestConfig("bandit"),
			BrakemanSecurityTest: dF.getSecurityTestConfig("brakeman"),
			RetirejsSecurityTest: dF.getSecurityTestConfig("retirejs"),
			NpmAuditSecurityTest: dF.getSecurityTestConfig("npmaudit"),
			SafetySecurityTest:   dF.getSecurityTestConfig("safety"),
		}
	})
}

func (dF DefaultConfig) GetAPIPort() int {
	apiPort, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_API_PORT"))
	if err != nil {
		apiPort = 8888
	}
	return apiPort
}

// GetAPIVersion returns current API version
func (dF DefaultConfig) GetAPIVersion() string {
	return "0.6.0"
}

// GetAPIReleaseDate returns current API release date
func (dF DefaultConfig) GetAPIReleaseDate() string {
	return "2019-07-18"
}

func (dF DefaultConfig) GetAPIUseTLS() bool {
	option := dF.Caller.GetEnvironmentVariable("HUSKYCI_API_ENABLE_HTTPS")
	if option == "true" || option == "1" || option == "TRUE" {
		return true
	}
	return false
}

func (dF DefaultConfig) getGitPrivateSSHKey() string {
	return dF.Caller.GetEnvironmentVariable("HUSKYCI_API_GIT_PRIVATE_SSH_KEY")
}

func (dF DefaultConfig) getGraylogConfig() *GraylogConfig {
	return &GraylogConfig{
		Address:        dF.Caller.GetEnvironmentVariable("HUSKYCI_LOGGING_GRAYLOG_ADDR"),
		Protocol:       dF.Caller.GetEnvironmentVariable("HUSKYCI_LOGGING_GRAYLOG_PROTO"),
		AppName:        dF.Caller.GetEnvironmentVariable("HUSKYCI_LOGGING_GRAYLOG_APP_NAME"),
		Tag:            dF.Caller.GetEnvironmentVariable("HUSKYCI_LOGGING_GRAYLOG_TAG"),
		DevelopmentEnv: GetGraylogIsDev(),
	}
}

func (dF DefaultConfig) GetGraylogIsDev() bool {
	option := dF.Caller.GetEnvironmentVariable("HUSKYCI_LOGGING_GRAYLOG_DEV")
	if option == "false" || option == "0" || option == "FALSE" {
		return false
	}
	return true
}

func (dF DefaultConfig) getMongoConfig() *MongoConfig {
	mongoHost := dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_ADDR")
	mongoPort := GetMongoPort()
	mongoAddress := fmt.Sprintf("%s:%d", mongoHost, mongoPort)
	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_DBNAME"),
		Username:     dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_DBUSERNAME"),
		Password:     dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_DBPASSWORD"),
		Port:         mongoPort,
		Timeout:      GetMongoTimeout(),
		PoolLimit:    GetMongoPoolLimit(),
	}
}

func (dF DefaultConfig) GetMongoPort() int {
	mongoPort, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_PORT"))
	if err != nil {
		return 27017
	}
	return mongoPort
}

func (dF DefaultConfig) GetMongoTimeout() time.Duration {
	mongoTimeout, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_TIMEOUT"))
	if err != nil {
		return dF.Caller.GetTimeDurationInSeconds(60)
	}
	return dF.Caller.GetTimeDurationInSeconds(mongoTimeout)
}

func (dF DefaultConfig) GetMongoPoolLimit() int {
	mongoPoolLimit, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_POOL_LIMIT"))
	if err != nil && mongoPoolLimit <= 0 {
		return 1000
	}
	return mongoPoolLimit
}

func (dF DefaultConfig) getDockerHostsConfig() *DockerHostsConfig {
	dockerAPIPort := GetDockerAPIPort()
	dockerHostsAddressesEnv := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_ADDR")
	dockerHostsAddresses := strings.Split(dockerHostsAddressesEnv, " ")
	dockerHostsCertificate := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_CERT_FILE")
	dockerHostsPathCertificates := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_CERT_PATH")
	dockerHostsKey := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_CERT_KEY")
	return &DockerHostsConfig{
		Address:              dockerHostsAddresses[0],
		DockerAPIPort:        dockerAPIPort,
		Certificate:          dockerHostsCertificate,
		PathCertificate:      dockerHostsPathCertificates,
		Key:                  dockerHostsKey,
		Host:                 fmt.Sprintf("%s:%d", dockerHostsAddresses[0], dockerAPIPort),
		TLSVerify:            GetDockerAPITLSVerify(),
		MaxContainersAllowed: GetMaxContainersAllowed(),
	}
}

func (dF DefaultConfig) GetDockerAPIPort() int {
	dockerAPIport, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_PORT"))
	if err != nil {
		return 2376
	}
	return dockerAPIport
}

func (dF DefaultConfig) GetDockerAPITLSVerify() int {
	option := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_TLS_VERIFY")
	if option == "false" || option == "0" || option == "FALSE" {
		return 0
	}
	return 1
}

func (dF DefaultConfig) getSecurityTestConfig(securityTestName string) *types.SecurityTest {
	return &types.SecurityTest{
		Name:             dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.name", securityTestName)),
		Image:            dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.image", securityTestName)),
		Cmd:              dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.cmd", securityTestName)),
		Language:         dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.language", securityTestName)),
		Default:          dF.Caller.GetBoolFromConfigFile(fmt.Sprintf("%s.default", securityTestName)),
		TimeOutInSeconds: dF.Caller.GetIntFromConfigFile(fmt.Sprintf("%s.timeOutInSeconds", securityTestName)),
	}
}

func (dF DefaultConfig) GetMaxContainersAllowed() int {
	maxContainersAllowed, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_MAX_CONTAINERS_BEFORE_CLEANING"))
	if err != nil {
		return 50
	}
	return maxContainersAllowed
}

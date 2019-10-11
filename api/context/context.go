// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/globocom/huskyCI/api/types"
)

// APIConfiguration holds all API configuration.
var (
	APIConfiguration *APIConfig
	onceConfig       sync.Once
	DefaultConf      *DefaultConfig
)

func init() {
	DefaultConf = &DefaultConfig{
		Caller: &ExternalCalls{},
	}
}

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
	Port                   int
	Version                string
	ReleaseDate            string
	AllowOriginValue       string
	UseTLS                 bool
	GitPrivateSSHKey       string
	GraylogConfig          *GraylogConfig
	MongoDBConfig          *MongoConfig
	DockerHostsConfig      *DockerHostsConfig
	EnrySecurityTest       *types.SecurityTest
	GitAuthorsSecurityTest *types.SecurityTest
	GosecSecurityTest      *types.SecurityTest
	BanditSecurityTest     *types.SecurityTest
	BrakemanSecurityTest   *types.SecurityTest
	NpmAuditSecurityTest   *types.SecurityTest
	YarnAuditSecurityTest  *types.SecurityTest
	SpotBugsSecurityTest   *types.SecurityTest
	GitleaksSecurityTest   *types.SecurityTest
	SafetySecurityTest     *types.SecurityTest
}

// DefaultConfig is the struct that stores the caller for testing.
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
	dF.SetOnceConfig()
	return APIConfiguration, nil
}

// SetOnceConfig sets APIConfiguration once
func (dF DefaultConfig) SetOnceConfig() {
	onceConfig.Do(func() {
		APIConfiguration = &APIConfig{
			Port:                   dF.GetAPIPort(),
			Version:                dF.GetAPIVersion(),
			ReleaseDate:            dF.GetAPIReleaseDate(),
			AllowOriginValue:       dF.GetAllowOriginValue(),
			UseTLS:                 dF.GetAPIUseTLS(),
			GitPrivateSSHKey:       dF.getGitPrivateSSHKey(),
			GraylogConfig:          dF.getGraylogConfig(),
			MongoDBConfig:          dF.getMongoConfig(),
			DockerHostsConfig:      dF.getDockerHostsConfig(),
			EnrySecurityTest:       dF.getSecurityTestConfig("enry"),
			GitAuthorsSecurityTest: dF.getSecurityTestConfig("gitauthors"),
			GosecSecurityTest:      dF.getSecurityTestConfig("gosec"),
			BanditSecurityTest:     dF.getSecurityTestConfig("bandit"),
			BrakemanSecurityTest:   dF.getSecurityTestConfig("brakeman"),
			NpmAuditSecurityTest:   dF.getSecurityTestConfig("npmaudit"),
			YarnAuditSecurityTest:  dF.getSecurityTestConfig("yarnaudit"),
			SpotbugsSecurityTest:   dF.getSecurityTestConfig("spotbugs"),
			GitleaksSecurityTest:   dF.getSecurityTestConfig("gitleaks"),
			SafetySecurityTest:     dF.getSecurityTestConfig("safety"),
		}
	})
}

// GetAPIPort will return the port number
// where HuskyCI will be listening to.
// If HUSKYCI_API_PORT is not set, it will
// return the default 8888 port.
func (dF DefaultConfig) GetAPIPort() int {
	apiPort, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_API_PORT"))
	if err != nil {
		apiPort = 8888
	}
	return apiPort
}

// GetAPIVersion returns current API version
func (dF DefaultConfig) GetAPIVersion() string {
	return "0.8.0"
}

// GetAPIReleaseDate returns current API release date
func (dF DefaultConfig) GetAPIReleaseDate() string {
	return "2019-09-30"
}

// GetAllowOriginValue returns the allow origin value
func (dF DefaultConfig) GetAllowOriginValue() string {
	return dF.Caller.GetEnvironmentVariable("HUSKYCI_API_ALLOW_ORIGIN_CORS")
}

// GetAPIUseTLS returns a boolean. If true, Husky API
// will be initialized with TLS. Otherwise, it won't.
// This depends on HUSKYCI_API_ENABLE_HTTPS variable.
func (dF DefaultConfig) GetAPIUseTLS() bool {
	option := dF.Caller.GetEnvironmentVariable("HUSKYCI_API_ENABLE_HTTPS")
	if strings.EqualFold(option, "true") || option == "1" {
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
		DevelopmentEnv: dF.GetGraylogIsDev(),
	}
}

// GetGraylogIsDev returns a true boolean if
// it is running in a development environment.
// This tells GlbGelf to generate logs only to
// stdout. Otherwise, it will return false. It
// depends on HUSKYCI_LOGGING_GRAYLOG_DEV env.
func (dF DefaultConfig) GetGraylogIsDev() bool {
	option := dF.Caller.GetEnvironmentVariable("HUSKYCI_LOGGING_GRAYLOG_DEV")
	if strings.EqualFold(option, "false") || option == "0" {
		return false
	}
	return true
}

func (dF DefaultConfig) getMongoConfig() *MongoConfig {
	mongoHost := dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_ADDR")
	mongoPort := dF.GetMongoPort()
	mongoAddress := fmt.Sprintf("%s:%d", mongoHost, mongoPort)
	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_DBNAME"),
		Username:     dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_DBUSERNAME"),
		Password:     dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_DBPASSWORD"),
		Port:         mongoPort,
		Timeout:      dF.GetMongoTimeout(),
		PoolLimit:    dF.GetMongoPoolLimit(),
	}
}

//GetMongoPort returns the port where MongoDB
// will be listening to. It depends on an env
// called HUSKYCI_DATABASE_MONGO_PORT.
func (dF DefaultConfig) GetMongoPort() int {
	mongoPort, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_PORT"))
	if err != nil {
		return 27017
	}
	return mongoPort
}

// GetMongoTimeout returns a time.Duration for
// duration of a connection with MongoDB. This
// depends on HUSKYCI_DATABASE_MONGO_TIMEOUT.
func (dF DefaultConfig) GetMongoTimeout() time.Duration {
	mongoTimeout, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_TIMEOUT"))
	if err != nil {
		return dF.Caller.GetTimeDurationInSeconds(60)
	}
	return dF.Caller.GetTimeDurationInSeconds(mongoTimeout)
}

// GetMongoPoolLimit returns an integer with
// the limit of pool of connections opened with
// MongoDB. This depends on an enviroment var
// called HUSKYCI_DATABASE_MONGO_POOL_LIMIT.
func (dF DefaultConfig) GetMongoPoolLimit() int {
	mongoPoolLimit, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_MONGO_POOL_LIMIT"))
	if err != nil || mongoPoolLimit <= 0 {
		return 1000
	}
	return mongoPoolLimit
}

func (dF DefaultConfig) getDockerHostsConfig() *DockerHostsConfig {
	dockerAPIPort := dF.GetDockerAPIPort()
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
		TLSVerify:            dF.GetDockerAPITLSVerify(),
		MaxContainersAllowed: dF.GetMaxContainersAllowed(),
	}
}

// GetDockerAPIPort will return the port number
// where Docker API will be listening to. This
// depends on HUSKYCI_DOCKERAPI_PORT.
func (dF DefaultConfig) GetDockerAPIPort() int {
	dockerAPIport, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_PORT"))
	if err != nil {
		return 2376
	}
	return dockerAPIport
}

// GetDockerAPITLSVerify returns an int that is
// interpreted as a boolean. If HUSKYCI_DOCKERAPI_TLS_VERIFY
// is false, it will return 0 and TLS won't be configured
// in the Docker API. Otherwise, it will return 1 and Docker
// API will use TLS protocol.
func (dF DefaultConfig) GetDockerAPITLSVerify() int {
	option := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_TLS_VERIFY")
	if strings.EqualFold(option, "false") || option == "0" {
		return 0
	}
	return 1
}

func (dF DefaultConfig) getSecurityTestConfig(securityTestName string) *types.SecurityTest {
	return &types.SecurityTest{
		Name:             dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.name", securityTestName)),
		Image:            dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.image", securityTestName)),
		ImageTag:         dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.imageTag", securityTestName)),
		Cmd:              dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.cmd", securityTestName)),
		Type:             dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.type", securityTestName)),
		Language:         dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.language", securityTestName)),
		Default:          dF.Caller.GetBoolFromConfigFile(fmt.Sprintf("%s.default", securityTestName)),
		TimeOutInSeconds: dF.Caller.GetIntFromConfigFile(fmt.Sprintf("%s.timeOutInSeconds", securityTestName)),
	}
}

// GetMaxContainersAllowed returns an interger the maximum number
// interpreted as the maximum number of containers initialized
// in parallel. It depends on the environment variable called
// HUSKYCI_DOCKERAPI_MAX_CONTAINERS_BEFORE_CLEANING.
func (dF DefaultConfig) GetMaxContainersAllowed() int {
	maxContainersAllowed, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_MAX_CONTAINERS_BEFORE_CLEANING"))
	if err != nil {
		return 50
	}
	return maxContainersAllowed
}

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
	Address         string
	DockerAPIPort   int
	Certificate     string
	PathCertificate string
	Key             string
	Host            string
	TLSVerify       int
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
	SafetySecurityTest   *types.SecurityTest
}

func init() {
	// load Viper using api/config.yml
	viper.SetConfigName("config")
	viper.AddConfigPath("api/")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading Viper config: ", err)
		os.Exit(1)
	}
	return
}

// GetAPIConfig returns the instance of an APIConfig.
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		APIConfiguration = &APIConfig{
			Port:                 getAPIPort(),
			Version:              getAPIVersion(),
			ReleaseDate:          getAPIReleaseDate(),
			UseTLS:               getAPIUseTLS(),
			GitPrivateSSHKey:     getGitPrivateSSHKey(),
			GraylogConfig:        getGraylogConfig(),
			MongoDBConfig:        getMongoConfig(),
			DockerHostsConfig:    getDockerHostsConfig(),
			EnrySecurityTest:     getSecurityTestConfig("enry"),
			GosecSecurityTest:    getSecurityTestConfig("gosec"),
			BanditSecurityTest:   getSecurityTestConfig("bandit"),
			BrakemanSecurityTest: getSecurityTestConfig("brakeman"),
			RetirejsSecurityTest: getSecurityTestConfig("retirejs"),
			SafetySecurityTest:   getSecurityTestConfig("safety"),
		}
	})
	return APIConfiguration
}

func getAPIPort() int {
	apiPort, err := strconv.Atoi(os.Getenv("HUSKYCI_API_PORT"))
	if err != nil {
		apiPort = 8888
	}
	return apiPort
}

func getAPIVersion() string {
	return "0.2.0"
}

func getAPIReleaseDate() string {
	return "2019-04-11"
}

func getAPIUseTLS() bool {
	option := os.Getenv("HUSKYCI_API_ENABLE_HTTPS")
	if option == "true" || option == "1" || option == "TRUE" {
		return true
	}
	return false
}

func getGitPrivateSSHKey() string {
	return os.Getenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY")
}

func getGraylogConfig() *GraylogConfig {
	return &GraylogConfig{
		Address:        os.Getenv("HUSKYCI_LOGGING_GRAYLOG_ADDR"),
		Protocol:       os.Getenv("HUSKYCI_LOGGING_GRAYLOG_PROTO"),
		AppName:        os.Getenv("HUSKYCI_LOGGING_GRAYLOG_APP_NAME"),
		Tag:            os.Getenv("HUSKYCI_LOGGING_GRAYLOG_TAG"),
		DevelopmentEnv: getGraylogIsDev(),
	}
}

func getGraylogIsDev() bool {
	option := os.Getenv("HUSKYCI_LOGGING_GRAYLOG_DEV")
	if option == "false" || option == "0" || option == "FALSE" {
		return false
	}
	return true
}

func getMongoConfig() *MongoConfig {
	mongoHost := os.Getenv("HUSKYCI_DATABASE_MONGO_ADDR")
	mongoPort := getMongoPort()
	mongoAddress := fmt.Sprintf("%s:%d", mongoHost, mongoPort)
	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: os.Getenv("HUSKYCI_DATABASE_MONGO_DBNAME"),
		Username:     os.Getenv("HUSKYCI_DATABASE_MONGO_DBUSERNAME"),
		Password:     os.Getenv("HUSKYCI_DATABASE_MONGO_DBPASSWORD"),
		Port:         mongoPort,
		Timeout:      getMongoTimeout(),
		PoolLimit:    getMongoPoolLimit(),
	}
}

func getMongoPort() int {
	mongoPort, err := strconv.Atoi(os.Getenv("HUSKYCI_DATABASE_MONGO_PORT"))
	if err != nil {
		return 27017
	}
	return mongoPort
}

func getMongoTimeout() time.Duration {
	mongoTimeout, err := strconv.Atoi(os.Getenv("HUSKYCI_DATABASE_MONGO_TIMEOUT"))
	if err != nil {
		return time.Duration(60) * time.Second
	}
	return time.Duration(mongoTimeout) * time.Second
}

func getMongoPoolLimit() int {
	mongoPoolLimit, err := strconv.Atoi(os.Getenv("HUSKYCI_DATABASE_MONGO_POOL_LIMIT"))
	if err != nil && mongoPoolLimit <= 0 {
		return 1000
	}
	return mongoPoolLimit
}

func getDockerHostsConfig() *DockerHostsConfig {
	dockerAPIPort := getDockerAPIPort()
	dockerHostsAddressesEnv := os.Getenv("HUSKYCI_DOCKERAPI_ADDR")
	dockerHostsAddresses := strings.Split(dockerHostsAddressesEnv, " ")
	dockerHostsCertificate := os.Getenv("HUSKYCI_DOCKERAPI_CERT_FILE")
	dockerHostsPathCertificates := os.Getenv("HUSKYCI_DOCKERAPI_CERT_PATH")
	dockerHostsKey := os.Getenv("HUSKYCI_DOCKERAPI_CERT_KEY")
	return &DockerHostsConfig{
		Address:         dockerHostsAddresses[0],
		DockerAPIPort:   dockerAPIPort,
		Certificate:     dockerHostsCertificate,
		PathCertificate: dockerHostsPathCertificates,
		Key:             dockerHostsKey,
		Host:            fmt.Sprintf("%s:%d", dockerHostsAddresses[0], dockerAPIPort),
		TLSVerify:       getDockerAPITLSVerify(),
	}
}

func getDockerAPIPort() int {
	dockerAPIport, err := strconv.Atoi(os.Getenv("HUSKYCI_DOCKERAPI_PORT"))
	if err != nil {
		return 2376
	}
	return dockerAPIport
}

func getDockerAPITLSVerify() int {
	option := os.Getenv("HUSKYCI_DOCKERAPI_TLS_VERIFY")
	if option == "false" || option == "0" || option == "FALSE" {
		return 0
	}
	return 1
}

func getSecurityTestConfig(securityTestName string) *types.SecurityTest {
	return &types.SecurityTest{
		Name:             viper.GetString(fmt.Sprintf("%s.name", securityTestName)),
		Image:            viper.GetString(fmt.Sprintf("%s.image", securityTestName)),
		Cmd:              viper.GetString(fmt.Sprintf("%s.cmd", securityTestName)),
		Language:         viper.GetString(fmt.Sprintf("%s.language", securityTestName)),
		Default:          viper.GetBool(fmt.Sprintf("%s.default", securityTestName)),
		TimeOutInSeconds: viper.GetInt(fmt.Sprintf("%s.timeOutInSeconds", securityTestName)),
	}
}

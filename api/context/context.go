// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/globocom/huskyCI/api/db"
	postgres "github.com/globocom/huskyCI/api/db/postgres"
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

// DBConfig represents DB configuration.
type DBConfig struct {
	Address         string
	DatabaseName    string
	Username        string
	Password        string
	Port            int
	Timeout         time.Duration
	PoolLimit       int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DockerHostsConfig represents Docker Hosts configuration.
type DockerHostsConfig struct {
	Address         string
	DockerAPIPort   int
	PathCertificate string
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
	Port                   int
	Version                string
	ReleaseDate            string
	AllowOriginValue       string
	UseTLS                 bool
	GitPrivateSSHKey       string
	GraylogConfig          *GraylogConfig
	DBConfig               *DBConfig
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
	TFSecSecurityTest      *types.SecurityTest
	DBInstance             db.Requests
	Cache                  *cache.Cache
}

// DefaultConfig is the struct that stores the caller for testing.
type DefaultConfig struct {
	Caller CallerInterface
}

// GetAPIConfig returns the instance of an APIConfig.
func (dF DefaultConfig) GetAPIConfig() (*APIConfig, error) {

	// load Viper using api/config.yml
	if err := dF.Caller.SetConfigFile("config", "."); err != nil {
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
			DBConfig:               dF.getDBConfig(),
			DockerHostsConfig:      dF.getDockerHostsConfig(),
			EnrySecurityTest:       dF.getSecurityTestConfig("enry"),
			GitAuthorsSecurityTest: dF.getSecurityTestConfig("gitauthors"),
			GosecSecurityTest:      dF.getSecurityTestConfig("gosec"),
			BanditSecurityTest:     dF.getSecurityTestConfig("bandit"),
			BrakemanSecurityTest:   dF.getSecurityTestConfig("brakeman"),
			NpmAuditSecurityTest:   dF.getSecurityTestConfig("npmaudit"),
			YarnAuditSecurityTest:  dF.getSecurityTestConfig("yarnaudit"),
			SpotBugsSecurityTest:   dF.getSecurityTestConfig("spotbugs"),
			GitleaksSecurityTest:   dF.getSecurityTestConfig("gitleaks"),
			SafetySecurityTest:     dF.getSecurityTestConfig("safety"),
			TFSecSecurityTest:      dF.getSecurityTestConfig("tfsec"),
			DBInstance:             dF.GetDB(),
			Cache:                  dF.GetCache(),
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
	return "0.14.0"
}

// GetAPIReleaseDate returns current API release date
func (dF DefaultConfig) GetAPIReleaseDate() string {
	return "2020-06-24"
}

// GetAllowOriginValue returns the allow origin value
func (dF DefaultConfig) GetAllowOriginValue() string {
	urlCORS := dF.Caller.GetEnvironmentVariable("HUSKYCI_API_ALLOW_ORIGIN_CORS")
	if urlCORS == "" {
		return "http://127.0.0.1:8888"
	}
	return urlCORS
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

func (dF DefaultConfig) getDBConfig() *DBConfig {
	return &DBConfig{
		Address:         dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_ADDR"),
		DatabaseName:    dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_NAME"),
		Username:        dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_USERNAME"),
		Password:        dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_PASSWORD"),
		Port:            dF.GetDBPort(),
		Timeout:         dF.GetDBTimeout(),
		PoolLimit:       dF.GetDBPoolLimit(),
		MaxOpenConns:    dF.GetMaxOpenConns(),
		MaxIdleConns:    dF.GetMaxIdleConns(),
		ConnMaxLifetime: dF.GetConnMaxLifetime(),
	}
}

// GetMaxOpenConns returns the maximum number
// of DB opened connections. It depends on an env
// called HUSKYCI_DATABASE_DB_MAX_OPEN_CONNS.
func (dF DefaultConfig) GetMaxOpenConns() int {
	maxOpenConns, err := dF.Caller.ConvertStrToInt(
		dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_MAX_OPEN_CONNS"))
	if err != nil {
		return 1
	}
	return maxOpenConns
}

// GetMaxIdleConns returns the maximum number
// of DB idle connections. It depends on an env
// called HUSKYCI_DATABASE_DB_MAX_IDLE_CONNS.
func (dF DefaultConfig) GetMaxIdleConns() int {
	maxIdleConns, err := dF.Caller.ConvertStrToInt(
		dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_MAX_IDLE_CONNS"))
	if err != nil {
		return 1
	}
	return maxIdleConns
}

// GetConnMaxLifetime returns the maximum duration
// of a DB connection. It depends on an env
// called HUSKYCI_DATABASE_DB_CONN_MAXLIFETIME.
func (dF DefaultConfig) GetConnMaxLifetime() time.Duration {
	connMaxLifetime, err := dF.Caller.ConvertStrToInt(
		dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_CONN_MAXLIFETIME"))
	if err != nil {
		return time.Hour
	}

	return time.Hour * time.Duration(connMaxLifetime)
}

// GetDBPort returns the port where DB
// will be listening to. It depends on an env
// called HUSKYCI_DATABASE_DB_PORT.
func (dF DefaultConfig) GetDBPort() int {
	dbPort, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_PORT"))
	if err != nil {
		return 27017
	}
	return dbPort
}

// GetDBTimeout returns a time.Duration for
// duration of a connection with DB. This
// depends on HUSKYCI_DATABASE_DB_TIMEOUT.
func (dF DefaultConfig) GetDBTimeout() time.Duration {
	dbTimeout, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_TIMEOUT"))
	if err != nil {
		return dF.Caller.GetTimeDurationInSeconds(60)
	}
	return dF.Caller.GetTimeDurationInSeconds(dbTimeout)
}

// GetDBPoolLimit returns an integer with
// the limit of pool of connections opened with
// DB. This depends on an enviroment var
// called HUSKYCI_DATABASE_DB_POOL_LIMIT.
func (dF DefaultConfig) GetDBPoolLimit() int {
	dbPoolLimit, err := dF.Caller.ConvertStrToInt(dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_DB_POOL_LIMIT"))
	if err != nil || dbPoolLimit <= 0 {
		return 1000
	}
	return dbPoolLimit
}

func (dF DefaultConfig) getDockerHostsConfig() *DockerHostsConfig {
	dockerAPIPort := dF.GetDockerAPIPort()
	dockerHostsAddressesEnv := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_ADDR")
	dockerHostsAddresses := strings.Split(dockerHostsAddressesEnv, " ")
	dockerHostsPathCertificates := dF.Caller.GetEnvironmentVariable("HUSKYCI_DOCKERAPI_CERT_PATH")
	return &DockerHostsConfig{
		Address:         dockerHostsAddresses[0],
		DockerAPIPort:   dockerAPIPort,
		PathCertificate: dockerHostsPathCertificates,
		Host:            fmt.Sprintf("%s:%d", dockerHostsAddresses[0], dockerAPIPort),
		TLSVerify:       dF.GetDockerAPITLSVerify(),
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

// GetDB returns a Requests implementation based on the
// on the type configured on HUSKYCI_DATABASE_TYPE env var.
// The default returns a MongoRequests that implements mongo
// queries.
func (dF DefaultConfig) GetDB() db.Requests {
	dB := dF.Caller.GetEnvironmentVariable("HUSKYCI_DATABASE_TYPE")
	if strings.EqualFold(dB, "postgres") {
		postgresOperations := postgres.PostgresHandler{}
		sqlConfig := postgres.SQLConfig{
			Postgres: &postgresOperations,
		}
		jsonHandler := db.JSONCaller{}
		sqlJSONRetriever := db.SQLJSONRetrieve{
			Psql:        &sqlConfig,
			JSONHandler: &jsonHandler,
		}
		postgres := db.PostgresRequests{
			DataRetriever: &sqlJSONRetriever,
			JSONHandler:   &jsonHandler,
		}
		return &postgres
	}
	return &db.MongoRequests{}
}

// GetCache returns a new cache based on the HUSKYCI_CACHE_DEFAULT_EXPIRATION
// and HUSKYCI_CACHE_CLEANUP_INTERVAL environment variables.
func (dF DefaultConfig) GetCache() *cache.Cache {
	var (
		defaultExpiration time.Duration
		cleanupInterval   time.Duration
		err               error
	)

	defaultExpiration, err = time.ParseDuration(
		dF.Caller.GetEnvironmentVariable("HUSKYCI_CACHE_DEFAULT_EXPIRATION"),
	)
	if err != nil {
		defaultExpiration = 5 * time.Minute
	}

	cleanupInterval, err = time.ParseDuration(
		dF.Caller.GetEnvironmentVariable("HUSKYCI_CACHE_CLEANUP_INTERVAL"),
	)
	if err != nil {
		cleanupInterval = 10 * time.Minute
	}

	return cache.New(defaultExpiration, cleanupInterval)
}

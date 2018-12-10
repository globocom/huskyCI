package context

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/globocom/huskyci/types"
	"github.com/spf13/viper"
)

// MongoConfig represents MongoDB configuration.
type MongoConfig struct {
	Address      string
	DatabaseName string
	Timeout      time.Duration
	Username     string
	Password     string
	PoolLimit    int
}

// DockerHostsConfig represents Docker Hosts configuration.
type DockerHostsConfig struct {
	Address       string
	DockerAPIPort int
	Certificate   string
	Key           string
	Host          string
}

// APIConfig represents API configuration.
type APIConfig struct {
	MongoDBConfig        *MongoConfig
	DockerHostsConfig    *DockerHostsConfig
	HuskyAPIPort         int
	EnrySecurityTest     *types.SecurityTest
	GosecSecurityTest    *types.SecurityTest
	BanditSecurityTest   *types.SecurityTest
	BrakemanSecurityTest *types.SecurityTest
}

var ApiConfig *APIConfig
var onceConfig sync.Once

func init() {
	// check if Viper's config.yaml is properly set.
	if err := loadViper(); err != nil {
		return
	}
}

// GetAPIConfig returns the instance of an APIConfig.
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		ApiConfig = &APIConfig{
			MongoDBConfig:        getMongoConfig(),
			DockerHostsConfig:    getDockerHostsConfig(),
			HuskyAPIPort:         getAPIHostPort(),
			EnrySecurityTest:     getEnryConfig(),
			GosecSecurityTest:    getGosecConfig(),
			BanditSecurityTest:   getBanditConfig(),
			BrakemanSecurityTest: getBrakemanConfig(),
		}
	})
	return ApiConfig
}

func getDockerHostsConfig() *DockerHostsConfig {

	dockerAPIPort := getDockerAPIPort()
	dockerHostsAddressesEnv := os.Getenv("DOCKER_HOSTS_LIST")
	dockerHostsAddresses := strings.Split(dockerHostsAddressesEnv, " ")
	dockerHostsCertificate := os.Getenv("DOCKER_HOSTS_CERT")
	dockerHostsKey := os.Getenv("DOCKER_HOSTS_KEY")

	return &DockerHostsConfig{
		Address:       dockerHostsAddresses[0],
		DockerAPIPort: dockerAPIPort,
		Certificate:   dockerHostsCertificate,
		Key:           dockerHostsKey,
		Host:          fmt.Sprintf("%s:%d", dockerHostsAddresses[0], dockerAPIPort),
	}
}

// getDockerAPIPort returns the port used by Docker API retrieved from an environment variable.
func getDockerAPIPort() int {
	dockerAPIport, err := strconv.Atoi(os.Getenv("DOCKER_API_PORT"))
	if err != nil {
		return 2376
	}
	return dockerAPIport
}

// getMongoConfig returns all MongoConfig retrieved from environment variables.
func getMongoConfig() *MongoConfig {

	mongoHost := os.Getenv("MONGO_HOST")
	mongoDatabaseName := os.Getenv("MONGO_DATABASE_NAME")
	mongoUserName := os.Getenv("MONGO_DATABASE_USERNAME")
	mongoPassword := os.Getenv("MONGO_DATABASE_PASSWORD")
	mongoTimeout := getMongoTimeout()
	mongoPort := getMongoPort()
	mongoAddress := fmt.Sprintf("%s:%d", mongoHost, mongoPort)
	mongoPoolLimit := getMongoPoolLimit()

	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: mongoDatabaseName,
		Timeout:      mongoTimeout,
		Username:     mongoUserName,
		Password:     mongoPassword,
		PoolLimit:    mongoPoolLimit,
	}
}

// getMongoPort returns the port used by MongoDB retrieved from an environment variable.
func getMongoPort() int {
	mongoPort, err := strconv.Atoi(os.Getenv("MONGO_PORT"))
	if err != nil {
		return 27017
	}
	return mongoPort
}

// getAPIHostPort returns the API port retrieved from an environment variable.
func getAPIHostPort() int {
	apiPort, err := strconv.Atoi(os.Getenv("HUSKY_API_PORT"))
	if err != nil {
		apiPort = 8888
	}
	return apiPort
}

func getMongoTimeout() time.Duration {
	mongoTimeout, err := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
	if err != nil {
		return time.Duration(60) * time.Second
	}
	return time.Duration(mongoTimeout) * time.Second
}

func getMongoPoolLimit() int {
	mongoPoolLimit, err := strconv.Atoi(os.Getenv("MONGO_POOL_LIMIT"))
	if err != nil {
		return 10
	}
	return mongoPoolLimit
}
func getEnryConfig() *types.SecurityTest {
	return &types.SecurityTest{
		Name:             viper.GetString("enry.name"),
		Image:            viper.GetString("enry.image"),
		Cmd:              viper.GetString("enry.cmd"),
		Language:         viper.GetString("enry.language"),
		Default:          viper.GetBool("enry.default"),
		TimeOutInSeconds: viper.GetInt("enry.timeOutInSeconds"),
	}
}

func getBrakemanConfig() *types.SecurityTest {
	return &types.SecurityTest{
		Name:             viper.GetString("brakeman.name"),
		Image:            viper.GetString("brakeman.image"),
		Cmd:              viper.GetString("brakeman.cmd"),
		Language:         viper.GetString("brakeman.language"),
		Default:          viper.GetBool("brakeman.default"),
		TimeOutInSeconds: viper.GetInt("brakeman.timeOutInSeconds"),
	}
}

func getGosecConfig() *types.SecurityTest {
	return &types.SecurityTest{
		Name:             viper.GetString("gosec.name"),
		Image:            viper.GetString("gosec.image"),
		Cmd:              viper.GetString("gosec.cmd"),
		Language:         viper.GetString("gosec.language"),
		Default:          viper.GetBool("gosec.default"),
		TimeOutInSeconds: viper.GetInt("gosec.timeOutInSeconds"),
	}
}

func getBanditConfig() *types.SecurityTest {
	return &types.SecurityTest{
		Name:             viper.GetString("bandit.name"),
		Image:            viper.GetString("bandit.image"),
		Cmd:              viper.GetString("bandit.cmd"),
		Language:         viper.GetString("bandit.language"),
		Default:          viper.GetBool("bandit.default"),
		TimeOutInSeconds: viper.GetInt("bandit.timeOutInSeconds"),
	}
}
func loadViper() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

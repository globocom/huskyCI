package context

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/globocom/husky/types"
	"github.com/spf13/viper"
)

// MongoConfig represents MongoDB configuration.
type MongoConfig struct {
	Address      string
	DatabaseName string
	Timeout      time.Duration
	Username     string
	Password     string
}

// DockerHostsConfig represents Docker Hosts configuration.
type DockerHostsConfig struct {
	Addresses     []string
	DockerAPIPort int
	Certificate   string
	Key           string
}

// APIConfig represents API configuration.
type APIConfig struct {
	MongoDBConfig      *MongoConfig
	DockerHostsConfig  *DockerHostsConfig
	HuskyAPIPort       int
	EnrySecurityTest   *types.SecurityTest
	GasSecurityTest    *types.SecurityTest
	BanditSecurityTest *types.SecurityTest
}

var apiConfig *APIConfig
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
		apiConfig = &APIConfig{
			MongoDBConfig:      getMongoConfig(),
			DockerHostsConfig:  getDockerHostsConfig(),
			HuskyAPIPort:       getAPIHostPort(),
			EnrySecurityTest:   getEnryConfig(),
			GasSecurityTest:    getGasConfig(),
			BanditSecurityTest: getBanditConfig(),
		}
	})
	return apiConfig
}

func getDockerHostsConfig() *DockerHostsConfig {

	dockerAPIPort := getDockerAPIPort()
	dockerHostsAddressesEnv := os.Getenv("DOCKER_HOSTS_LIST")
	dockerHostsAddresses := strings.Split(dockerHostsAddressesEnv, " ")
	dockerHostsCertificate := os.Getenv("DOCKER_HOSTS_CERT")
	dockerHostsKey := os.Getenv("DOCKER_HOSTS_KEY")

	return &DockerHostsConfig{
		Addresses:     dockerHostsAddresses,
		DockerAPIPort: dockerAPIPort,
		Certificate:   dockerHostsCertificate,
		Key:           dockerHostsKey,
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

	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: mongoDatabaseName,
		Timeout:      mongoTimeout,
		Username:     mongoUserName,
		Password:     mongoPassword,
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

func getGasConfig() *types.SecurityTest {
	return &types.SecurityTest{
		Name:             viper.GetString("gas.name"),
		Image:            viper.GetString("gas.image"),
		Cmd:              viper.GetString("gas.cmd"),
		Language:         viper.GetString("gas.language"),
		Default:          viper.GetBool("gas.default"),
		TimeOutInSeconds: viper.GetInt("gas.timeOutInSeconds"),
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

package context

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// ExternalCalls is the extruct that performs exernal calls.
type ExternalCalls struct{}

// SetConfigFile will set a configfile into a path.
func (eC *ExternalCalls) SetConfigFile(configName, configPath string) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig()
}

// GetEnvironmentVariable will return the value of an env var.
func (eC *ExternalCalls) GetEnvironmentVariable(envName string) string {
	return os.Getenv(envName)
}

// ConvertStrToInt converts a string into int.
func (eC *ExternalCalls) ConvertStrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// GetTimeDurationInSeconds returnin the number of seconds of a duration.
func (eC *ExternalCalls) GetTimeDurationInSeconds(duration int) time.Duration {
	return time.Duration(duration) * time.Second
}

// GetStringFromConfigFile returns a string from a config file.
func (eC *ExternalCalls) GetStringFromConfigFile(value string) string {
	return viper.GetString(value)
}

// GetBoolFromConfigFile returns a bool from a config file.
func (eC *ExternalCalls) GetBoolFromConfigFile(value string) bool {
	return viper.GetBool(value)
}

// GetIntFromConfigFile returns a int from a config file.
func (eC *ExternalCalls) GetIntFromConfigFile(value string) int {
	return viper.GetInt(value)
}

// CallerInterface is the interface that stores all external call functions.
type CallerInterface interface {
	SetConfigFile(configName, configPath string) error
	GetStringFromConfigFile(value string) string
	GetBoolFromConfigFile(value string) bool
	GetIntFromConfigFile(value string) int
	GetEnvironmentVariable(envName string) string
	ConvertStrToInt(str string) (int, error)
	GetTimeDurationInSeconds(duration int) time.Duration
}

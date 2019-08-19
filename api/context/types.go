package context

import (
	"github.com/spf13/viper"
	"os"
	"strconv"
	"time"
)

type ExternalCalls struct{}

func (eC *ExternalCalls) SetConfigFile(configName, configPath string) error {
	viper.SetConfigName(configName)
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

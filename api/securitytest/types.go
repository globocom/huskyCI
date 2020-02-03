package securitytest

import (
	"github.com/spf13/viper"
)

// ViperCalls is the extruct that performs exernal calls.
type ViperCalls struct{}

// SetConfigFile will set a configfile into a path.
func (vC *ViperCalls) SetConfigFile(configName, configPath string) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig()
}

// GetStringFromConfigFile returns a string from a config file.
func (vC *ViperCalls) GetStringFromConfigFile(value string) string {
	return viper.GetString(value)
}

// GetBoolFromConfigFile returns a bool from a config file.
func (vC *ViperCalls) GetBoolFromConfigFile(value string) bool {
	return viper.GetBool(value)
}

// GetIntFromConfigFile returns a int from a config file.
func (vC *ViperCalls) GetIntFromConfigFile(value string) int {
	return viper.GetInt(value)
}

// ViperInterface is the interface that stores all external call functions.
type ViperInterface interface {
	SetConfigFile(configName, configPath string) error
	GetStringFromConfigFile(value string) string
	GetBoolFromConfigFile(value string) bool
	GetIntFromConfigFile(value string) int
}

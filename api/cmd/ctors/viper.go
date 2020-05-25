package ctors

import (
	"time"

	"github.com/spf13/viper"
)

// NewSettings retuns a new Viper instance and an error
func NewSettings() (*viper.Viper, error) {

	viperInstance := viper.New()
	viperInstance.AutomaticEnv()
	viperInstance.AddConfigPath(".")
	viperInstance.SetConfigName("config")
	viperInstance.ReadInConfig()

	viperInstance.SetDefault("HUSKYCI_ARTIFACTORY_URL", "docker.io")
	viperInstance.SetDefault("HUSKYCI_API_PORT", 8888)
	viperInstance.SetDefault("HUSKYCI_API_ALLOW_ORIGIN_CORS", "http://127.0.0.1:1981")
	viperInstance.SetDefault("HUSKYCI_API_ENABLE_HTTPS", false)
	viperInstance.SetDefault("HUSKYCI_LOGGING_GRAYLOG_DEV", true)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_MAX_OPEN_CONNS", 1)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_MAX_IDLE_CONNS", 1)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_CONN_MAXLIFETIME", 10*time.Minute)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_PORT", 27017)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_TIMEOUT", 30*time.Second)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_FAIL_FAST", true)
	viperInstance.SetDefault("HUSKYCI_DATABASE_DB_POOL_LIMIT", 1000)
	viperInstance.SetDefault("HUSKYCI_DOCKERAPI_PORT", 2376)
	viperInstance.SetDefault("HUSKYCI_DOCKERAPI_TLS_VERIFY", true)
	viperInstance.SetDefault("HUSKYCI_DATABASE_TYPE", "mongo")
	viperInstance.SetDefault("HUSKYCI_RUNNER_TYPE", "dockerapi")

	return viperInstance, nil
}

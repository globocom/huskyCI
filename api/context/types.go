package context

import (
	"os"
	"strconv"
	"time"
)

// ExternalCalls is the extruct that performs exernal calls.
type ExternalCalls struct{}

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

// CallerInterface is the interface that stores all external call functions.
type CallerInterface interface {
	GetEnvironmentVariable(envName string) string
	GetTimeDurationInSeconds(duration int) time.Duration
	ConvertStrToInt(str string) (int, error)
}

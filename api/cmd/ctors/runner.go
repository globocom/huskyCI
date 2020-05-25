package ctors

import (
	"errors"

	"github.com/globocom/huskyCI/api/runner"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// NewRunnerSession starts a new runner session.
func NewRunnerSession(lc fx.Lifecycle, settings *viper.Viper) (runner.RSession, error) {

	runnerType := settings.GetString("HUSKYCI_RUNNER_TYPE")

	if runnerType != "dockerapi" {
		return nil, errors.New("runner not dockerapi")
	}

	return runner.NewDockerAPISession(lc, settings)
}

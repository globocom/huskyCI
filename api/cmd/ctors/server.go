package ctors

import (
	"log"

	"github.com/globocom/huskyCI/api/database"
	"github.com/globocom/huskyCI/api/runner"
	"github.com/globocom/huskyCI/api/securitytest"
	"github.com/globocom/huskyCI/api/server"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

// EchoServerParams is the struct that hold all server params
type EchoServerParams struct {
	dig.In

	Settings          *viper.Viper
	Logger            *log.Logger
	DatabaseSession   database.DBSession
	RunnerSession     runner.RSession
	SecurityTestStore securitytest.Store
}

// NewEchoServer creates a new server
func NewEchoServer(params EchoServerParams) (*server.EchoServer, error) {
	server, err := server.New(
		params.Settings,
		params.Logger,
		params.DatabaseSession,
		params.RunnerSession,
		params.SecurityTestStore,
	)
	if err != nil {
		return nil, err
	}
	return server, nil
}

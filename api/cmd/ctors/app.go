package ctors

import (
	"context"
	"log"

	"github.com/globocom/huskyCI/api/app"
	"github.com/globocom/huskyCI/api/server"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// ApplicationParams is the params
type ApplicationParams struct {
	dig.In

	Logger   *log.Logger
	Settings *viper.Viper
	Server   *server.EchoServer
}

// NewApplication is the app ctor
func NewApplication(lc fx.Lifecycle, params ApplicationParams) (*app.Application, error) {
	app := app.New(
		params.Logger,
		params.Settings,
		params.Server,
	)
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return app.Start()
		},
		OnStop: func(context.Context) error {
			return app.Stop()
		},
	})
	return app, nil
}

package main

import (
	"github.com/globocom/huskyCI/api/app"
	"github.com/globocom/huskyCI/api/cmd/ctors"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			// app
			ctors.NewApplication,

			// core
			ctors.NewSettings,
			ctors.NewLogger,
			ctors.NewEchoServer,

			// sessions
			ctors.NewDatabaseSession,
			ctors.NewRunnerSession,

			// security tests
			ctors.NewSecurityTestStore,
		),
		fx.Invoke(func(app *app.Application) {}),
	).Run()
}

package main

import (
	"github.com/globocom/huskyCI/api/app"
	"github.com/globocom/huskyCI/api/cmd/ctors"
	"go.uber.org/fx"
)

func main() {

	// time.Sleep(100000 * time.Hour)

	fx.New(
		fx.Provide(

			// the following functions will start each application dependency
			// before invoking it. As the output of some of them are required
			// by each other, the fx lib will handle this properly. For example:
			// NewSecurityTestStore needs the output of NewSettings to be invoked.

			// core
			ctors.NewSettings,
			ctors.NewLogger,
			ctors.NewEchoServer,

			// sessions
			ctors.NewDatabaseSession,
			ctors.NewRunnerSession,

			// security tests default values
			ctors.NewSecurityTestStore,

			// app
			ctors.NewApplication,
		),
		fx.Invoke(func(app *app.Application) {}),
	).Run()
}

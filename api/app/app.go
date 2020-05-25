package app

import (
	"fmt"
	"log"

	"github.com/globocom/huskyCI/api/server"
	"github.com/spf13/viper"
)

// Application is the application
type Application struct {
	Logger   *log.Logger
	Settings *viper.Viper
	Server   *server.EchoServer
}

// New returns a new application
func New(logger *log.Logger, settings *viper.Viper, server *server.EchoServer) *Application {
	return &Application{
		Logger:   logger,
		Settings: settings,
		Server:   server,
	}
}

// Start starts the application
func (a *Application) Start() error {
	fmt.Println("Starting the application...")
	a.Server.SetRoutes()
	return a.Server.Run()
}

// Stop stops the application
func (a *Application) Stop() error {
	fmt.Println("Stopping the application...")
	return a.Server.Stop()
}

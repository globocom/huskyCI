package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/globocom/huskyCI/api/database"
	"github.com/globocom/huskyCI/api/runner"
	"github.com/globocom/huskyCI/api/securitytest"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// EchoServer is the struct that holds the server info
type EchoServer struct {
	Port              string
	Engine            *echo.Echo
	DatabaseSession   database.DBSession
	RunnerSession     runner.RSession
	SecurityTestStore securitytest.Store
}

// New returns a new Echo server
func New(
	settings *viper.Viper,
	logger *log.Logger,
	databaseSession database.DBSession,
	runnerSession runner.RSession,
	securityTestStore securitytest.Store) (*EchoServer, error) {

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	// Middlewares
	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())
	echoInstance.Use(middleware.Secure())

	server := &EchoServer{
		Port:              fmt.Sprintf(":%d", settings.GetInt("HUSKYCI_API_PORT")),
		Engine:            echoInstance,
		DatabaseSession:   databaseSession,
		RunnerSession:     runnerSession,
		SecurityTestStore: securityTestStore,
	}

	return server, nil
}

// SetRoutes register all echo routes
func (es *EchoServer) SetRoutes() {
	es.Engine.GET("/healthcheck", es.HealthCheck)
}

// Run runs the server
func (es *EchoServer) Run() error {
	go func() {
		err := es.Engine.Start(es.Port)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

// Stop stops the server gracefully
func (es *EchoServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go func() {
		err := es.Engine.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

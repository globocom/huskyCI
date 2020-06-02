package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Defining postgres plugin
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Postgres implements the Database Interface
type Postgres struct {
	Session *sql.DB
}

// NewPostgresSession starts a new Postgres session
func NewPostgresSession(lc fx.Lifecycle, settings *viper.Viper) (*Postgres, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		settings.GetString("HUSKYCI_DATABASE_DB_USERNAME"),
		settings.GetString("HUSKYCI_DATABASE_DB_PASSWORD"),
		settings.GetString("HUSKYCI_DATABASE_DB_ADDR"),
		settings.GetString("HUSKYCI_DATABASE_DB_NAME"),
	)

	postgresSession, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	postgresSession.SetMaxOpenConns(settings.GetInt("HUSKYCI_DATABASE_DB_MAX_OPEN_CONNS"))
	postgresSession.SetMaxIdleConns(settings.GetInt("HUSKYCI_DATABASE_DB_MAX_IDLE_CONNS"))
	postgresSession.SetConnMaxLifetime(10 * time.Minute)

	databaseSession := &Postgres{
		Session: postgresSession,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return databaseSession.Ping()
		},
		OnStop: func(ctx context.Context) error {
			return databaseSession.Close()
		},
	})

	return databaseSession, nil
}

// Ping checks the Postgres session
func (p *Postgres) Ping() error {
	fmt.Println("Checking Postgres Session...")
	return p.Session.Ping()
}

// Close closes the Postgres session
func (p *Postgres) Close() error {
	fmt.Println("Closing Postgres Session...")
	return p.Session.Close()
}

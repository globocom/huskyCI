package database

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	mgo "gopkg.in/mgo.v2"
)

// MongoDB implements the Database Interface
type MongoDB struct {
	Session *mgo.Session
}

// NewMongoSession starts a new MongoDB session.
func NewMongoSession(lc fx.Lifecycle, settings *viper.Viper) (*MongoDB, error) {

	dialInfo, err := mgo.ParseURL(settings.GetString("HUSKYCI_DATABASE_DB_ADDR"))
	if err != nil {
		return nil, err
	}

	dialInfo.Timeout = 30 * time.Second
	dialInfo.Username = settings.GetString("HUSKYCI_DATABASE_DB_USERNAME")
	dialInfo.Password = settings.GetString("HUSKYCI_DATABASE_DB_PASSWORD")
	dialInfo.Database = settings.GetString("HUSKYCI_DATABASE_DB_NAME")
	dialInfo.PoolLimit = settings.GetInt("HUSKYCI_DATABASE_DB_POOL_LIMIT")
	dialInfo.FailFast = settings.GetBool("HUSKYCI_DATABASE_DB_FAIL_FAST")

	mongoSession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	databaseSession := &MongoDB{
		Session: mongoSession,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return databaseSession.Ping()
		},
		OnStop: func(ctx context.Context) error {
			databaseSession.Close()
			return nil
		},
	})

	return databaseSession, nil
}

// Ping checks the MongoDB session
func (m *MongoDB) Ping() error {
	fmt.Println("Checking MongoDB Session...")
	return m.Session.Ping()
}

// Close closes the MongoDB session
func (m *MongoDB) Close() error {
	fmt.Println("Closing MongoDB Session...")
	m.Session.Close()
	return nil
}

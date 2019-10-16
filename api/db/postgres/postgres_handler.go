package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	_ "github.com/lib/pq" // Defining postgres plugin
)

var (
	doOnce sync.Once
	db     *sql.DB = nil
	dbErr  error   = nil
)

// ConfigureDB will establish a new connection
// with the postgres DB referenced in the arguments
func (pHandler *PostgresHandler) ConfigureDB(
	address string,
	username string,
	password string,
	dbName string) error {
	doOnce.Do(func() {
		connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			username,
			password,
			address,
			dbName)
		db, dbErr = sql.Open("posgres", connStr)
	})
	return dbErr
}

// ConfigurePool will set the pool of connections
// Postgres based on the values passed in its arguments.
func (pHandler *PostgresHandler) ConfigurePool(maxOpenConns, maxIdleConns int, connLT time.Duration) {
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connLT)
}

// CloseDB will call the function to close
// Postgres connection that has an opened status
// in the pool of connections.
func (pHandler *PostgresHandler) CloseDB() error {
	return db.Close()
}

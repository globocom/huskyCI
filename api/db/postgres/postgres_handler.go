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
		db, dbErr = sql.Open("postgres", connStr)
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

// ConfigureQuery will call the function Query
// to make a query to DB and it will return the
// rows from DB. This will be stored in Rows pointer
// and the error will be returned.
func (pHandler *PostgresHandler) ConfigureQuery(query string, args ...interface{}) error {
	var err error
	pHandler.Rows, err = db.Query(query, args...)
	return err
}

// CloseRows will call Close function from db.Rows struct.
// It will destroy all rows returned from the query.
func (pHandler *PostgresHandler) CloseRows() error {
	return pHandler.Rows.Close()
}

// GetColumns will return all the columns stored in the
// Rows. If it won't be possible, it will return an empty
// slice and an error.
func (pHandler *PostgresHandler) GetColumns() ([]string, error) {
	return pHandler.Rows.Columns()
}

// HasNextRow will call Next from sql.Rows and will return
// a boolean. True if has a next row to be read. Otherwise,
// the boolean will be false.
func (pHandler *PostgresHandler) HasNextRow() bool {
	return pHandler.Rows.Next()
}

// ScanRow will call Scan from db.Rows and return an error
// if the scan fails. The result of the scan will be returned
// on dest slice interface.
func (pHandler *PostgresHandler) ScanRow(dest ...interface{}) error {
	return pHandler.Rows.Scan(dest...)
}

// GetRowsErr will call Err from db.Rows to verify if an
// error ocurred during rows retrieval.
func (pHandler *PostgresHandler) GetRowsErr() error {
	return pHandler.Rows.Err()
}

// Exec will run Exec function with. the query passed in the argument.
// It will store db.Result struct in PostgresHandler.
func (pHandler *PostgresHandler) Exec(query string, args ...interface{}) error {
	var err error
	pHandler.Result, err = db.Exec(query, args...)
	return err
}

// GetRowsAffected will call RowsAffected to return the number of
// rows affected during query process in Exec. In case of any failures in query,
// an error will be returned.
func (pHandler *PostgresHandler) GetRowsAffected() (int64, error) {
	return pHandler.Result.RowsAffected()
}

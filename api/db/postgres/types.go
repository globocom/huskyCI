package db

import (
	"database/sql"
)

// PostgresHandler implements
// the PostgresOperations interface.
type PostgresHandler struct {
	Rows   *sql.Rows
	Result sql.Result
}

// PostgresOperations defines the functions
// that will call the postgres database and
// deal with the generated queries. It makes
// all the required interactions with Postgres
// database directly.
type PostgresOperations interface {
	ConfigureDB() error
	ConfigurePool()
	CloseDB() error
	ConfigureQuery(query string, args ...interface{}) error
	CloseRows() error
	GetColumns() ([]string, error)
	HasNextRow() bool
	ScanRow(dest ...interface{}) error
	GetRowsErr() error
	Exec(query string, args ...interface{}) error
	GetRowsAffected() (int64, error)
}

// SQLConfig will implement the required logic
// for data interaction with Postgres. It will
// make the calls generic so any data could be
// retrieved or any valid queries requested without
// broken contract with DB.
type SQLConfig struct {
	Postgres PostgresOperations
}

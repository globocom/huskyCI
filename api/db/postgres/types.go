package db


import (
	"database/sql"
)

type PostgresHandler struct {
	Rows *sql.Rows
	Result sql.Result
}

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

type SqlConfig struct {
	Postgres PostgresOperations
}
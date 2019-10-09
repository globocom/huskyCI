package db


type PostgresHandler struct{}

type PostgresOperations interface {
	ConfigureDB() error
	ConfigurePool()
	CloseDB() error
}

type SqlConfig struct {
	Postgres PostgresOperations
}
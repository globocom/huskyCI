package db


import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
	"fmt"
	config "github.com/globocom/huskyCI/api/context"
)


var (
	doOnce sync.Once
	db *sql.DB = nil
	dbErr error = nil
	postgresConfig *config.MongoConfig = nil
)

func (pHandler *PostgresHandler) ConfigureDB() error {
	doOnce.Do(func() {
		postgresConfig = config.APIConfiguration.MongoDBConfig
		connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			postgresConfig.Username, 
			postgresConfig.Password,
			postgresConfig.Address,
			postgresConfig.DatabaseName)

		db, dbErr = sql.Open("posgres", connStr)
	})
	return dbErr
}


func (pHandler *PostgresHandler) ConfigurePool() {
	db.SetMaxOpenConns(postgresConfig.MaxOpenConns)
	db.SetMaxIdleConns(postgresConfig.MaxIdleConns)
	db.SetConnMaxLifetime(postgresConfig.ConnMaxLifetime)
}

func (pHandler *PostgresHandler) CloseDB() error {
	return db.Close()
}
package db

import (
	"time"
)

func (sR *SqlJsonRetrieve) Connect(
	address string,
	dbName string,
	username string,
	password string,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration) error {
	return sR.Psql.Connect(
		address,
		username,
		password,
		dbName,
		maxOpenConns,
		maxIdleConns,
		connMaxLifetime)
}

func (sR *SqlJsonRetrieve) RetrieveFromDB(
	query string, response interface{}, params ...interface{}) error {
	values, err := sR.Psql.GetValuesFromDB(query, params)
	if err != nil {
		return err
	}
	jsonValues, err := sR.JsonHandler.Marshal(values)
	if err != nil {
		return err
	}
	if err := sR.JsonHandler.Unmarshal(jsonValues, response); err != nil {
		return err
	}
	return nil
}

func (sR *SqlJsonRetrieve) WriteInDB(query string, args ...interface{}) (int64, error) {
	return sR.Psql.WriteInDB(query, args)
}

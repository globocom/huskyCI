package db

import (
	"time"
)

// Connect will establish and a connection with Postgres DB based on the received arguments
func (sR *SQLJSONRetrieve) Connect(
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

// RetrieveFromDB will get values from Postgres DB through a SELECT query. The response will
// be returned in response interface. It needs to be a pointer to the struct that will map the
// returned values. The params passed are the ones that will be listed in WHERE argument of the
// query if it exists.
func (sR *SQLJSONRetrieve) RetrieveFromDB(
	query string, response interface{}, params ...interface{}) error {
	values, err := sR.Psql.GetValuesFromDB(query, params)
	if err != nil {
		return err
	}
	jsonValues, err := sR.JSONHandler.Marshal(values)
	if err != nil {
		return err
	}
	if err := sR.JSONHandler.Unmarshal(jsonValues, response); err != nil {
		return err
	}
	return nil
}

// WriteInDB will call Postgres with INSERT and UPDATE queries. The values to be inserted or
// updated are passed in args variable.
func (sR *SQLJSONRetrieve) WriteInDB(query string, args ...interface{}) (int64, error) {
	return sR.Psql.WriteInDB(query, args)
}

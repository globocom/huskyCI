package db

import (
	"errors"
	"log"
	"time"
)

// Connect will call Postgres and establish
// a new connection considering the pool of
// connections configured in the enviroment.
func (sqlConfig *SQLConfig) Connect(address string,
	username string,
	password string,
	dbName string,
	maxOpenConns int,
	maxIdleConns int,
	connLT time.Duration) error {
	if err := sqlConfig.Postgres.ConfigureDB(
		address,
		username,
		password,
		dbName); err != nil {
		return err
	}
	sqlConfig.Postgres.ConfigurePool(maxOpenConns, maxIdleConns, connLT)
	return nil
}

// CloseDB will call Postgres and finish
// its connection.
func (sqlConfig *SQLConfig) CloseDB() error {
	return sqlConfig.Postgres.CloseDB()
}

// GetValuesFromDB will call Postgres through
// SELECT query passed as an argument and return
// all data found in the query. The returned struct
// is an array of map. Each element of the array is
// a row of the returned query. The key is the
// name of column and the value is the data stored in
// the key. If no rows are found in the query, an error
// will be dropped stating that no data were found.
func (sqlConfig *SQLConfig) GetValuesFromDB(query string,
	args ...interface{}) ([]map[string]interface{}, error) {
	err := sqlConfig.Postgres.ConfigureQuery(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := sqlConfig.Postgres.CloseRows()
		if err != nil {
			log.Fatal(err)
		}
	}()
	columns, err := sqlConfig.Postgres.GetColumns()
	if err != nil {
		return nil, err
	}
	results := make([]map[string]interface{}, 0)
	for sqlConfig.Postgres.HasNextRow() {
		rowPointers := generateRowPointers(len(columns))
		if err = sqlConfig.Postgres.ScanRow(rowPointers...); err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		for i, colName := range columns {
			rowVal := rowPointers[i].(*interface{})
			m[colName] = *rowVal
		}
		results = append(results, m)
	}
	if err = sqlConfig.Postgres.GetRowsErr(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("No data found")
	}
	return results, nil
}

// WriteInDB will make Insert and Update queries
// to DB and return the number of rows affected
// during query process.
func (sqlConfig *SQLConfig) WriteInDB(query string, args ...interface{}) (int64, error) {
	err := sqlConfig.Postgres.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return sqlConfig.Postgres.GetRowsAffected()
}

// generateRowPointers returns a slice of interfaces. Each
// element has a memory location related to an interface
// type. So, it will return a slice of pointers of interface
// type.
func generateRowPointers(numPointers int) []interface{} {
	rowResults := make([]interface{}, numPointers)
	rowPointers := make([]interface{}, numPointers)
	for i := range rowResults {
		rowPointers[i] = &rowResults[i]
	}
	return rowPointers
}

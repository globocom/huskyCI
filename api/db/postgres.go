package db

import (
	"errors"
	"time"

	"github.com/globocom/huskyCI/api/types"
)

func (pR *PostgresRequests) ConnectDB(
	address string,
	dbName string,
	username string,
	password string,
	timeout time.Duration,
	poolLimit int,
	port int,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration) error {
	if err := pR.Psql.Connect(
		address,
		username,
		password,
		dbName,
		maxOpenConns,
		maxIdleConns,
		connMaxLifetime); err != nil {
		return err
	}
	return nil
}

func (pR *PostgresRequests) FindOneDBRepository(
	mapParams map[string]interface{}) (types.Repository, error) {
	repositoryResponse := types.Repository{}
	repository, ok := mapParams["repositoryURL"]
	if !ok {
		return repositoryResponse, errors.New("Could not find repository URL")
	}
	myQuery := `SELECT 
					repositoryURL,
					repositoryBranch,
					createdAt
				FROM
					repository
				WHERE
					repositoryURL = $1`

	values, err := pR.Psql.GetValuesFromDB(myQuery, repository)
	if err != nil {
		return repositoryResponse, err
	}
	if len(values) != 1 {
		return repositoryResponse, errors.New("Data returned in a not expected format")
	}
	jsonValues, err := pR.JsonHandler.Marshal(values[0])
	if err != nil {
		return repositoryResponse, err
	}
	if err := pR.JsonHandler.Unmarshal(jsonValues, &repositoryResponse); err != nil {
		return repositoryResponse, err
	}
	return repositoryResponse, nil
}

func (pR *PostgresRequests) FindOneDBSecurityTest(
	mapParams map[string]interface{}) (types.SecurityTest, error) {
	securityResponse := types.SecurityTest{}
	securityTest, ok := mapParams["name"]
	if !ok {
		return securityResponse, errors.New("Could not find securityTest name field")
	}
	myQuery := `SELECT 
					name,
					image,
					imageTag,
					cmd,
					type,
					language,
					default,
					timeOutSeconds
				FROM
					securityTest
				WHERE
					name = $1`

	values, err := pR.Psql.GetValuesFromDB(myQuery, securityTest)
	if err != nil {
		return securityResponse, err
	}
	if len(values) != 1 {
		return securityResponse, errors.New("Data returned in a not expected format")
	}
	jsonValues, err := pR.JsonHandler.Marshal(values[0])
	if err != nil {
		return securityResponse, err
	}
	if err := pR.JsonHandler.Unmarshal(jsonValues, &securityResponse); err != nil {
		return securityResponse, err
	}
	return securityResponse, nil
}

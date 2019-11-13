package db

import (
	"time"

	postgres "github.com/globocom/huskyCI/api/db/postgres"
	"github.com/globocom/huskyCI/api/types"
)

// Requests defines all functions
// that make interactions with the
// database. Based on this, any kind
// of new database support can be done
// implementing Requests.
type Requests interface {
	ConnectDB(address string, dbName string, username string, password string, timeout time.Duration, poolLimit int, port int, maxOpenConns int, maxIdleConns int, connMaxLifetime time.Duration) error
	FindOneDBRepository(mapParams map[string]interface{}) (types.Repository, error)
	FindOneDBSecurityTest(mapParams map[string]interface{}) (types.SecurityTest, error)
	FindOneDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error)
	FindOneDBUser(mapParams map[string]interface{}) (types.User, error)
	FindOneDBAccessToken(mapParams map[string]interface{}) (types.DBToken, error)
	FindAllDBRepository(mapParams map[string]interface{}) ([]types.Repository, error)
	FindAllDBSecurityTest(mapParams map[string]interface{}) ([]types.SecurityTest, error)
	FindAllDBAnalysis(mapParams map[string]interface{}) ([]types.Analysis, error)
	InsertDBRepository(repository types.Repository) error
	InsertDBSecurityTest(securityTest types.SecurityTest) error
	InsertDBAnalysis(analysis types.Analysis) error
	InsertDBUser(user types.User) error
	InsertDBAccessToken(accessToken types.DBToken) error
	UpdateOneDBRepository(mapParams, updateQuery map[string]interface{}) error
	UpsertOneDBSecurityTest(mapParams map[string]interface{}, updatedSecurityTest types.SecurityTest) (interface{}, error)
	UpdateOneDBAnalysis(mapParams map[string]interface{}, updatedAnalysis map[string]interface{}) error
	UpdateOneDBUser(mapParams map[string]interface{}, updatedUser types.User) error
	UpdateOneDBAnalysisContainer(mapParams, updateQuery map[string]interface{}) error
	UpdateOneDBAccessToken(mapParams map[string]interface{}, updatedAccessToken types.DBToken) error
	GetMetricByType(metricType string, queryStringParams map[string][]string) (interface{}, error)
}

// MongoRequests implements Requests
// for Mongo, a non-relational DB.
type MongoRequests struct{}

// JSON interface defines the functions that will threat data
// to be transformed to JSON or a JSON that will be mapped in
// a struct.
type JSON interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// DataGenerator defines the functions that will interact directly
// with DB functions. It abstracts the database functions
type DataGenerator interface {
	Connect(address string, username string, password string, dbName string, maxOpenConns int, maxIdleConns int, connLT time.Duration) error
	RetrieveFromDB(query string, response interface{}, arrayColumns []string, params ...interface{}) error
	WriteInDB(query string, args ...interface{}) (int64, error)
	PqArray(values []string) interface{}
}

// JSONCaller implements JSON interface calling functions
// from encoding/json package.
type JSONCaller struct{}

// SQLJSONRetrieve implements DataGenerator that will interact with
// the Postgres functions. This struct will DB data in JSON format.
type SQLJSONRetrieve struct {
	Psql        postgres.SQLGen
	JSONHandler JSON
}

// PostgresRequests implements Requests
// for Postgres, a relational DB.
type PostgresRequests struct {
	DataRetriever DataGenerator
	JSONHandler   JSON
}

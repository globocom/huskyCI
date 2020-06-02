package database

import "github.com/globocom/huskyCI/api/analysis"

// DBSession is the interface that holds the database session
type DBSession interface {
	Ping() error
	Close() error
	InsertAnalysis(analysis *analysis.Analysis) error
	// FindOneDBRepository(mapParams map[string]interface{}) (types.Repository, error)
	// FindOneDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error)
	// FindOneDBUser(mapParams map[string]interface{}) (types.User, error)
	// FindOneDBAccessToken(mapParams map[string]interface{}) (types.DBToken, error)
	// FindAllDBRepository(mapParams map[string]interface{}) ([]types.Repository, error)
	// FindAllDBSecurityTest(mapParams map[string]interface{}) ([]types.SecurityTest, error)
	// FindAllDBAnalysis(mapParams map[string]interface{}) ([]types.Analysis, error)
	// InsertDBRepository(repository types.Repository) error
	// InsertDBSecurityTest(securityTest types.SecurityTest) error
	// InsertDBUser(user types.User) error
	// InsertDBAccessToken(accessToken types.DBToken) error
	// UpdateOneDBRepository(mapParams, updateQuery map[string]interface{}) error
	// UpsertOneDBSecurityTest(mapParams map[string]interface{}, updatedSecurityTest types.SecurityTest) (interface{}, error)
	// UpdateOneDBAnalysis(mapParams map[string]interface{}, updatedAnalysis map[string]interface{}) error
	// UpdateOneDBUser(mapParams map[string]interface{}, updatedUser types.User) error
	// UpdateOneDBAnalysisContainer(mapParams, updateQuery map[string]interface{}) error
	// UpdateOneDBAccessToken(mapParams map[string]interface{}, updatedAccessToken types.DBToken) error
	// GetMetricByType(metricType string, queryStringParams map[string][]string) (interface{}, error)
}

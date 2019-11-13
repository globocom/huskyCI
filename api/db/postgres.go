package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/types"
)

// ConnectDB will call Connect function
// and try to establish a connection with
// Postgres.
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
	return pR.DataRetriever.Connect(
		address,
		username,
		password,
		dbName,
		maxOpenConns,
		maxIdleConns,
		connMaxLifetime)
}

// FindOneDBRepository checks if a given repository is present into repository table.
func (pR *PostgresRequests) FindOneDBRepository(
	mapParams map[string]interface{}) (types.Repository, error) {
	repositoryResponse := []types.Repository{}
	repository, ok := mapParams["repositoryURL"]
	if !ok {
		return types.Repository{}, errors.New("Could not find repository URL")
	}
	myQuery := `SELECT 
					"repositoryURL",
					"repositoryBranch",
					"createdAt"
				FROM
					repository
				WHERE
					"repositoryURL" = $1`

	if err := pR.DataRetriever.RetrieveFromDB(
		myQuery, &repositoryResponse, []string{}, repository); err != nil {
		return types.Repository{}, err
	}
	return repositoryResponse[0], nil
}

// FindOneDBSecurityTest checks if a given securityTest is present into securityTest table.
func (pR *PostgresRequests) FindOneDBSecurityTest(
	mapParams map[string]interface{}) (types.SecurityTest, error) {
	securityResponse := []types.SecurityTest{}
	securityTest, ok := mapParams["name"]
	if !ok {
		return types.SecurityTest{}, errors.New("Could not find securityTest name field")
	}
	myQuery := `SELECT
					name,
					image,
					"imageTag",
					cmd,
					type,
					language,
					default,
					"timeOutSeconds"
				FROM
					"securityTest"
				WHERE
					name = $1`
	if err := pR.DataRetriever.RetrieveFromDB(
		myQuery, &securityResponse, []string{}, securityTest); err != nil {
		return types.SecurityTest{}, err
	}
	return securityResponse[0], nil
}

// FindOneDBAnalysis checks if a given analysis is present into analysis table.
func (pR *PostgresRequests) FindOneDBAnalysis(
	mapParams map[string]interface{}) (types.Analysis, error) {
	analysisResponse := []types.Analysis{}
	analysis, ok := mapParams["RID"]
	if !ok {
		return types.Analysis{}, errors.New("Could not find RID field")
	}
	myQuery := `SELECT
					"RID",
					"repositoryURL",
					"repositoryBranch",
					"commitAuthors",
					"status",
					"result",
					"errorFound",
					"containers",
					"startedAt",
					"finishedAt",
					"codes",
					"huskyciresults"
				FROM
					analysis
				WHERE
					RID = $1`

	if err := pR.DataRetriever.RetrieveFromDB(
		myQuery, &analysisResponse, []string{"commitAuthors"}, analysis); err != nil {
		return types.Analysis{}, err
	}
	return analysisResponse[0], nil
}

// FindOneDBUser checks if a given user is present into user table.
func (pR *PostgresRequests) FindOneDBUser(
	mapParams map[string]interface{}) (types.User, error) {
	userResponse := []types.User{}
	user, ok := mapParams["username"]
	if !ok {
		return types.User{}, errors.New("Could not find user in DB")
	}
	myQuery := `SELECT
					username,
					password,
					salt,
					interations,
					keylen,
					hashfunction
				FROM
					user
				WHERE
					username = $1`

	if err := pR.DataRetriever.RetrieveFromDB(
		myQuery, &userResponse, []string{}, user); err != nil {
		return types.User{}, err
	}
	return userResponse[0], nil
}

// FindOneDBAccessToken checks if a given accessToken exists in accessToken table.
func (pR *PostgresRequests) FindOneDBAccessToken(
	mapParams map[string]interface{}) (types.DBToken, error) {
	tokenResponse := []types.DBToken{}
	token, ok := mapParams["uuid"]
	if !ok {
		return types.DBToken{}, errors.New("Could not find uuid parameter")
	}
	myQuery := `SELECT
					huskytoken,
					"repositoryURL",
					"isValid",
					"createdAt",
					salt,
					uuid
				FROM
					"accessToken"
				WHERE
					uuid = $1`
	if err := pR.DataRetriever.RetrieveFromDB(
		myQuery, &tokenResponse, []string{}, token); err != nil {
		return types.DBToken{}, err
	}
	return tokenResponse[0], nil
}

// FindAllDBRepository returns all Repository of a given query present into repository table.
func (pR *PostgresRequests) FindAllDBRepository(
	mapParams map[string]interface{}) ([]types.Repository, error) {
	repositoryResponse := []types.Repository{}
	query, params := ConfigureQuery(`SELECT * FROM repository`, mapParams)
	if err := pR.DataRetriever.RetrieveFromDB(
		query, &repositoryResponse, []string{}, params); err != nil {
		return repositoryResponse, err
	}
	return repositoryResponse, nil
}

// FindAllDBSecurityTest returns all SecurityTests of a given query present
// into security Test table.
func (pR *PostgresRequests) FindAllDBSecurityTest(
	mapParams map[string]interface{}) ([]types.SecurityTest, error) {
	securityResponse := []types.SecurityTest{}
	query, params := ConfigureQuery(`SELECT * FROM securityTest`, mapParams)
	if err := pR.DataRetriever.RetrieveFromDB(
		query, &securityResponse, []string{}, params); err != nil {
		return securityResponse, err
	}
	return securityResponse, nil
}

// FindAllDBAnalysis returns all Analysis of a given query present into analysis table.
func (pR *PostgresRequests) FindAllDBAnalysis(
	mapParams map[string]interface{}) ([]types.Analysis, error) {
	analysisResponse := []types.Analysis{}
	query, params := ConfigureQuery(`SELECT * FROM analysis`, mapParams)
	if err := pR.DataRetriever.RetrieveFromDB(
		query, &analysisResponse, []string{}, params); err != nil {
		return analysisResponse, err
	}
	return analysisResponse, nil
}

// InsertDBRepository inserts a new repository into repository table.
func (pR *PostgresRequests) InsertDBRepository(repository types.Repository) error {
	if (types.Repository{}) == repository {
		return errors.New("Empty repository data")
	}
	repositoryMap := map[string]interface{}{
		"repositoryURL": repository.URL,
		"createdAt":     repository.CreatedAt,
	}
	finalQuery, values := ConfigureInsertQuery(
		`INSERT into repository`, repositoryMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was inserted")
	}
	return nil
}

// InsertDBSecurityTest inserts a new securityTest into securityTest table.
func (pR *PostgresRequests) InsertDBSecurityTest(securityTest types.SecurityTest) error {
	if (types.SecurityTest{}) == securityTest {
		return errors.New("Empty SecurityTest data")
	}
	securityMap := map[string]interface{}{
		"name":           securityTest.Name,
		"image":          securityTest.Image,
		"cmd":            securityTest.Cmd,
		"language":       securityTest.Language,
		"type":           securityTest.Type,
		"default":        securityTest.Default,
		"timeOutSeconds": securityTest.TimeOutInSeconds,
	}
	finalQuery, values := ConfigureInsertQuery(
		`INSERT into "securityTest"`, securityMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was inserted")
	}
	return nil
}

// InsertDBAnalysis inserts a new analysis into analysis table.
func (pR *PostgresRequests) InsertDBAnalysis(analysis types.Analysis) error {
	if analysis.URL == "" {
		return errors.New("Empty Analysis data")
	}
	analysisMap := map[string]interface{}{
		"RID":              analysis.RID,
		"repositoryURL":    analysis.URL,
		"repositoryBranch": analysis.Branch,
		"status":           analysis.Status,
		"result":           analysis.Result,
		"containers":       analysis.Containers,
		"startedAt":        analysis.StartedAt,
	}
	analysisMap, err := pR.ConfigureAnalysisData(analysisMap)
	if err != nil {
		return err
	}
	finalQuery, values := ConfigureInsertQuery(
		`INSERT into analysis`, analysisMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was inserted")
	}
	return nil
}

// InsertDBUser inserts a new user into user table.
func (pR *PostgresRequests) InsertDBUser(user types.User) error {
	if (types.User{}) == user {
		return errors.New("Empty User data")
	}
	userMap := map[string]interface{}{
		"username":     user.Username,
		"password":     user.Password,
		"salt":         user.Salt,
		"iterations":   user.Iterations,
		"keylen":       user.KeyLen,
		"hashfunction": user.HashFunction,
	}
	finalQuery, values := ConfigureInsertQuery(
		`INSERT into user`, userMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was inserted")
	}
	return nil
}

// InsertDBAccessToken inserts a new access into accessToken table.
func (pR *PostgresRequests) InsertDBAccessToken(accessToken types.DBToken) error {
	if (types.DBToken{}) == accessToken {
		return errors.New("Empty DBToken data")
	}
	accessTokenMap := map[string]interface{}{
		"huskytoken":    accessToken.HuskyToken,
		"repositoryURL": accessToken.URL,
		"isValid":       accessToken.IsValid,
		"createdAt":     accessToken.CreatedAt,
		"salt":          accessToken.Salt,
		"uuid":          accessToken.UUID,
	}
	finalQuery, values := ConfigureInsertQuery(
		`INSERT into "accessToken"`, accessTokenMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was inserted")
	}
	return nil
}

// UpdateOneDBRepository checks if a given repository is present into repository table
// and update it.
func (pR *PostgresRequests) UpdateOneDBRepository(
	mapParams, updateQuery map[string]interface{}) error {
	if len(updateQuery) == 0 {
		return errors.New("Empty fields to be updated")
	}
	if len(mapParams) == 0 {
		return errors.New("Empty fields to search")
	}
	finalQuery, values := ConfigureUpdateQuery(
		`UPDATE repository`, mapParams, updateQuery)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was updated")
	}
	return nil
}

// UpsertOneDBSecurityTest checks if a given securityTest is present into securityTest
// and update it. If not, it will insert a new entry.
func (pR *PostgresRequests) UpsertOneDBSecurityTest(
	mapParams map[string]interface{}, updatedSecurityTest types.SecurityTest) (interface{}, error) {
	if (types.SecurityTest{}) == updatedSecurityTest {
		return nil, errors.New("Empty fields to be updated")
	}
	if len(mapParams) == 0 {
		return nil, errors.New("Empty fields to search")
	}
	updatedSecurityMap := map[string]interface{}{
		"name":           updatedSecurityTest.Name,
		"image":          updatedSecurityTest.Image,
		"imageTag":       updatedSecurityTest.ImageTag,
		"cmd":            updatedSecurityTest.Cmd,
		"type":           updatedSecurityTest.Type,
		"language":       updatedSecurityTest.Language,
		"default":        updatedSecurityTest.Default,
		"timeOutSeconds": updatedSecurityTest.TimeOutInSeconds,
	}
	finalQuery, values := ConfigureUpsertQuery(
		`INSERT into "securityTest"`, mapParams, updatedSecurityMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return nil, err
	}
	if rowsAff == int64(0) {
		return nil, errors.New("No data was updated")
	}
	return rowsAff, nil
}

// ConfigureUpsertQuery will receive a partial query and mount the final query
// CONFLICT statement so it will allow Postgres make an Upsert in the entry
// based on the conflicted columns passed in this statement. An UPDATE query is build
// with the related values to be updated in case of a conflict.
func ConfigureUpsertQuery(
	query string, searchValues, newValues map[string]interface{}) (string, []interface{}) {
	insertQuery, values := ConfigureInsertQuery(query, newValues)
	conflictQuery := ""
	i := 1
	for k := range searchValues {
		if !strings.Contains(conflictQuery, "CONFLICT") {
			conflictQuery = `ON CONFLICT (`
		}
		if i == len(searchValues) {
			conflictQuery = fmt.Sprintf(`%s"%s")`, conflictQuery, k)
		} else {
			conflictQuery = fmt.Sprintf(`%s"%s", `, conflictQuery, k)
		}
		i++
	}
	updateQuery := ""
	for k := range newValues {
		if !strings.Contains(updateQuery, "UPDATE") {
			updateQuery = `DO UPDATE SET`
		}
		if strings.Contains(updateQuery, "=") {
			updateQuery = fmt.Sprintf("%s,", updateQuery)
		}
		updateQuery = fmt.Sprintf(`%s "%s" = EXCLUDED."%s"`, updateQuery, k, k)
	}
	if conflictQuery != "" {
		insertQuery = fmt.Sprintf(`%s %s`, insertQuery, conflictQuery)
	}
	if updateQuery != "" {
		insertQuery = fmt.Sprintf(`%s %s`, insertQuery, updateQuery)
	}
	return insertQuery, values
}

// UpdateOneDBAnalysis checks if a given analysis is present into analysis table and update it.
func (pR *PostgresRequests) UpdateOneDBAnalysis(
	mapParams map[string]interface{}, updatedAnalysis map[string]interface{}) error {
	if len(updatedAnalysis) == 0 {
		return errors.New("Empty fields to be updated")
	}
	if len(mapParams) == 0 {
		return errors.New("Empty fields to search")
	}
	// Convert commitAuthors to a valid type for psql
	// understand that it is an array.
	updatedAnalysis, err := pR.ConfigureAnalysisData(updatedAnalysis)
	if err != nil {
		return err
	}
	finalQuery, values := ConfigureUpdateQuery(
		`UPDATE analysis`, mapParams, updatedAnalysis)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was updated")
	}
	return nil
}

// UpdateOneDBUser checks if a given user is present into user table and update it.
func (pR *PostgresRequests) UpdateOneDBUser(
	mapParams map[string]interface{}, updatedUser types.User) error {
	if (types.User{}) == updatedUser {
		return errors.New("Empty fields to be updated")
	}
	if len(mapParams) == 0 {
		return errors.New("Empty fields to search")
	}
	updatedUserMap := map[string]interface{}{
		"username":     updatedUser.Username,
		"password":     updatedUser.Password,
		"salt":         updatedUser.Salt,
		"iterations":   updatedUser.Iterations,
		"keylen":       updatedUser.KeyLen,
		"hashfunction": updatedUser.HashFunction,
	}
	finalQuery, values := ConfigureUpdateQuery(
		`UPDATE user`, mapParams, updatedUserMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was updated")
	}
	return nil
}

// UpdateOneDBAnalysisContainer checks if a given analysis is present into analysis table
// and update the container associated in it.
func (pR *PostgresRequests) UpdateOneDBAnalysisContainer(
	mapParams, updateQuery map[string]interface{}) error {
	if len(updateQuery) == 0 {
		return errors.New("Empty fields to be updated")
	}
	if len(mapParams) == 0 {
		return errors.New("Empty fields to search")
	}
	updateQuery, err := pR.ConfigureAnalysisData(updateQuery)
	if err != nil {
		return err
	}
	finalQuery, values := ConfigureUpdateQuery(
		`UPDATE analysis`, mapParams, updateQuery)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was updated")
	}
	return nil
}

// UpdateOneDBAccessToken checks if a given access token is present into accessToken and update it.
func (pR *PostgresRequests) UpdateOneDBAccessToken(
	mapParams map[string]interface{}, updatedAccessToken types.DBToken) error {
	if (types.DBToken{}) == updatedAccessToken {
		return errors.New("Empty fields to be updated")
	}
	if len(mapParams) == 0 {
		return errors.New("Empty fields to search")
	}
	updatedAccessTokenMap := map[string]interface{}{
		"huskytoken":    updatedAccessToken.HuskyToken,
		"repositoryURL": updatedAccessToken.URL,
		"isValid":       updatedAccessToken.IsValid,
		"createdAt":     updatedAccessToken.CreatedAt,
		"salt":          updatedAccessToken.Salt,
		"uuid":          updatedAccessToken.UUID,
	}
	finalQuery, values := ConfigureUpdateQuery(
		`UPDATE "accessToken"`, mapParams, updatedAccessTokenMap)
	rowsAff, err := pR.DataRetriever.WriteInDB(finalQuery, values)
	if err != nil {
		return err
	}
	if rowsAff == int64(0) {
		return errors.New("No data was updated")
	}
	return nil
}

// GetMetricByType returns data about the metric received
func (pR *PostgresRequests) GetMetricByType(
	metricType string, queryStringParams map[string][]string) (interface{}, error) {
	// TODO: Need to know how to generate the same statistics
	// as on Mongo
	return nil, errors.New("Function not supported")
}

// ConfigureUpdateQuery will receive a partial update query and mount the final query with
// all data to be set and the search parameters related to the row to be changed.
func ConfigureUpdateQuery(
	query string, searchValues, newValues map[string]interface{}) (string, []interface{}) {
	valuesQuery := ""
	searchQuery := ""
	values := make([]interface{}, 0)
	i := 1
	for k, v := range searchValues {
		if !strings.Contains(searchQuery, "WHERE") {
			searchQuery = `WHERE`
		}
		if strings.Contains(searchQuery, "=") {
			searchQuery = fmt.Sprintf("%s AND", searchQuery)
		}
		searchQuery = fmt.Sprintf(`%s "%s" = $%d`, searchQuery, k, i)
		i++
		values = append(values, v)
	}
	for k, v := range newValues {
		if !strings.Contains(valuesQuery, "SET") {
			valuesQuery = `SET`
		}
		if strings.Contains(valuesQuery, "=") {
			valuesQuery = fmt.Sprintf("%s,", valuesQuery)
		}
		valuesQuery = fmt.Sprintf(`%s "%s" = $%d`, valuesQuery, k, i)
		i++
		values = append(values, v)
	}
	if valuesQuery != "" {
		query = fmt.Sprintf("%s %s", query, valuesQuery)
	}
	if searchQuery != "" {
		query = fmt.Sprintf("%s %s", query, searchQuery)
	}
	return query, values
}

// ConfigureInsertQuery will receive a partial query and mount the final query with all data
// to be inserted in the table with its related fields.
func ConfigureInsertQuery(query string, params map[string]interface{}) (string, []interface{}) {
	values := make([]interface{}, 0)
	i := 1
	argsQuery := `(`
	valuesQuery := `VALUES (`
	for k, v := range params {
		if i == len(params) {
			argsQuery = fmt.Sprintf(`%s"%s")`, argsQuery, k)
			valuesQuery = fmt.Sprintf("%s$%d)", valuesQuery, i)
		} else {
			argsQuery = fmt.Sprintf(`%s"%s", `, argsQuery, k)
			valuesQuery = fmt.Sprintf("%s$%d, ", valuesQuery, i)
		}
		values = append(values, v)
		i++
	}
	query = fmt.Sprintf("%s %s %s", query, argsQuery, valuesQuery)
	return query, values
}

// ConfigureQuery will receive a partial search query and will mount the final query with the
// search if it exists.
func ConfigureQuery(query string, params map[string]interface{}) (string, []interface{}) {
	if len(params) != 0 {
		query = fmt.Sprintf("%s WHERE", query)
	}
	values := make([]interface{}, 0)
	i := 1
	for k, v := range params {
		if strings.Contains(query, "=") {
			query = fmt.Sprintf("%s AND", query)
		}
		query = fmt.Sprintf(`%s "%s" = $%d`, query, k, i)
		values = append(values, v)
		i++
	}
	return query, values
}

func (pR *PostgresRequests) ConfigureAnalysisData(
	updatedAnalysis map[string]interface{}) (map[string]interface{}, error) {
	if authors, ok := updatedAnalysis["commitAuthors"].([]string); ok {
		updatedAnalysis["commitAuthors"] = pR.DataRetriever.PqArray(authors)
	}
	if containers, ok := updatedAnalysis["containers"].([]types.Container); ok {
		containerJson, err := pR.JSONHandler.Marshal(containers)
		if err != nil {
			return updatedAnalysis, err
		}
		updatedAnalysis["containers"] = containerJson
	}
	if huskyciresults, ok := updatedAnalysis["huskyciresults"].(types.HuskyCIResults); ok {
		huskyJson, err := pR.JSONHandler.Marshal(huskyciresults)
		if err != nil {
			return updatedAnalysis, err
		}
		updatedAnalysis["huskyciresults"] = huskyJson
	}
	if myCodes, ok := updatedAnalysis["codes"].([]types.Code); ok {
		codeJson, err := pR.JSONHandler.Marshal(myCodes)
		if err != nil {
			return updatedAnalysis, err
		}
		updatedAnalysis["codes"] = codeJson
	}
	return updatedAnalysis, nil
}

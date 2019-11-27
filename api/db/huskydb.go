// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"time"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
)

// ConnectDB will call Connect function
// and return a nil error if connection
// with MongoDB was succeeded.
func (mR *MongoRequests) ConnectDB(address string, dbName string,
	username string,
	password string,
	timeout time.Duration,
	poolLimit int,
	port int,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration) error {
	return mongoHuskyCI.Connect(
		address,
		dbName,
		username,
		password,
		poolLimit,
		port,
		timeout)
}

// FindOneDBRepository checks if a given repository is present into RepositoryCollection.
func (mR *MongoRequests) FindOneDBRepository(mapParams map[string]interface{}) (types.Repository, error) {
	repositoryResponse := types.Repository{}
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	err := mongoHuskyCI.Conn.SearchOne(repositoryFinalQuery, nil, mongoHuskyCI.RepositoryCollection, &repositoryResponse)
	return repositoryResponse, err
}

// FindOneDBSecurityTest checks if a given securityTest is present into SecurityTestCollection.
func (mR *MongoRequests) FindOneDBSecurityTest(mapParams map[string]interface{}) (types.SecurityTest, error) {
	securityTestResponse := types.SecurityTest{}
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	err := mongoHuskyCI.Conn.SearchOne(securityTestFinalQuery, nil, mongoHuskyCI.SecurityTestCollection, &securityTestResponse)
	return securityTestResponse, err
}

// FindOneDBAnalysis checks if a given analysis is present into AnalysisCollection.
func (mR *MongoRequests) FindOneDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error) {
	analysisResponse := types.Analysis{}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}

	err := mongoHuskyCI.Conn.SearchOne(analysisFinalQuery, nil, mongoHuskyCI.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// FindOneDBUser checks if a given user is present into UserCollection.
func (mR *MongoRequests) FindOneDBUser(mapParams map[string]interface{}) (types.User, error) {
	userResponse := types.User{}
	userQuery := []bson.M{}
	for k, v := range mapParams {
		userQuery = append(userQuery, bson.M{k: v})
	}
	userFinalQuery := bson.M{"$and": userQuery}
	err := mongoHuskyCI.Conn.SearchOne(userFinalQuery, nil, mongoHuskyCI.UserCollection, &userResponse)
	return userResponse, err
}

// FindOneDBAccessToken checks if a given accessToken exists in AccessTokenCollection.
func (mR *MongoRequests) FindOneDBAccessToken(mapParams map[string]interface{}) (types.DBToken, error) {
	aTokenResponse := types.DBToken{}
	aTokenQuery := []bson.M{}
	for k, v := range mapParams {
		aTokenQuery = append(aTokenQuery, bson.M{k: v})
	}
	aTokenFinalQuery := bson.M{"$and": aTokenQuery}
	err := mongoHuskyCI.Conn.SearchOne(aTokenFinalQuery, nil, mongoHuskyCI.AccessTokenCollection, &aTokenResponse)
	return aTokenResponse, err
}

// FindAllDBRepository returns all Repository of a given query present into RepositoryCollection.
func (mR *MongoRequests) FindAllDBRepository(mapParams map[string]interface{}) ([]types.Repository, error) {
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	repositoryResponse := []types.Repository{}
	err := mongoHuskyCI.Conn.Search(repositoryFinalQuery, nil, mongoHuskyCI.RepositoryCollection, &repositoryResponse)
	return repositoryResponse, err
}

// FindAllDBSecurityTest returns all SecurityTests of a given query present into SecurityTestCollection.
func (mR *MongoRequests) FindAllDBSecurityTest(mapParams map[string]interface{}) ([]types.SecurityTest, error) {
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	securityTestResponse := []types.SecurityTest{}
	err := mongoHuskyCI.Conn.Search(securityTestFinalQuery, nil, mongoHuskyCI.SecurityTestCollection, &securityTestResponse)
	return securityTestResponse, err
}

// FindAllDBAnalysis returns all Analysis of a given query present into AnalysisCollection.
func (mR *MongoRequests) FindAllDBAnalysis(mapParams map[string]interface{}) ([]types.Analysis, error) {
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	analysisResponse := []types.Analysis{}
	err := mongoHuskyCI.Conn.Search(analysisFinalQuery, nil, mongoHuskyCI.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// InsertDBRepository inserts a new repository into RepositoryCollection.
func (mR *MongoRequests) InsertDBRepository(repository types.Repository) error {
	newRepository := bson.M{
		"repositoryURL": repository.URL,
		"createdAt":     repository.CreatedAt,
	}
	err := mongoHuskyCI.Conn.Insert(newRepository, mongoHuskyCI.RepositoryCollection)
	return err
}

// InsertDBSecurityTest inserts a new securityTest into SecurityTestCollection.
func (mR *MongoRequests) InsertDBSecurityTest(securityTest types.SecurityTest) error {
	newSecurityTest := bson.M{
		"name":           securityTest.Name,
		"image":          securityTest.Image,
		"cmd":            securityTest.Cmd,
		"language":       securityTest.Language,
		"type":           securityTest.Type,
		"default":        securityTest.Default,
		"timeOutSeconds": securityTest.TimeOutInSeconds,
	}
	err := mongoHuskyCI.Conn.Insert(newSecurityTest, mongoHuskyCI.SecurityTestCollection)
	return err
}

// InsertDBAnalysis inserts a new analysis into AnalysisCollection.
func (mR *MongoRequests) InsertDBAnalysis(analysis types.Analysis) error {
	newAnalysis := bson.M{
		"RID":              analysis.RID,
		"repositoryURL":    analysis.URL,
		"repositoryBranch": analysis.Branch,
		"status":           analysis.Status,
		"result":           analysis.Result,
		"containers":       analysis.Containers,
		"startedAt":        analysis.StartedAt,
	}
	err := mongoHuskyCI.Conn.Insert(newAnalysis, mongoHuskyCI.AnalysisCollection)
	return err
}

// InsertDBUser inserts a new user into UserCollection.
func (mR *MongoRequests) InsertDBUser(user types.User) error {
	newUser := bson.M{
		"username":     user.Username,
		"password":     user.Password,
		"salt":         user.Salt,
		"iterations":   user.Iterations,
		"keylen":       user.KeyLen,
		"hashfunction": user.HashFunction,
	}
	err := mongoHuskyCI.Conn.Insert(newUser, mongoHuskyCI.UserCollection)
	return err
}

// InsertDBAccessToken inserts a new access into AccessTokenCollection.
func (mR *MongoRequests) InsertDBAccessToken(accessToken types.DBToken) error {
	newAccessToken := bson.M{
		"huskytoken":    accessToken.HuskyToken,
		"repositoryURL": accessToken.URL,
		"isValid":       accessToken.IsValid,
		"createdAt":     accessToken.CreatedAt,
		"salt":          accessToken.Salt,
		"uuid":          accessToken.UUID,
	}
	err := mongoHuskyCI.Conn.Insert(newAccessToken, mongoHuskyCI.AccessTokenCollection)
	return err
}

// UpdateOneDBRepository checks if a given repository is present into RepositoryCollection and update it.
func (mR *MongoRequests) UpdateOneDBRepository(mapParams, updateQuery map[string]interface{}) error {
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	err := mongoHuskyCI.Conn.Update(repositoryFinalQuery, updateQuery, mongoHuskyCI.RepositoryCollection)
	return err
}

// UpsertOneDBSecurityTest checks if a given securityTest is present into SecurityTestCollection and update it.
func (mR *MongoRequests) UpsertOneDBSecurityTest(mapParams map[string]interface{}, updatedSecurityTest types.SecurityTest) (interface{}, error) {
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	changeInfo, err := mongoHuskyCI.Conn.Upsert(securityTestFinalQuery, updatedSecurityTest, mongoHuskyCI.SecurityTestCollection)
	return changeInfo, err
}

// UpdateOneDBAnalysis checks if a given analysis is present into AnalysisCollection and update it.
func (mR *MongoRequests) UpdateOneDBAnalysis(mapParams map[string]interface{}, updatedAnalysis map[string]interface{}) error {
	updatedQuery := bson.M{
		"$set": updatedAnalysis,
	}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	err := mongoHuskyCI.Conn.Update(analysisFinalQuery, updatedQuery, mongoHuskyCI.AnalysisCollection)
	return err
}

// UpdateOneDBUser checks if a given user is present into UserCollection and update it.
func (mR *MongoRequests) UpdateOneDBUser(mapParams map[string]interface{}, updatedUser types.User) error {
	userQuery := []bson.M{}
	for k, v := range mapParams {
		userQuery = append(userQuery, bson.M{k: v})
	}
	userFinalQuery := bson.M{"$and": userQuery}
	err := mongoHuskyCI.Conn.Update(userFinalQuery, updatedUser, mongoHuskyCI.UserCollection)
	return err
}

// UpdateOneDBAnalysisContainer checks if a given analysis is present into AnalysisCollection and update the container associated in it.
func (mR *MongoRequests) UpdateOneDBAnalysisContainer(mapParams, updateQuery map[string]interface{}) error {
	updatedQuery := bson.M{
		"$set": updateQuery,
	}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	err := mongoHuskyCI.Conn.Update(analysisFinalQuery, updatedQuery, mongoHuskyCI.AnalysisCollection)
	return err
}

// UpdateOneDBAccessToken checks if a given access token is present into AccessTokenCollection and update it.
func (mR *MongoRequests) UpdateOneDBAccessToken(mapParams map[string]interface{}, updatedAccessToken types.DBToken) error {
	aTokenQuery := []bson.M{}
	for k, v := range mapParams {
		aTokenQuery = append(aTokenQuery, bson.M{k: v})
	}
	aTokenFinalQuery := bson.M{"$and": aTokenQuery}
	err := mongoHuskyCI.Conn.Update(aTokenFinalQuery, updatedAccessToken, mongoHuskyCI.AccessTokenCollection)
	return err
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/globocom/huskyCI/api/analysis"
	"github.com/globocom/huskyCI/api/auth"
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

// FindOneDBAnalysis checks if a given analysis is present into AnalysisCollection.
func (mR *MongoRequests) FindOneDBAnalysis(mapParams map[string]interface{}) (analysis.Analysis, error) {
	analysisResponse := analysis.Analysis{}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}

	err := mongoHuskyCI.Conn.SearchOne(analysisFinalQuery, nil, mongoHuskyCI.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// FindOneDBUser checks if a given user is present into UserCollection.
func (mR *MongoRequests) FindOneDBUser(mapParams map[string]interface{}) (auth.User, error) {
	userResponse := auth.User{}
	userQuery := []bson.M{}
	for k, v := range mapParams {
		userQuery = append(userQuery, bson.M{k: v})
	}
	userFinalQuery := bson.M{"$and": userQuery}
	err := mongoHuskyCI.Conn.SearchOne(userFinalQuery, nil, mongoHuskyCI.UserCollection, &userResponse)
	return userResponse, err
}

// FindOneDBAccessToken checks if a given accessToken exists in TokenCollection.
func (mR *MongoRequests) FindOneDBAccessToken(mapParams map[string]interface{}) (types.DBToken, error) {
	aTokenResponse := types.DBToken{}
	aTokenQuery := []bson.M{}
	for k, v := range mapParams {
		aTokenQuery = append(aTokenQuery, bson.M{k: v})
	}
	aTokenFinalQuery := bson.M{"$and": aTokenQuery}
	err := mongoHuskyCI.Conn.SearchOne(aTokenFinalQuery, nil, mongoHuskyCI.TokenCollection, &aTokenResponse)
	return aTokenResponse, err
}

// FindAllDBAnalysis returns all Analysis of a given query present into AnalysisCollection.
func (mR *MongoRequests) FindAllDBAnalysis(mapParams map[string]interface{}) ([]analysis.Analysis, error) {
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	analysisResponse := []analysis.Analysis{}
	err := mongoHuskyCI.Conn.Search(analysisFinalQuery, nil, mongoHuskyCI.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// InsertDBAnalysis inserts a new analysis into AnalysisCollection.
func (mR *MongoRequests) InsertDBAnalysis(analysis analysis.Analysis) error {
	newAnalysis := bson.M{
		"ID":              analysis.ID,
		"repository":      analysis.Repository,
		"result":          analysis.Result,
		"startedAt":       analysis.StartedAt,
		"finishedAt":      analysis.FinishedAt,
		"errorsFound":     analysis.ErrorsFound,
		"vulnerabilities": analysis.Vulnerabilities,
		"securityTests":   analysis.SecurityTests,
	}
	err := mongoHuskyCI.Conn.Insert(newAnalysis, mongoHuskyCI.AnalysisCollection)
	return err
}

// InsertDBUser inserts a new user into UserCollection.
func (mR *MongoRequests) InsertDBUser(user auth.User) error {
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

// InsertDBAccessToken inserts a new access into TokenCollection.
func (mR *MongoRequests) InsertDBAccessToken(accessToken types.DBToken) error {
	newAccessToken := bson.M{
		"huskytoken":    accessToken.HuskyToken,
		"repositoryURL": accessToken.URL,
		"isValid":       accessToken.IsValid,
		"createdAt":     accessToken.CreatedAt,
		"salt":          accessToken.Salt,
		"uuid":          accessToken.UUID,
	}
	err := mongoHuskyCI.Conn.Insert(newAccessToken, mongoHuskyCI.TokenCollection)
	return err
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
func (mR *MongoRequests) UpdateOneDBUser(mapParams map[string]interface{}, updatedUser auth.User) error {
	userQuery := []bson.M{}
	for k, v := range mapParams {
		userQuery = append(userQuery, bson.M{k: v})
	}
	userFinalQuery := bson.M{"$and": userQuery}
	err := mongoHuskyCI.Conn.Update(userFinalQuery, updatedUser, mongoHuskyCI.UserCollection)
	return err
}

// UpdateOneDBAccessToken checks if a given access token is present into TokenCollection and update it.
func (mR *MongoRequests) UpdateOneDBAccessToken(mapParams map[string]interface{}, updatedAccessToken types.DBToken) error {
	aTokenQuery := []bson.M{}
	for k, v := range mapParams {
		aTokenQuery = append(aTokenQuery, bson.M{k: v})
	}
	aTokenFinalQuery := bson.M{"$and": aTokenQuery}
	err := mongoHuskyCI.Conn.Update(aTokenFinalQuery, updatedAccessToken, mongoHuskyCI.TokenCollection)
	return err
}

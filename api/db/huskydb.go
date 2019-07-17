// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"time"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

// FindOneDBRepository checks if a given repository is present into RepositoryCollection.
func FindOneDBRepository(mapParams map[string]interface{}) (types.Repository, error) {
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
func FindOneDBSecurityTest(mapParams map[string]interface{}) (types.SecurityTest, error) {
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
func FindOneDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error) {
	analysisResponse := types.Analysis{}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}

	err := mongoHuskyCI.Conn.SearchOne(analysisFinalQuery, nil, mongoHuskyCI.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// FindAllDBRepository returns all Repository of a given query present into RepositoryCollection.
func FindAllDBRepository(mapParams map[string]interface{}) ([]types.Repository, error) {
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
func FindAllDBSecurityTest(mapParams map[string]interface{}) ([]types.SecurityTest, error) {
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
func FindAllDBAnalysis(mapParams map[string]interface{}) ([]types.Analysis, error) {
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	analysisResponse := []types.Analysis{}
	err := mongoHuskyCI.Conn.Search(analysisFinalQuery, nil, mongoHuskyCI.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// InsertDBRepository inserts a new repository with default securityTests into RepositoryCollection.
func InsertDBRepository(repository types.Repository) error {
	repository.CreatedAt = time.Now()
	securityTestList := []types.SecurityTest{}
	err := errors.New("")
	maxQueryAllowed := 8

	if len(repository.SecurityTestName) == 0 {
		// checking default and Generic securityTests in SecurityTestCollection
		securityTestQuery := map[string]interface{}{"default": true, "language": "Generic"}
		securityTestList, err = FindAllDBSecurityTest(securityTestQuery)
		if err != nil {
			log.Error("InsertDBRepository", "HUSKYDB", 2005, err)
		}
	} else {
		// checking if a given securityTestName matches a securityTest
		repository.SecurityTestName = util.RemoveDuplicates(repository.SecurityTestName)
		if len(repository.SecurityTestName) > maxQueryAllowed {
			for _, securityTestName := range repository.SecurityTestName[:maxQueryAllowed] {
				securityTestQuery := map[string]interface{}{"name": securityTestName}
				securityTestResult, err := FindOneDBSecurityTest(securityTestQuery)
				if err != nil {
					log.Error("InsertDBRepository", "HUSKYDB", 2006, err)
				} else {
					securityTestList = append(securityTestList, securityTestResult)
				}
			}
		} else {
			for _, securityTestName := range repository.SecurityTestName {
				securityTestQuery := map[string]interface{}{"name": securityTestName}
				securityTestResult, err := FindOneDBSecurityTest(securityTestQuery)
				if err != nil {
					log.Error("InsertDBRepository", "HUSKYDB", 2006, err)
				} else {
					securityTestList = append(securityTestList, securityTestResult)
				}
			}
		}
	}

	newRepository := types.Repository{
		URL:            repository.URL,
		Branch:         repository.Branch,
		SecurityTests:  securityTestList,
		CreatedAt:      repository.CreatedAt,
		DeletedAt:      repository.DeletedAt,
		Languages:      repository.Languages,
		InternalDepURL: repository.InternalDepURL,
	}

	err = mongoHuskyCI.Conn.Insert(newRepository, mongoHuskyCI.RepositoryCollection)
	return err
}

// InsertDBSecurityTest inserts a new securityTest into SecurityTestCollection.
func InsertDBSecurityTest(securityTest types.SecurityTest) error {
	newSecurityTest := bson.M{
		"name":           securityTest.Name,
		"image":          securityTest.Image,
		"cmd":            securityTest.Cmd,
		"language":       securityTest.Language,
		"default":        securityTest.Default,
		"timeOutSeconds": securityTest.TimeOutInSeconds,
	}
	err := mongoHuskyCI.Conn.Insert(newSecurityTest, mongoHuskyCI.SecurityTestCollection)
	return err
}

// InsertDBAnalysis inserts a new analysis into AnalysisCollection.
func InsertDBAnalysis(analysis types.Analysis) error {
	newAnalysis := bson.M{
		"RID":            analysis.RID,
		"URL":            analysis.URL,
		"Branch":         analysis.Branch,
		"securityTest":   analysis.SecurityTests,
		"status":         analysis.Status,
		"result":         analysis.Result,
		"containers":     analysis.Containers,
		"startedAt":      analysis.StartedAt,
		"internaldepURL": analysis.InternalDepURL,
	}
	err := mongoHuskyCI.Conn.Insert(newAnalysis, mongoHuskyCI.AnalysisCollection)
	return err
}

// UpdateOneDBRepository checks if a given repository is present into RepositoryCollection and update it.
func UpdateOneDBRepository(mapParams, updateQuery map[string]interface{}) error {
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	err := mongoHuskyCI.Conn.Update(repositoryFinalQuery, updateQuery, mongoHuskyCI.RepositoryCollection)
	return err
}

// UpdateOneDBSecurityTest checks if a given securityTest is present into SecurityTestCollection and update it.
func UpdateOneDBSecurityTest(mapParams map[string]interface{}, updatedSecurityTest types.SecurityTest) error {
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	_, err := mongoHuskyCI.Conn.Upsert(securityTestFinalQuery, updatedSecurityTest, mongoHuskyCI.SecurityTestCollection)
	return err
}

// UpdateOneDBAnalysis checks if a given analysis is present into AnalysisCollection and update it.
func UpdateOneDBAnalysis(mapParams map[string]interface{}, updatedAnalysis types.Analysis) error {
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	err := mongoHuskyCI.Conn.Update(analysisFinalQuery, updatedAnalysis, mongoHuskyCI.AnalysisCollection)
	return err
}

// UpdateOneDBAnalysisContainer checks if a given analysis is present into AnalysisCollection and update the container associated in it.
func UpdateOneDBAnalysisContainer(mapParams, updateQuery map[string]interface{}) error {
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	err := mongoHuskyCI.Conn.Update(analysisFinalQuery, updateQuery, mongoHuskyCI.AnalysisCollection)
	return err
}

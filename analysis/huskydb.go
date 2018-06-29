package analysis

import (
	"fmt"
	"time"

	db "github.com/globocom/husky/db/mongo"
	"github.com/globocom/husky/types"
	"gopkg.in/mgo.v2/bson"
)

// FindOneDBRepository checks if a given repository is present into RepositoryCollection.
func FindOneDBRepository(mapParams map[string]interface{}) (types.Repository, error) {
	session := db.Connect()
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	repositoryResponse := types.Repository{}
	err := session.SearchOne(repositoryFinalQuery, nil, db.RepositoryCollection, &repositoryResponse)
	return repositoryResponse, err
}

// FindOneDBSecurityTest checks if a given securityTest is present into SecurityTestCollection.
func FindOneDBSecurityTest(mapParams map[string]interface{}) (types.SecurityTest, error) {
	session := db.Connect()
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	securityTestResponse := types.SecurityTest{}
	err := session.SearchOne(securityTestFinalQuery, nil, db.SecurityTestCollection, &securityTestResponse)
	return securityTestResponse, err
}

// FindOneDBAnalysis checks if a given analysis is present into AnalysisCollection.
func FindOneDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error) {
	session := db.Connect()
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	analysisResponse := types.Analysis{}
	err := session.SearchOne(analysisFinalQuery, nil, db.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// FindAllDBRepository returns all Repository of a given query present into RepositoryCollection.
func FindAllDBRepository(mapParams map[string]interface{}) ([]types.Repository, error) {
	session := db.Connect()
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	repositoryResponse := []types.Repository{}
	err := session.Search(repositoryFinalQuery, nil, db.RepositoryCollection, &repositoryResponse)
	return repositoryResponse, err
}

// FindAllDBSecurityTest returns all SecurityTests of a given query present into SecurityTestCollection.
func FindAllDBSecurityTest(mapParams map[string]interface{}) ([]types.SecurityTest, error) {
	session := db.Connect()
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	securityTestResponse := []types.SecurityTest{}
	err := session.Search(securityTestFinalQuery, nil, db.SecurityTestCollection, &securityTestResponse)
	return securityTestResponse, err
}

// FindAllDBAnalysis returns all Analysis of a given query present into AnalysisCollection.
func FindAllDBAnalysis(mapParams map[string]interface{}) ([]types.Analysis, error) {
	session := db.Connect()
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	analysisResponse := []types.Analysis{}
	err := session.Search(analysisFinalQuery, nil, db.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// InsertDBRepository inserts a new repository with default securityTests into RepositoryCollection.
func InsertDBRepository(repository types.Repository) (types.Repository, error) {
	session := db.Connect()
	repository.CreatedAt = time.Now()

	// checking default and generic securityTests in SecurityTestCollection
	securityTestDefaultIDs := []bson.ObjectId{}
	securityTestDefaultQuery := map[string]interface{}{"default": true, "language": "generic"}
	securityTestDefaultResult, err := FindAllDBSecurityTest(securityTestDefaultQuery)
	if err != nil {
		fmt.Println("Err:", err)
	}
	for _, securityTest := range securityTestDefaultResult {
		securityTestDefaultIDs = append(securityTestDefaultIDs, securityTest.ID)
	}

	newRepository := bson.M{
		"URL":          repository.URL,
		"VM":           repository.VM,
		"createdAt":    repository.CreatedAt,
		"deletedAt":    repository.DeletedAt,
		"language":     repository.Language,
		"securityTest": securityTestDefaultIDs,
	}

	err = session.Insert(newRepository, db.RepositoryCollection)
	return repository, err
}

// InsertDBSecurityTest inserts a new securityTest into SecurityTestCollection.
func InsertDBSecurityTest(securityTest types.SecurityTest) (types.SecurityTest, error) {
	session := db.Connect()
	newSecurityTest := bson.M{
		"name":     securityTest.Name,
		"image":    securityTest.Image,
		"cmd":      securityTest.Cmd,
		"language": securityTest.Language,
		"default":  securityTest.Default,
	}
	err := session.Insert(newSecurityTest, db.SecurityTestCollection)
	return securityTest, err
}

// InsertDBAnalysis inserts a new analysis into AnalysisCollection.
func InsertDBAnalysis(analysis types.Analysis) (types.Analysis, error) {
	session := db.Connect()
	newAnalysis := bson.M{
		"RID":          analysis.RID,
		"URL":          analysis.URL,
		"securityTest": analysis.SecurityTestID,
		"status":       analysis.Status,
		"result":       analysis.Result,
		"container":    analysis.CID,
	}
	err := session.Insert(newAnalysis, db.AnalysisCollection)
	return analysis, err
}

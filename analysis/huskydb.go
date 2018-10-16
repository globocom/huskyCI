package analysis

import (
	"errors"
	"fmt"
	"time"

	db "github.com/globocom/husky/db/mongo"
	"github.com/globocom/husky/types"
	"gopkg.in/mgo.v2/bson"
)

// FindOneDBRepository checks if a given repository is present into RepositoryCollection.
func FindOneDBRepository(mapParams map[string]interface{}) (types.Repository, error) {
	repositoryResponse := types.Repository{}
	session, err := db.Connect()
	if err != nil {
		return repositoryResponse, err
	}
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	err = session.SearchOne(repositoryFinalQuery, nil, db.RepositoryCollection, &repositoryResponse)
	return repositoryResponse, err
}

// FindOneDBSecurityTest checks if a given securityTest is present into SecurityTestCollection.
func FindOneDBSecurityTest(mapParams map[string]interface{}) (types.SecurityTest, error) {
	securityTestResponse := types.SecurityTest{}
	session, err := db.Connect()
	if err != nil {
		return securityTestResponse, err
	}
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	err = session.SearchOne(securityTestFinalQuery, nil, db.SecurityTestCollection, &securityTestResponse)
	return securityTestResponse, err
}

// FindOneDBAnalysis checks if a given analysis is present into AnalysisCollection.
func FindOneDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error) {
	analysisResponse := types.Analysis{}
	session, err := db.Connect()
	if err != nil {
		return analysisResponse, err
	}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}

	err = session.SearchOne(analysisFinalQuery, nil, db.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// FindAllDBRepository returns all Repository of a given query present into RepositoryCollection.
func FindAllDBRepository(mapParams map[string]interface{}) ([]types.Repository, error) {
	session, err := db.Connect()
	if err != nil {
		return nil, err
	}
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	repositoryResponse := []types.Repository{}
	err = session.Search(repositoryFinalQuery, nil, db.RepositoryCollection, &repositoryResponse)
	return repositoryResponse, err
}

// FindAllDBSecurityTest returns all SecurityTests of a given query present into SecurityTestCollection.
func FindAllDBSecurityTest(mapParams map[string]interface{}) ([]types.SecurityTest, error) {
	session, err := db.Connect()
	if err != nil {
		return nil, err
	}
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	securityTestResponse := []types.SecurityTest{}
	err = session.Search(securityTestFinalQuery, nil, db.SecurityTestCollection, &securityTestResponse)
	return securityTestResponse, err
}

// FindAllDBAnalysis returns all Analysis of a given query present into AnalysisCollection.
func FindAllDBAnalysis(mapParams map[string]interface{}) ([]types.Analysis, error) {
	session, err := db.Connect()
	if err != nil {
		return nil, err
	}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	analysisResponse := []types.Analysis{}
	err = session.Search(analysisFinalQuery, nil, db.AnalysisCollection, &analysisResponse)
	return analysisResponse, err
}

// InsertDBRepository inserts a new repository with default securityTests into RepositoryCollection.
func InsertDBRepository(repository types.Repository) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	repository.CreatedAt = time.Now()
	securityTestList := []types.SecurityTest{}
	err = errors.New("")
	maxQueryAllowed := 8

	if len(repository.SecurityTestName) == 0 {
		// checking default and Generic securityTests in SecurityTestCollection
		securityTestQuery := map[string]interface{}{"default": true, "language": "Generic"}
		securityTestList, err = FindAllDBSecurityTest(securityTestQuery)
		if err != nil {
			fmt.Println("Could not find default securityTests:", err)
		}
	} else {
		// checking if a given securityTestName matches a securityTest
		repository.SecurityTestName = removeDuplicates(repository.SecurityTestName)
		if len(repository.SecurityTestName) > maxQueryAllowed {
			for _, securityTestName := range repository.SecurityTestName[:maxQueryAllowed] {
				securityTestQuery := map[string]interface{}{"name": securityTestName}
				securityTestResult, err := FindOneDBSecurityTest(securityTestQuery)
				if err != nil {
					fmt.Println("Could not find securityTestName:", securityTestName)
				} else {
					securityTestList = append(securityTestList, securityTestResult)
				}
			}
		} else {
			for _, securityTestName := range repository.SecurityTestName {
				securityTestQuery := map[string]interface{}{"name": securityTestName}
				securityTestResult, err := FindOneDBSecurityTest(securityTestQuery)
				if err != nil {
					fmt.Println("Could not find securityTestName:", securityTestName)
				} else {
					securityTestList = append(securityTestList, securityTestResult)
				}
			}
		}
	}

	newRepository := types.Repository{
		URL:           repository.URL,
		Branch:        repository.Branch,
		SecurityTests: securityTestList,
		VM:            repository.VM,
		CreatedAt:     repository.CreatedAt,
		DeletedAt:     repository.DeletedAt,
		Languages:     repository.Languages,
	}

	err = session.Insert(newRepository, db.RepositoryCollection)
	return err
}

// InsertDBSecurityTest inserts a new securityTest into SecurityTestCollection.
func InsertDBSecurityTest(securityTest types.SecurityTest) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	newSecurityTest := bson.M{
		"name":           securityTest.Name,
		"image":          securityTest.Image,
		"cmd":            securityTest.Cmd,
		"language":       securityTest.Language,
		"default":        securityTest.Default,
		"timeOutSeconds": securityTest.TimeOutInSeconds,
	}
	err = session.Insert(newSecurityTest, db.SecurityTestCollection)
	return err
}

// InsertDBAnalysis inserts a new analysis into AnalysisCollection.
func InsertDBAnalysis(analysis types.Analysis) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	newAnalysis := bson.M{
		"RID":          analysis.RID,
		"URL":          analysis.URL,
		"Branch":       analysis.Branch,
		"securityTest": analysis.SecurityTests,
		"status":       analysis.Status,
		"result":       analysis.Result,
		"containers":   analysis.Containers,
	}
	err = session.Insert(newAnalysis, db.AnalysisCollection)
	return err
}

// UpdateOneDBRepository checks if a given repository is present into RepositoryCollection and update it.
func UpdateOneDBRepository(mapParams, updateQuery map[string]interface{}) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	repositoryQuery := []bson.M{}
	for k, v := range mapParams {
		repositoryQuery = append(repositoryQuery, bson.M{k: v})
	}
	repositoryFinalQuery := bson.M{"$and": repositoryQuery}
	err = session.Update(repositoryFinalQuery, updateQuery, db.RepositoryCollection)
	return err
}

// UpdateOneDBSecurityTest checks if a given securityTest is present into SecurityTestCollection and update it.
func UpdateOneDBSecurityTest(mapParams map[string]interface{}, updatedSecurityTest types.SecurityTest) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	securityTestQuery := []bson.M{}
	for k, v := range mapParams {
		securityTestQuery = append(securityTestQuery, bson.M{k: v})
	}
	securityTestFinalQuery := bson.M{"$and": securityTestQuery}
	err = session.Update(securityTestFinalQuery, updatedSecurityTest, db.SecurityTestCollection)
	return err
}

// UpdateOneDBAnalysis checks if a given analysis is present into AnalysisCollection and update it.
func UpdateOneDBAnalysis(mapParams map[string]interface{}, updatedAnalysis types.Analysis) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	err = session.Update(analysisFinalQuery, updatedAnalysis, db.AnalysisCollection)
	return err
}

// UpdateOneDBAnalysisContainer checks if a given analysis is present into AnalysisCollection and update the container associated in it.
func UpdateOneDBAnalysisContainer(mapParams, updateQuery map[string]interface{}) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}
	analysisQuery := []bson.M{}
	for k, v := range mapParams {
		analysisQuery = append(analysisQuery, bson.M{k: v})
	}
	analysisFinalQuery := bson.M{"$and": analysisQuery}
	err = session.Update(analysisFinalQuery, updateQuery, db.AnalysisCollection)
	return err
}

// removeDuplicates remove duplicated itens from a slice.
func removeDuplicates(s []string) []string {
	mapS := make(map[string]string, len(s))
	i := 0
	for _, v := range s {
		if _, ok := mapS[v]; !ok {
			mapS[v] = v
			s[i] = v
			i++
		}
	}
	return s[:i]
}

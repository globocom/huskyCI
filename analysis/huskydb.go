package analysis

import (
	"fmt"
	"time"

	db "github.com/globocom/husky/db/mongo"
	"github.com/globocom/husky/types"
	"gopkg.in/mgo.v2/bson"
)

// FindDBRepository checks if a given repository is present into RepositoryCollection.
func FindDBRepository(mapParams map[string]interface{}) (types.Repository, error) {
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

// FindDBSecurityTest checks if a given securityTest is present into SecurityTestCollection.
func FindDBSecurityTest(mapParams map[string]interface{}) (types.SecurityTest, error) {
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

// FindDBAnalysis checks if a given analysis is present into AnalysisCollection.
func FindDBAnalysis(mapParams map[string]interface{}) (types.Analysis, error) {
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

// FindDBContainer checks if a given container is present into ContainerCollection.
func FindDBContainer(mapParams map[string]interface{}) (types.Container, error) {
	session := db.Connect()
	containerQuery := []bson.M{}
	for k, v := range mapParams {
		containerQuery = append(containerQuery, bson.M{k: v})
	}
	containerFinalQuery := bson.M{"$and": containerQuery}
	containerResponse := types.Container{}
	err := session.SearchOne(containerFinalQuery, nil, db.ContainerCollection, &containerResponse)
	return containerResponse, err
}

// InsertDBNewRepository inserts a new repository with default securityTests into RepositoryCollection.
func InsertDBNewRepository(repository types.Repository) (types.Repository, error) {
	session := db.Connect()

	securityTestEnryQuery := map[string]interface{}{"name": "enry"}
	securityTestEnryResponse, err := FindDBSecurityTest(securityTestEnryQuery)
	if err != nil {
		return repository, err
	}

	securityTestGasQuery := map[string]interface{}{"name": "gas"}
	securityTestGasResponse, err := FindDBSecurityTest(securityTestGasQuery)
	if err != nil {
		return repository, err
	}

	repository.CreatedAt = time.Now()
	repository.SecurityTestID = []bson.ObjectId{securityTestEnryResponse.ID, securityTestGasResponse.ID}

	newRepository := bson.M{
		"URL":          repository.URL,
		"VM":           repository.VM,
		"createdAt":    repository.CreatedAt,
		"deletedAt":    repository.DeletedAt,
		"securityTest": repository.SecurityTestID,
	}
	err = session.Insert(newRepository, db.RepositoryCollection)
	return repository, err
}

// InsertDBNewSecurityTest inserts a new securityTest into SecurityTestCollection.
func InsertDBNewSecurityTest(securityTest types.SecurityTest) (types.SecurityTest, error) {
	session := db.Connect()
	newSecurityTest := bson.M{
		"name":  securityTest.Name,
		"image": securityTest.Image,
		"cmd":   securityTest.Cmd,
	}
	err := session.Insert(newSecurityTest, db.SecurityTestCollection)
	return securityTest, err
}

// InsertDBNewAnalysis inserts a new analysis into AnalysisCollection.
func InsertDBNewAnalysis(analysis types.Analysis) (types.Analysis, error) {
	session := db.Connect()
	newAnalysis := bson.M{
		"RID":          analysis.RID,
		"URL":          analysis.URL,
		"securityTest": analysis.SecurityTestID,
		"status":       analysis.Status,
		"result":       analysis.Result,
		"output":       analysis.Output,
		"container":    analysis.CID,
	}
	err := session.Insert(newAnalysis, db.AnalysisCollection)
	return analysis, err
}

// InsertDBNewContainer inserts a new container into ContainerCollection's db.
func InsertDBNewContainer(container types.Container) (types.Container, error) {
	session := db.Connect()
	newContainer := bson.M{
		"CID":          container.CID,
		"RID":          container.RID,
		"VM":           container.VM,
		"securityTest": container.SecurityTestID,
		"cStatus":      container.CStatus,
		"cOutput":      container.COuput,
	}
	err := session.Insert(newContainer, db.ContainerCollection)
	return container, err
}

// InitDBSecurityTestCollection initiates SecurityTestCollection.
func InitDBSecurityTestCollection() error {
	session := db.Connect()
	securityTestEnry := types.SecurityTest{Name: "enry", Image: "huskyci/enry", Cmd: []string{"ls", "whoami"}}
	err := session.Insert(securityTestEnry, db.SecurityTestCollection)
	if err != nil {
		fmt.Println("InitSecurityTestCollection() - Error inserting securityTest ", securityTestEnry.Name, ":", err)
		return err
	}
	securityTestGAS := types.SecurityTest{Name: "gas", Image: "huskyci/gas", Cmd: []string{"command1", "ls"}}
	err = session.Insert(securityTestGAS, db.SecurityTestCollection)
	if err != nil {
		fmt.Println("InitSecurityTestCollection() - Error inserting securityTest ", securityTestEnry.Name, ":", err)
		return err
	}
	return err
}

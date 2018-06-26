package analysis

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/globocom/husky/db/mongo"
	"github.com/globocom/husky/types"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// HealthCheck is the heath check function.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!\n")
}

// StartAnalysis starts the analysis.
func StartAnalysis(c echo.Context) error {
	RID := c.Response().Header()["X-Request-Id"][0]
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		c.Logger()
		return c.JSON(http.StatusBadRequest, map[string]string{"RID": RID, "result": "error", "details": "Error binding repository."})
	}
	_, err = FindRepository(repository)
	if err != nil {
		err := InsertNewRepository(repository)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"RID": RID, "result": "error", "details": "Internal error."})
		}
	} else {
		for _, securityTestID := range repository.SecurityTestID {
			err = StartSecurityTest(RID, securityTestID, repository)
			if err != nil {
				fmt.Println("StartSecurityTest() - Could not start ", string(securityTestID), ":", err)
			}
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "received", "details": "Request received."})
}

// StartSecurityTest starts a given securityTestID in a given repository.
func StartSecurityTest(RID string, securityTestID bson.ObjectId, r types.Repository) error {
	// securityTestQuery := types.SecurityTest{ID: securityTestID}
	// securityTestReponse, err := FindSecurityTest(securityTestQuery)
	// docker := dockerapi.Docker{}
	// err = docker.RunContainer(securityTestReponse.Image, securityTestReponse.Cmd)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// FindRepository checks if a given repository is present in db.
func FindRepository(repository types.Repository) (types.Repository, error) {
	session := db.Connect()
	repositoryQuery := bson.M{
		"ID":             repository.ID,
		"URL":            repository.URL,
		"VM":             repository.VM,
		"SecurityTestID": repository.SecurityTestID,
		"CreatedAt":      repository.CreatedAt,
		"DeletedAt":      repository.DeletedAt,
	}
	repositoryResponse := types.Repository{}
	err := session.SearchOne(repositoryQuery, nil, db.RepositoryCollection, &repositoryResponse)
	if err != nil {
		fmt.Println("FindRepository() - Error finding ", repositoryResponse.URL, ":", err)
		return repositoryResponse, err
	}
	return repositoryResponse, err
}

// FindSecurityTest checks if a given securityTest is present in DB.
func FindSecurityTest(securityTest types.SecurityTest) (types.SecurityTest, error) {
	session := db.Connect()
	securityTestQuery := bson.M{
		"ID":    securityTest.ID,
		"Name":  securityTest.Name,
		"Image": securityTest.Image,
		"Cmd":   securityTest.Cmd,
	}
	securityTestResponse := types.SecurityTest{}
	err := session.SearchOne(securityTestQuery, nil, db.SecurityTestCollection, &securityTestResponse)
	if err != nil {
		fmt.Println("FindSecurityTest() - Error finding ", securityTest.Name, ":", err)
		return securityTestResponse, err
	}
	return securityTestResponse, err
}

// InsertNewRepository inserts a new repository with default securityTests into DB.
func InsertNewRepository(repository types.Repository) error {
	session := db.Connect()
	securityTestQueryEnry := types.SecurityTest{Name: "enry"}
	securityTestResponse, err := FindSecurityTest(securityTestQueryEnry)
	if err != nil {
		return err
	}
	repository.CreatedAt = time.Now().Format(time.RFC850)
	newRepository := bson.M{
		"URL":          repository.URL,
		"VM":           repository.VM,
		"createdAt":    repository.CreatedAt,
		"deletedAt":    repository.DeletedAt,
		"securityTest": []bson.ObjectId{securityTestResponse.ID},
	}
	err = session.Insert(newRepository, db.RepositoryCollection)
	if err != nil {
		fmt.Println("InsertNewRepository() - Error inserting repository:", err)
		return err
	}
	return err
}

// InitSecurityTestCollection initiates SecurityTestCollection.
func InitSecurityTestCollection() error {
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

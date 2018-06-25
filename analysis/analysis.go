package analysis

import (
	"fmt"
	"net/http"

	db "github.com/globocom/husky/db/mongo"
	"github.com/globocom/husky/types"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// HealthCheck is the heath check function
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!\n")
}

// StartAnalysis starts the analysis
func StartAnalysis(c echo.Context) error {

	repository := types.Repository{}
	requestID := c.Response().Header()["X-Request-Id"][0]

	if err := c.Bind(&repository); err != nil {
		c.Logger()
		return c.JSON(http.StatusBadRequest, map[string]string{"RID": requestID, "result": "error", "details": "Error binding repository."})
	}

	repositoryChecked, err := CheckRepository(repository)
	if err != nil {
		_, err := InsertRepository(repositoryChecked)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"RID": requestID, "result": "error", "details": "Internal error."})
		}
	} else {
		fmt.Println(repositoryChecked.URL, "found! Have to execute:", repositoryChecked.SecurityTest)
	}

	return c.JSON(http.StatusOK, map[string]string{"RID": requestID, "result": "received", "details": "Request received."})
}

// CheckRepository checks if repository.URL is present in db.
func CheckRepository(r types.Repository) (types.Repository, error) {

	session := db.Connect()
	query := bson.M{"URL": r.URL}

	err := session.SearchOne(query, nil, db.RepositoryCollection, &r)
	if err != nil {
		fmt.Println("Error SearchOne() - CheckRepository:", r.URL, err)
		return r, err
	}
	return r, err
}

// CheckSecurityTest checks if securityTest.Name is present in db.
func CheckSecurityTest(t types.SecurityTest) (types.SecurityTest, error) {

	session := db.Connect()
	query := bson.M{"name": t.Name}

	err := session.SearchOne(query, nil, db.SecurityTestCollection, &t)
	if err != nil {
		fmt.Println("Error SearchOne() - CheckSecurityTest:", t.Name, err)
		return t, err
	}
	return t, err
}

// InsertRepository inserts repositoryURL received from POST into DB.
func InsertRepository(r types.Repository) (types.Repository, error) {

	session := db.Connect()
	initialTests := []string{"123", "4321"}
	r.SecurityTest = initialTests
	query := bson.M{
		"URL":          r.URL,
		"VM":           r.VM,
		"createdAt":    r.CreatedAt,
		"deletedAt":    r.DeletedAt,
		"securityTest": r.SecurityTest,
	}

	_, err := session.Upsert(query, &r, db.RepositoryCollection)
	if err != nil {
		fmt.Println("Error Upsert() - InsertRepository:", err)
		return r, err
	}
	return r, err
}

// InitSecurityTestCollection initiates SecurityTestCollection
func InitSecurityTestCollection() error {

	session := db.Connect()

	securityTestEnry := types.SecurityTest{Name: "enry", Image: "huskyci/enry", Cmd: []string{"ls", "whoami"}}
	queryEnry := bson.M{"name": securityTestEnry.Name, "image": securityTestEnry.Image, "cmd": securityTestEnry.Cmd}
	_, err := session.Upsert(queryEnry, &securityTestEnry, db.SecurityTestCollection)
	if err != nil {
		fmt.Println("Error Upsert() - queryEnry - InitSecurityTestCollection:", err)
		return err
	}

	securityTestGAS := types.SecurityTest{Name: "gas", Image: "huskyci/gas", Cmd: []string{"command1", "ls"}}
	queryGAS := bson.M{"name": securityTestGAS.Name, "image": securityTestGAS.Image, "cmd": securityTestGAS.Cmd}
	_, err = session.Upsert(queryGAS, &securityTestGAS, db.SecurityTestCollection)
	if err != nil {
		fmt.Println("Error Upsert() - queryGAS - InitSecurityTestCollection:", err)
		return err
	}

	return err
}

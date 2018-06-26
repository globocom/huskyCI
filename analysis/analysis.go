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

	_, err = FindRepository(repository.URL)
	if err != nil {
		err := InsertNewRepository(repository.URL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"RID": RID, "result": "error", "details": "Internal error."})
		}
	} else {
		fmt.Println(repository.URL, "found! Have to execute:", repository.SecurityTest)
	}

	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "received", "details": "Request received."})
}

// FindRepository checks if repository.URL is present in db.
func FindRepository(repositoryURL string) (types.Repository, error) {

	session := db.Connect()
	query := bson.M{"URL": repositoryURL}
	r := types.Repository{URL: repositoryURL}

	err := session.SearchOne(query, nil, db.RepositoryCollection, &r)
	if err != nil {
		fmt.Println("Error SearchOne() - CheckRepository:", r.URL, err)
		return r, err
	}
	return r, err
}

// FindSecurityTest checks if securityTest.Name is present in db.
func FindSecurityTest(securityTestName string) (types.SecurityTest, error) {

	session := db.Connect()
	query := bson.M{"name": securityTestName}
	t := types.SecurityTest{Name: securityTestName}

	err := session.SearchOne(query, nil, db.SecurityTestCollection, &t)
	if err != nil {
		fmt.Println("Error SearchOne() - CheckSecurityTest:", t.Name, err)
		return t, err
	}
	return t, err
}

// InsertNewRepository inserts a new repository with the default tests into DB.
func InsertNewRepository(repositoryURL string) error {

	session := db.Connect()

	// Getting objectID of enry.
	securityTest, err := FindSecurityTest("enry")
	if err != nil {
		fmt.Println("Could not find enry in MongoDB:", err)
		return err
	}

	r := types.Repository{URL: repositoryURL, CreatedAt: time.Now().Format(time.RFC850)}

	newRepository := bson.M{
		"URL":          r.URL,
		"VM":           r.VM,
		"createdAt":    r.CreatedAt,
		"deletedAt":    r.DeletedAt,
		"securityTest": []bson.ObjectId{securityTest.ID},
	}

	err = session.Insert(newRepository, db.RepositoryCollection)
	if err != nil {
		fmt.Println("Error Insert() - InsertRepository:", err)
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
		fmt.Println("Error Upsert() - queryEnry - InitSecurityTestCollection:", err)
		return err
	}

	securityTestGAS := types.SecurityTest{Name: "gas", Image: "huskyci/gas", Cmd: []string{"command1", "ls"}}
	err = session.Insert(securityTestGAS, db.SecurityTestCollection)
	if err != nil {
		fmt.Println("Error Upsert() - queryGAS - InitSecurityTestCollection:", err)
		return err
	}

	return err
}

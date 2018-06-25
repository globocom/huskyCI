package analysis

import (
	"fmt"
	"net/http"
	"os"

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

// CheckRepository will check if repositoryURL is present in BD.
func CheckRepository(r types.Repository) (types.Repository, error) {

	session := db.Connect()
	collection := os.Getenv("MONGO_COLLECTION_REPOSITORY")
	query := bson.M{"URL": r.URL}

	err := session.SearchOne(query, nil, collection, &r)
	if err != nil {
		fmt.Println("Error SearchOne():", r.URL, err)
		return r, err
	}

	return r, err
}

// InsertRepository will insert repositoryURL received from POST into DB.
func InsertRepository(r types.Repository) (types.Repository, error) {

	session := db.Connect()
	collection := os.Getenv("MONGO_COLLECTION_REPOSITORY")
	initialTests := []string{"123", "4321"}
	r.SecurityTest = initialTests
	query := bson.M{
		"URL":          r.URL,
		"VM":           r.VM,
		"createdAt":    r.CreatedAt,
		"deletedAt":    r.DeletedAt,
		"securityTest": r.SecurityTest,
	}

	_, err := session.Upsert(query, &r, collection)
	if err != nil {
		fmt.Println("Error Upsert():", err)
		return r, err
	}
	return r, err
}

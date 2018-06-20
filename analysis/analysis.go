package analysis

import (
	"fmt"
	"net/http"
	"os"

	docker "github.com/globocom/husky/dockers"

	db "github.com/globocom/husky/db/mongo"
	"github.com/globocom/husky/types"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// ResultDB represents all data retreived from mongo
type ResultDB struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	URL string        `bson:"URL"`
}

// HealthCheck is the heath check function
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!\n")
}

// StartAnalysis starts the analysis
func StartAnalysis(c echo.Context) error {

	// Parsing input
	repo := new(types.Repository)
	if err := c.Bind(repo); err != nil {
		fmt.Println("Error binding repositoryURL:", err)
	}

	// err := CheckMongoRepoURL(repo.URL)
	// if err != nil {
	// 	return c.String(http.StatusOK, repo.URL+" not found.\n")
	// }

	// err = d.PullImage("ubuntu")
	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// }

	// cmd := "whoami"
	// err = d.RunContainer(c, "ubuntu", cmd)
	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// }

	d := new(docker.Docker)
	images := d.ListImages()
	fmt.Println(images)

	return c.String(http.StatusOK, repo.URL+" found!\n")
}

// CheckMongoRepoURL will query mongo to check if repositoryURL is present.
func CheckMongoRepoURL(repositoryURL string) error {

	session := db.Connect()
	query := bson.M{"URL": repositoryURL}
	result := ResultDB{}
	collection := os.Getenv("MONGO_COLLECTION")

	err := session.SearchOne(query, nil, collection, &result)
	if err != nil {
		fmt.Println("Error SearchOne():", repositoryURL, err)
		result.URL = repositoryURL
		_, err = session.Upsert(query, &result, collection)
		if err != nil {
			fmt.Println("Error Upser():", err)
			return err
		}
		return err
	}
	return err
}

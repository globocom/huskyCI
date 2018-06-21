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

// ResultDB represents all data retreived from mongo
type ResultDB struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	URL   string        `bson:"URL"`
	VM    string        `bson:"VM"`
	Tests []string      `bson:"tests"`
}

// HealthCheck is the heath check function
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!\n")
}

// StartAnalysis starts the analysis
func StartAnalysis(c echo.Context) error {

	repo := new(types.Repository)
	if err := c.Bind(repo); err != nil {
		fmt.Println("Error binding repositoryURL:", err)
	}

	tests, err := GetRepoTests(repo.URL)

	if err != nil {
		fmt.Println("Error checking Mongo for repositoryURL.", err)
	}
	if tests != nil {
		for i, teste := range tests {
			fmt.Println("Teste:", i, teste)
		}
	} else {
		fmt.Println("Nenhum teste encontrado! Associa o primeiro testo ao respositório: ENRY:", tests)

	}

	// err = d.PullImage("ubuntu")
	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// }

	// cmd := "whoami"
	// err = d.RunContainer(c, "ubuntu", cmd)
	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// }

	// d := new(docker.Docker)
	// images := d.ListImages()
	// fmt.Println(images)

	return c.String(http.StatusOK, "Request received!\n")
}

// GetRepoTests will query mongo to check if repositoryURL is present. If it is not, it will include it.
func GetRepoTests(repositoryURL string) ([]string, error) {

	session := db.Connect()
	query := bson.M{"URL": repositoryURL}
	result := ResultDB{}
	collection := os.Getenv("MONGO_COLLECTION_REPOSITORY")

	err := session.SearchOne(query, nil, collection, &result)
	if err != nil {
		fmt.Println("Error SearchOne():", repositoryURL, err)
		result.URL = repositoryURL
		_, err = session.Upsert(query, &result, collection)
		if err != nil {
			fmt.Println("Error Upser():", err)
		}
	}
	return result.Tests, err
}

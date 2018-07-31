package analysis

import (
	"fmt"
	"net/http"
	"regexp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/globocom/husky/types"
	"github.com/labstack/echo"
)

// HealthCheck is the heath check function.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!\n")
}

// ReceiveRequest receives the request and performs several checks before starting a new analysis.
func ReceiveRequest(c echo.Context) error {
	RID := c.Response().Header()["X-Request-Id"][0]

	// check-00: is this a valid JSON?
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Error binding repository."})
	}

	// check-01: is this a git repository URL?
	regexpGit := `^(?:git|https?|ssh|git@[-\w.]+):(//)?(.*?)(\.git)(/?|#[-\d\w._]+?)$`
	valid, err := regexp.MatchString(regexpGit, repository.URL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Internal error."})
	}
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "URL received is not a git repository."})
	}

	// check-02: is this repository in MongoDB?
	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	repositoryResult, err := FindOneDBRepository(repositoryQuery)
	if err == nil {
		// check-03: repository found! does it have a running status analysis? (for the future: check commits and not URLs?)
		analysisQuery := map[string]interface{}{"URL": repository.URL}
		analysisResult, err := FindOneDBAnalysis(analysisQuery)
		if err != mgo.ErrNotFound {
			if analysisResult.Status == "running" {
				return c.JSON(http.StatusConflict, map[string]string{"result": "error", "details": "An analysis is already in place for this URL."})
			}
		}
	} else {
		// repository not found! insert it into MongoDB with default securityTests
		err = InsertDBRepository(repository)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"result": "error", "details": "Internal error inserting repository."})
		}
		repositoryQuery := map[string]interface{}{"URL": repository.URL}
		repositoryResult, err = FindOneDBRepository(repositoryQuery)
		if err != nil {
			// well it was supposed to be there, after all, we just inserted it.
			return c.JSON(http.StatusInternalServerError, map[string]string{"result": "error", "details": "Internal error finding repository."})
		}
	}

	go StartAnalysis(RID, repositoryResult)

	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "ok", "details": "Request received."})
}

// StartAnalysis starts the analysis given a RID and a repository.
func StartAnalysis(RID string, repository types.Repository) {

	// step 0: create a new analysis struct
	newAnalysis := types.Analysis{
		RID:           RID,
		URL:           repository.URL,
		SecurityTests: repository.SecurityTests,
		Status:        "started",
		Containers:    make([]types.Container, 0),
	}

	// step 1: insert new analysis into MongoDB
	err := InsertDBAnalysis(newAnalysis)
	if err != nil {
		fmt.Println("Error inserting new analysis.", err)
	}

	// step 2: increment -1 to repository.LimitEnryScan
	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	updateRepositoryQuery := bson.M{
		"$inc": bson.M{
			"limitEnryScan": -1,
		},
	}
	err = UpdateOneDBRepository(repositoryQuery, updateRepositoryQuery)
	if err != nil {
		fmt.Println("Could not increment repository.LimitEnryScan:", err)
		return
	}

	// step 3: check if a new enry scan is needed
	if repository.LimitEnryScan == 0 {
		// new enry scan is needed
		limitEnryScan := 10
		repository.SecurityTests = nil
		enrySecurityTestQuery := map[string]interface{}{"name": "enry"}
		enrySecurityTestResult, err := FindOneDBSecurityTest(enrySecurityTestQuery)
		if err != nil {
			fmt.Println("Error finding enry securityTest:", err)
		}
		repository.SecurityTests = append(repository.SecurityTests, enrySecurityTestResult)
		// set repository.LimitEnryScan to its default value
		repositoryQuery := map[string]interface{}{"URL": repository.URL}
		updateRepositoryQuery := bson.M{
			"$set": bson.M{
				"limitEnryScan": limitEnryScan,
			},
		}
		err = UpdateOneDBRepository(repositoryQuery, updateRepositoryQuery)
		if err != nil {
			fmt.Println("Could not set repository.LimitEnryScan to its default value:", err)
			return
		}
	}

	// step 4: start each securityTest set
	for _, securityTest := range repository.SecurityTests {
		go DockerRun(RID, &newAnalysis, securityTest)
	}

	// step 5: worker will check if the jobs are done to set newAnalysis.Status = "finished"
}

// StatusAnalysis returns the status of a given analysis (via RID).
func StatusAnalysis(c echo.Context) error {
	RID := c.Param("id")
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err == mgo.ErrNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"result": "error", "details": "Analysis not found."})
	} // What if DB is not reachable!? else { }
	return c.JSON(http.StatusFound, analysisResult)
}

// CreateNewSecurityTest inserts the given securityTest into SecurityTestCollection.
func CreateNewSecurityTest(c echo.Context) error {
	securityTest := types.SecurityTest{}
	err := c.Bind(&securityTest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Error binding securityTest."})
	}

	securityTestQuery := map[string]interface{}{"name": securityTest.Name}
	_, err = FindOneDBSecurityTest(securityTestQuery)
	if err != mgo.ErrNotFound {
		return c.JSON(http.StatusConflict, map[string]string{"result": "error", "details": "This securityTest is already in MongoDB."})
	}

	err = InsertDBSecurityTest(securityTest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": "error", "details": "Error creating new securityTest."})
	}

	return c.JSON(http.StatusCreated, map[string]string{"result": "created", "details": "securityTest sucessfully created."})
}

// CreateNewRepository inserts the given repository into RepositoryCollection.
func CreateNewRepository(c echo.Context) error {
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Error binding repository."})
	}

	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	_, err = FindOneDBRepository(repositoryQuery)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"result": "error", "details": "Repository found."})
	}

	err = InsertDBRepository(repository)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": "error", "details": "Error creating new repository."})
	}

	return c.JSON(http.StatusCreated, map[string]string{"result": "created", "details": "repository sucessfully created."})
}

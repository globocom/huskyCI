package analysis

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	docker "github.com/globocom/husky/dockers"
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

	// check-01: is this a valid JSON?
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Error binding repository."})
	}

	// check-02: is this repository in MongoDB?
	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	repositoryResult, err := FindOneDBRepository(repositoryQuery)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"result": "error", "details": "Repository not found."})
	} else {
		// ok let's insert into MongoDB the repository with default securityTests
		err = InsertDBRepository(repositoryResult)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"result": "error", "details": "Internal error."})
		}
	}

	// check-03: does this repository have a running status analysis? (for the future: check commits and not URLs?)
	analysisQuery := map[string]interface{}{"URL": repository.URL}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err != mgo.ErrNotFound {
		if analysisResult.Status == "running" {
			return c.JSON(http.StatusConflict, map[string]string{"result": "error", "details": "An analysis is already in place for this URL."})
		}
	}

	go StartAnalysis(RID, repositoryResult)

	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "ok", "details": "Request received."})
}

// StartAnalysis starts the analysis given a RID and a repository.
func StartAnalysis(RID string, repository types.Repository) {

	newAnalysis := types.Analysis{
		RID:           RID,
		URL:           repository.URL,
		SecurityTests: repository.SecurityTests,
		Status:        "started",
		Containers:    make([]types.Container, 0),
	}

	err := InsertDBAnalysis(newAnalysis)
	if err != nil {
		fmt.Println("Error inserting new analysis.", err)
	}

	for _, securityTest := range repository.SecurityTests {
		go dockerRun(RID, &newAnalysis, securityTest)
	}

	// worker will check if the jobs are done to set newAnalysis.Status = "finished"
}

func dockerRun(RID string, newAnalysis *types.Analysis, securityTest types.SecurityTest) {

	d := docker.Docker{}

	newContainer := types.Container{
		SecurityTest: securityTest,
		StartedAt:    time.Now(),
		CStatus:      "started",
	}

	CID, err := d.CreateContainer(*newAnalysis, securityTest.Image, securityTest.Cmd)
	if err != nil {
		fmt.Println("Error creating new container:", err)
	}

	newContainer.CID = CID
	newContainer.CStatus = "running"

	err = d.StartContainer(CID)
	if err != nil {
		fmt.Println("Error starting the container", CID, ":", err)
		newContainer.CStatus = "error"
	}

	err = d.WaitContainer(CID)
	if err != nil {
		fmt.Println("Error waiting the container", CID, ":", err)
	}

	newContainer.COuput = d.ReadOutput(CID)
	newContainer.CStatus = "finished"
	newContainer.FinishedAt = time.Now()

	newAnalysis.Containers = append(newAnalysis.Containers, newContainer)

	analysisQuery := map[string]interface{}{"RID": RID}
	err = UpdateOneDBAnalysis(analysisQuery, *newAnalysis)
	if err != nil {
		fmt.Println("Error updating Analysis", err)
	}

}

// StatusAnalysis returns the status of a given analysis (via RID).
func StatusAnalysis(c echo.Context) error {
	RID := c.Param("id")
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err == mgo.ErrNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"result": "error", "details": "Analysis not found."})
	} else {
		// What if DB is not reachable!?
	}
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

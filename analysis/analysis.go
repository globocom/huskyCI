package analysis

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

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
		return c.JSON(http.StatusBadRequest, map[string]string{"RID": RID, "result": "error", "details": "Error binding repository."})
	}

	// essa URL é um repositorio git? verificar bibliotecas default.

	// does this URL have already a running status analysis?
	// same commit ok! check this
	analysisQuery := map[string]interface{}{"URL": repository.URL}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err != mgo.ErrNotFound {
		// found an analysis for this URL. Is it running?
		if analysisResult.Status == "running" {
			return c.JSON(http.StatusConflict, map[string]string{"RID": RID, "result": "error", "details": "An analysis is already in place for this URL."})
		}
	}

	// does this URL have already a document in MongoDB?
	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	repositoryResult, err := FindOneDBRepository(repositoryQuery)
	if err == mgo.ErrNotFound {
		// inserting a new document for this URL.
		repositoryResult, err = InsertDBRepository(repository)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"RID": RID, "result": "error", "details": "Internal error. InsertDBNewRepository()."})
		}
	} else {
		// what if MongoDB is not reachable!?
		// case! check this
		//return c.JSON(http.StatusInternalServerError, map[string]string{"RID": RID, "result": "error", "details": "Internal error. Check MongoDB. FindBDRepository()."})
	}

	// analysis info to be included later into MongoDB.
	newAnalysis := types.Analysis{
		RID:            RID,
		URL:            repositoryResult.URL,
		SecurityTestID: repositoryResult.SecurityTestID,
		Status:         "running",
		Result:         "",
	}

	// starting each securityTest for URL
	for _, securityTestID := range repositoryResult.SecurityTestID {
		CID, err := StartSecurityTest(RID, securityTestID, repository)
		if err != nil {
			// change this! if one fails, try again. maxTries
			// glbgelf include!!!! https://gitlab.globoi.com/supseg/glbgelf
			c.Logger().Info("e")
			return c.JSON(http.StatusInternalServerError, map[string]string{"RID": RID, "result": "error", "details": "Internal error. StartSecurityTest()."})
		}
		newAnalysis.CID = append(newAnalysis.CID, CID)
	}

	_, err = InsertDBAnalysis(newAnalysis)
	if err != nil {
		fmt.Print("Error", err)
	}

	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "ok", "details": "Request received."})
}

// StartSecurityTest starts a given securityTestID in a given repository and returns the containerID.
func StartSecurityTest(RID string, securityTestID bson.ObjectId, r types.Repository) (string, error) {
	// securityTestQuery := map[string]interface{}{"_id": securityTestID}
	// securityTestReponse, err := FindSecurityTest(securityTestQuery)
	// if err != nil {
	// 	return err
	// }

	// docker := dockerapi.Docker{}
	// containerID, err := docker.CreateContainer(securityTestReponse.Name, securityTestReponse.Cmd)
	// if err != nil {
	// 	fmt.Println("Erro!", err)
	// }
	// fmt.Println("Sucesso!", containerID)

	newContainer := types.Container{
		CID:            "123asdvxcv12",
		RID:            RID,
		VM:             "10.10.10.10",
		SecurityTestID: securityTestID,
		CStatus:        "running",
		COuput:         []string{},
	}

	_, err := InsertDBContainer(newContainer)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return newContainer.CID, nil
}

// StatusAnalysis returns the status of a given analysis (via RID).
func StatusAnalysis(c echo.Context) error {

	// checking if given RID is present into AnalysisCollection.
	RID := c.Param("id")
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err == mgo.ErrNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"RID": RID, "result": "error", "details": "Analysis not found."})
	} else {
		// What if DB is not reachable!?
	}

	// checking all securityTests associated with this analysis into SecurityTestCollection.
	securityTestList := []types.SecurityTest{}
	for _, securityID := range analysisResult.SecurityTestID {
		securityTestQuery := map[string]interface{}{"_id": securityID}
		securityTestResult, err := FindOneDBSecurityTest(securityTestQuery)
		if err != nil {
			fmt.Println("Error", err)
			return err
		}
		securityTestList = append(securityTestList, securityTestResult)
	}

	// checking all containers associated with this analysis into ContainerCollection
	containerList := []types.Container{}
	for _, CID := range analysisResult.CID {
		containerQuery := map[string]interface{}{"CID": CID}
		containerResult, err := FindOneDBContainer(containerQuery)
		if err != nil {
			fmt.Println("Error", err)
			return err
		}
		containerList = append(containerList, containerResult)
	}

	statusFound := map[string]interface{}{
		"RID":          RID,
		"URL":          analysisResult.URL,
		"result":       analysisResult.Result,
		"status":       analysisResult.Status,
		"securityTest": securityTestList,
		"container":    containerList,
	}

	return c.JSON(http.StatusFound, statusFound)
}

// CreateNewSecurityTest inserts the given securityTest into SecurityTestCollection.
func CreateNewSecurityTest(c echo.Context) error {
	RID := c.Response().Header()["X-Request-Id"][0]
	securityTest := types.SecurityTest{}
	err := c.Bind(&securityTest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"RID": RID, "result": "error", "details": "Error binding securityTest."})
	}

	securityTestQuery := map[string]interface{}{"name": securityTest.Name}
	_, err = FindOneDBSecurityTest(securityTestQuery)
	if err != mgo.ErrNotFound {
		return c.JSON(http.StatusConflict, map[string]string{"RID": RID, "result": "error", "details": "This securityTest is already in MongoDB."})
	}

	_, err = InsertDBSecurityTest(securityTest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"RID": RID, "result": "error", "details": "Error creating new securityTest."})
	}

	return c.JSON(http.StatusCreated, map[string]string{"RID": RID, "result": "created", "details": "securityTest sucessfully created."})
}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"net/http"
	"regexp"
	"time"

	"github.com/globocom/huskyci/api/log"
	"github.com/globocom/huskyci/api/types"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Version holds the API version to be returned in /version route.
var Version types.VersionAPI

// HealthCheck is the heath check function.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!\n")
}

//VersionHandler returns the API version
func VersionHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Version)
}

// ReceiveRequest receives the request and performs several checks before starting a new analysis.
func ReceiveRequest(c echo.Context) error {
	RID := c.Response().Header().Get(echo.HeaderXRequestID)

	// check-00: is this a valid JSON?
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1015, err)
		return c.String(http.StatusBadRequest, "This is an invalid JSON.\n")
	}

	// check-01: is this a git repository URL and a branch?
	regexpGit := `((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`
	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repository.URL)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository URL regexp ", err)
		return c.String(http.StatusInternalServerError, "Internal error 1008.\n")
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1016, repository.URL)
		return c.String(http.StatusBadRequest, "This is not a valid repository URL.\n")
	}
	matches := r.FindString(repository.URL)
	repository.URL = matches

	regexpBranch := `^[a-zA-Z0-9_\.-]*$`
	valid, err = regexp.MatchString(regexpBranch, repository.Branch)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository Branch regexp ", err)
		return c.String(http.StatusInternalServerError, "Internal error 1008.\n")
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1017, repository.Branch)
		return c.String(http.StatusBadRequest, "This is not a valid branch.\n")
	}

	// check-02: is this repository in MongoDB?
	repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
	repositoryResult, err := FindOneDBRepository(repositoryQuery)
	if err == nil {
		// check-03: repository found! does it have a running status analysis?
		analysisQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		analysisResult, err := FindOneDBAnalysis(analysisQuery)
		if err != nil {
			if err != mgo.ErrNotFound {
				if analysisResult.Status == "running" {
					log.Warning("ReceiveRequest", "ANALYSIS", 104, analysisResult.URL)
					return c.String(http.StatusConflict, "An analysis is already in place for this URL.\n")
				}
			}
			log.Error("ReceiveRequest", "ANALYSIS", 1009, err)
		}
	} else {
		// repository not found! insert it into MongoDB with default securityTests
		err = InsertDBRepository(repository)
		if err != nil {
			log.Error("ReceiveRequest", "ANALYSIS", 1010, err)
			return c.String(http.StatusInternalServerError, "Internal error 1010.\n")
		}
		repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		repositoryResult, err = FindOneDBRepository(repositoryQuery)
		if err != nil {
			// well it was supposed to be there, after all, we just inserted it.
			log.Error("ReceiveRequest", "ANALYSIS", 1011, err)
			return c.String(http.StatusInternalServerError, "Internal error 1011.\n")
		}
	}

	log.Info("ReceiveRequest", "ANALYSIS", 16, repository.Branch, repository.URL)
	go StartAnalysis(RID, repositoryResult)
	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "ok", "details": "Request received."})
}

// StartAnalysis starts the analysis given a RID and a repository.
func StartAnalysis(RID string, repository types.Repository) {

	// step 0: create a new analysis struct
	newAnalysis := types.Analysis{
		RID:        RID,
		URL:        repository.URL,
		Branch:     repository.Branch,
		Status:     "running",
		Containers: make([]types.Container, 0),
	}

	// step 1: insert new analysis into MongoDB.
	err := InsertDBAnalysis(newAnalysis)
	if err != nil {
		log.Error("StartAnalysis", "ANALYSIS", 2011, err)
		return
	}

	// step 2: start enry and EnryStartAnalysis will start all others securityTests
	enryQuery := map[string]interface{}{"name": "enry"}
	enrySecurityTest, err := FindOneDBSecurityTest(enryQuery)
	if err != nil {
		log.Error("StartAnalysis", "ANALYSIS", 2011, "enry", err)
		return
	}
	DockerRun(RID, &newAnalysis, enrySecurityTest)

	// step 3: worker will check if jobs are done to set newAnalysis.Status = "finished".
	go MonitorAnalysis(&newAnalysis)

}

// MonitorAnalysis querys an analysis every retryTick seconds to check if it has already finished.
func MonitorAnalysis(analysis *types.Analysis) {

	timeout := time.After(10 * time.Minute)
	retryTick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeout:
			// cenario 1: MonitorAnalysis has timed out!
			log.Warning("MonitorAnalysis", "ANALYSIS", 105, analysis.RID)
			if err := registerAnalysisTimedOut(analysis.RID); err != nil {
				log.Error("MonitorAnalysis", "ANALYSIS", 2013, err)
				return
			}
			return
		case <-retryTick:
			// check if analysis has already finished.
			analysisHasFinished, err := monitorAnalysisCheckStatus(analysis.RID)
			if err != nil {
				// already being logged
			}
			// cenario 2: analysis has finished!
			if analysisHasFinished {
				err := monitorAnalysisUpdateStatus(analysis.RID)
				if err != nil {
					// already being logged
				}
			} // cenario 3: retry after retryTick seconds!
		}
	}

}

// registerAnalysisTimedOut updates the status of a given analysis to "timedout".
func registerAnalysisTimedOut(RID string) error {
	analysisQuery := map[string]interface{}{"RID": RID}
	updateAnalysisQuery := bson.M{
		"$set": bson.M{
			"status": "timedout",
		},
	}
	err := UpdateOneDBAnalysisContainer(analysisQuery, updateAnalysisQuery)
	return err
}

// monitorAnalysisUpdateStatus updates status and result of a given analysis.
func monitorAnalysisUpdateStatus(RID string) error {
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("monitorAnalysisUpdateStatus", "ANALYSIS", 2014, RID, err)
		return err
	}
	// analyze each cResult from each container to determine what is the value of analysis.Result
	finalResult := "passed"
	for _, container := range analysisResult.Containers {
		if container.CResult == "failed" {
			finalResult = "failed"
			break
		}
	}
	updateAnalysisQuery := bson.M{
		"$set": bson.M{
			"status": "finished",
			"result": finalResult,
		},
	}
	err = UpdateOneDBAnalysisContainer(analysisQuery, updateAnalysisQuery)
	if err != nil {
		log.Error("monitorAnalysisUpdateStatus", "ANALYSIS", 2007, err)
	}
	return err
}

// monitorAnalysisCheckStatus checks if an analysis has already finished and returns the correspoding boolean.
func monitorAnalysisCheckStatus(RID string) (bool, error) {
	analysisFinished := false
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("monitorAnalysisCheckStatus", "ANALYSIS", 2014, RID, err)
	}
	for _, container := range analysisResult.Containers {
		if container.CStatus != "finished" {
			analysisFinished = false
			break
		} else {
			analysisFinished = true
		}
	}
	return analysisFinished, err
}

// StatusAnalysis returns the status of a given analysis (via RID).
func StatusAnalysis(c echo.Context) error {

	RID := c.Param("id")
	regexpRID := `^[a-zA-Z0-9]*$`
	valid, err := regexp.MatchString(regexpRID, RID)
	if err != nil {
		log.Error("StatusAnalysis", "ANALYSIS", 1008, "RID regexp ", err)
		return c.String(http.StatusInternalServerError, "Internal error 1008.\n")
	}
	if !valid {
		log.Warning("StatusAnalysis", "ANALYSIS", 107, RID)
		return c.String(http.StatusBadRequest, "This is not a valid RID.\n")
	}

	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := FindOneDBAnalysis(analysisQuery)
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Warning("StatusAnalysis", "ANALYSIS", 106, RID)
			return c.String(http.StatusNotFound, "Analysis not found.\n")
		}
		return c.String(http.StatusInternalServerError, "Internal Error.\n")
	}
	return c.JSON(http.StatusOK, analysisResult)
}

// CreateNewSecurityTest inserts the given securityTest into SecurityTestCollection.
func CreateNewSecurityTest(c echo.Context) error {
	securityTest := types.SecurityTest{}
	err := c.Bind(&securityTest)
	if err != nil {
		log.Warning("CreateNewSecurityTest", "ANALYSIS", 108)
		return c.String(http.StatusBadRequest, "This is not a valid securityTest JSON.\n")
	}

	securityTestQuery := map[string]interface{}{"name": securityTest.Name}
	_, err = FindOneDBSecurityTest(securityTestQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewSecurityTest", "ANALYSIS", 109, securityTest.Name)
			return c.String(http.StatusConflict, "This securityTest is already in MongoDB.\n")
		}
		log.Error("CreateNewSecurityTest", "ANALYSIS", 1012, err)
	}

	err = InsertDBSecurityTest(securityTest)
	if err != nil {
		log.Error("CreateNewSecurityTest", "ANALYSIS", 2016, err)
		return c.String(http.StatusInternalServerError, "Internal error 2015.\n")
	}

	log.Info("CreateNewSecurityTest", "ANALYSIS", 18, securityTest.Name)
	return c.String(http.StatusCreated, "SecurityTest sucessfully created.\n")
}

// CreateNewRepository inserts the given repository into RepositoryCollection.
func CreateNewRepository(c echo.Context) error {
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		log.Warning("CreateNewRepository", "ANALYSIS", 101)
		return c.String(http.StatusBadRequest, "This is not a valid repository JSON.\n")
	}

	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	_, err = FindOneDBRepository(repositoryQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewRepository", "ANALYSIS", 110, repository.URL)
			return c.String(http.StatusConflict, "This repository is already in MongoDB.\n")
		}
		log.Error("CreateNewRepository", "ANALYSIS", 1013, err)
	}

	err = InsertDBRepository(repository)
	if err != nil {
		log.Error("CreateNewRepository", "ANALYSIS", 2015, err)
		return c.String(http.StatusInternalServerError, "Internal error 2015.\n")
	}

	log.Info("CreateNewRepository", "ANALYSIS", 17, repository.URL)
	return c.String(http.StatusCreated, "Repository sucessfully created.\n")
}

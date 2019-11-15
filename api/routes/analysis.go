// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"
	"time"

	"github.com/globocom/huskyCI/api/analysis"
	"github.com/globocom/huskyCI/api/auth"
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/token"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

var (
	tokenValidator token.TValidator
)

func init() {
	tokenCaller := token.TCaller{}
	hashGen := auth.Pbkdf2Caller{}
	tokenHandler := token.THandler{
		External: &tokenCaller,
		HashGen:  &hashGen,
	}
	tokenValidator = token.TValidator{
		TokenVerifier: &tokenHandler,
	}
}

const logActionReceiveRequest = "ReceiveRequest"
const logActionGetAnalysis = "GetAnalysis"
const logInfoAnalysis = "ANALYSIS"

// GetAnalysis returns the status of a given analysis given a RID.
func GetAnalysis(c echo.Context) error {

	RID := c.Param("id")
	attemptToken := c.Request().Header.Get("Husky-Token")
	if err := util.CheckMaliciousRID(RID, c); err != nil {
		return err
	}
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := apiContext.APIConfiguration.DBInstance.FindOneDBAnalysis(analysisQuery)
	if !tokenValidator.HasAuthorization(attemptToken, analysisResult.URL) {
		log.Error(logActionGetAnalysis, logInfoAnalysis, 1027, RID)
		reply := map[string]interface{}{"success": false, "error": "permission denied"}
		return c.JSON(http.StatusUnauthorized, reply)
	}
	if err != nil {
		if err == mgo.ErrNotFound || err.Error() == "No data found" {
			log.Warning(logActionGetAnalysis, logInfoAnalysis, 106, RID)
			reply := map[string]interface{}{"success": false, "error": "analysis not found"}
			return c.JSON(http.StatusNotFound, reply)
		}
		log.Error(logActionGetAnalysis, logInfoAnalysis, 1020, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	return c.JSON(http.StatusOK, analysisResult)
}

// ReceiveRequest receives the request and performs several checks before starting a new analysis.
func ReceiveRequest(c echo.Context) error {

	RID := c.Response().Header().Get(echo.HeaderXRequestID)
	attemptToken := c.Request().Header.Get("Husky-Token")

	// step-00: is this a valid JSON?
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		log.Error(logActionReceiveRequest, logInfoAnalysis, 1015, err)
		reply := map[string]interface{}{"success": false, "error": "invalid repository JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}
	if !tokenValidator.HasAuthorization(attemptToken, repository.URL) {
		log.Error("ReceivedRequest", logInfoAnalysis, 1027, RID)
		reply := map[string]interface{}{"success": false, "error": "permission denied"}
		return c.JSON(http.StatusUnauthorized, reply)
	}
	// step-01: Check malicious inputs
	sanitizedRepoURL, err := util.CheckValidInput(repository, c)
	if err != nil {
		return err
	}
	repository.URL = sanitizedRepoURL

	// step-02: is this repository already in MongoDB?
	repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL}
	_, err = apiContext.APIConfiguration.DBInstance.FindOneDBRepository(repositoryQuery)
	if err != nil {
		if err == mgo.ErrNotFound || err.Error() == "No data found" {
			// step-02-o1: repository not found! insert it into MongoDB
			repository.CreatedAt = time.Now()
			err = apiContext.APIConfiguration.DBInstance.InsertDBRepository(repository)
			if err != nil {
				log.Error(logActionReceiveRequest, logInfoAnalysis, 1010, err)
				reply := map[string]interface{}{"success": false, "error": "internal error"}
				return c.JSON(http.StatusInternalServerError, reply)
			}
		} else {
			// step-02-o2: another error searching for repositoryQuery
			log.Error(logActionReceiveRequest, logInfoAnalysis, 1013, err)
			reply := map[string]interface{}{"success": false, "error": "internal error"}
			return c.JSON(http.StatusInternalServerError, reply)
		}
	} else { // err == nil
		// step-03: repository found! does it have a running status analysis?
		analysisQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		analysisResult, err := apiContext.APIConfiguration.DBInstance.FindOneDBAnalysis(analysisQuery)
		if err != nil {
			if err == mgo.ErrNotFound || err.Error() == "No data found" {
				// nice! we can start this analysis!
			} else {
				// step-03-err: another error searching for analysisQuery
				log.Error(logActionReceiveRequest, logInfoAnalysis, 1009, err)
				reply := map[string]interface{}{"success": false, "error": "internal error"}
				return c.JSON(http.StatusInternalServerError, reply)
			}
		} else { // err == nil
			// step 03-a: Ops, this analysis is already running!
			if analysisResult.Status == "running" {
				log.Warning(logActionReceiveRequest, logInfoAnalysis, 104, analysisResult.URL)
				reply := map[string]interface{}{"success": false, "error": "an analysis is already in place for this URL and branch"}
				return c.JSON(http.StatusConflict, reply)
			}
		}
	}

	// step 04: lets start this analysis!
	log.Info(logActionReceiveRequest, logInfoAnalysis, 16, repository.Branch, repository.URL)
	go analysis.StartAnalysis(RID, repository)
	reply := map[string]interface{}{"success": true, "error": ""}
	return c.JSON(http.StatusCreated, reply)
}

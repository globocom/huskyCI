// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"fmt"
	"net/http"

	"github.com/globocom/huskyCI/api/analysis"
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/repository"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

// var (
// 	tokenValidator token.TValidator
// )

// func init() {
// 	tokenCaller := token.TCaller{}
// 	hashGen := auth.Pbkdf2Caller{}
// 	tokenHandler := token.THandler{
// 		External: &tokenCaller,
// 		HashGen:  &hashGen,
// 	}
// 	tokenValidator = token.TValidator{
// 		TokenVerifier: &tokenHandler,
// 	}
// }

const logActionReceiveRequest = "ReceiveRequest"
const logActionGetAnalysis = "GetAnalysis"
const logInfoAnalysis = "ANALYSIS"

// GetAnalysis returns the status of a given analysis given a RID.
func GetAnalysis(c echo.Context) error {

	analysisID := c.Param("id")

	// attemptToken := c.Request().Header.Get("Husky-Token")
	// if err := util.CheckMaliciousRID(RID, c); err != nil {
	// 	return err
	// }
	// if !tokenValidator.HasAuthorization(attemptToken, analysisResult.Repository.URL) {
	// 	log.Error(logActionGetAnalysis, logInfoAnalysis, 1027, RID)
	// 	reply := map[string]interface{}{"error": "permission denied"}
	// 	return c.JSON(http.StatusUnauthorized, reply)
	// }

	// analysisQuery := map[string]interface{}{"ID": analysisID}
	// analysisResult, err := apiContext.APIConfiguration.DBInstance.FindOneDBAnalysis(analysisQuery)
	// if err == mgo.ErrNotFound || err.Error() == "No data found" {
	// 	log.Warning(logActionGetAnalysis, logInfoAnalysis, 106, analysisID)
	// 	reply := map[string]interface{}{"error": "analysis not found"}
	// 	return c.JSON(http.StatusNotFound, reply)
	// }

	// log.Error(logActionGetAnalysis, logInfoAnalysis, 1020, err)
	// reply := map[string]interface{}{"error": "internal error"}
	// return c.JSON(http.StatusInternalServerError, reply)

	return c.JSON(http.StatusOK, analysisID)
}

// ReceiveRequest receives the request and performs several checks before starting a new analysis.
func ReceiveRequest(c echo.Context) error {

	// is this a valid token?
	// attemptToken := c.Request().Header.Get("Husky-Token")
	// if !tokenValidator.HasAuthorization(attemptToken, repository.URL) {
	// 	log.Error("ReceivedRequest", logInfoAnalysis, 1027, RID)
	// 	reply := map[string]interface{}{"success": false, "error": "permission denied"}
	// 	return c.JSON(http.StatusUnauthorized, reply)
	// }

	// is this a valid JSON?
	repositoryReceived := repository.Repository{}
	if err := c.Bind(&repositoryReceived); err != nil {
		log.Error(logActionReceiveRequest, logInfoAnalysis, 1015, err)
		reply := map[string]interface{}{"error": "invalid repository JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// is this a malicious payload?
	if err := repositoryReceived.CheckInput(); err != nil {
		log.Error(logActionReceiveRequest, logInfoAnalysis, 1015, err)
		reply := map[string]interface{}{"error": "invalid repository JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// is there an analysis already being running for this repository?
	analysisQuery := map[string]interface{}{"repositoryURL": repositoryReceived.URL, "repositoryBranch": repositoryReceived.Branch}
	analysisResult, err := apiContext.APIConfiguration.DBInstance.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		if err == mgo.ErrNotFound || err.Error() == "No data found" {
			// nice! we can start this analysis!
		} else {
			log.Error(logActionReceiveRequest, logInfoAnalysis, 1009, err)
			reply := map[string]interface{}{"error": "internal error"}
			return c.JSON(http.StatusInternalServerError, reply)
		}
	} else {
		// this analysis is already running!
		if analysisResult.Result.Status == "running" {
			log.Warning(logActionReceiveRequest, logInfoAnalysis, 104, analysisResult.Repository.URL)
			reply := map[string]interface{}{"error": "an analysis is already in place for this URL and branch"}
			return c.JSON(http.StatusConflict, reply)
		}
	}

	// good to go, start analysis!
	newAnalysis := analysis.New(&repositoryReceived)
	go newAnalysis.Start()

	// analysis started successfully
	log.Info(logActionReceiveRequest, logInfoAnalysis, 16, repositoryReceived.Branch, repositoryReceived.URL)
	startedAnalysisMessage := fmt.Sprintf("analysis %s created", newAnalysis.ID)
	reply := map[string]interface{}{"info": startedAnalysisMessage}
	return c.JSON(http.StatusCreated, reply)
}

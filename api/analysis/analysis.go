// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"time"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
)

// securityTestDoneCounter keeps track of all the security tests that were already done,
// it always starts with 0 and goes up to len(securityTests) - 1.
// Enry is not taken into consideration, as all other security tests are called by it.
var securityTestDoneCounter int

// StartAnalysis starts the analysis given a RID and a repository.
func StartAnalysis(RID string, repository types.Repository) {

	// step 0: create a new analysis struct
	newAnalysis := types.Analysis{
		RID:        RID,
		URL:        repository.URL,
		Branch:     repository.Branch,
		Status:     "running",
		Containers: make([]types.Container, 0),
		StartedAt:  time.Now(),
	}

	if repository.InternalDepURL != "" {
		newAnalysis.InternalDepURL = repository.InternalDepURL
	}

	// step 1: insert new analysis into MongoDB.
	err := db.InsertDBAnalysis(newAnalysis)
	if err != nil {
		log.Error("StartAnalysis", "ANALYSIS", 2011, err)
		return
	}

	// step 2: start enry and EnryStartAnalysis will start all others securityTests
	enryQuery := map[string]interface{}{"name": "enry"}
	enrySecurityTest, err := db.FindOneDBSecurityTest(enryQuery)
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

	timeout := time.After(90 * time.Minute)
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
				return
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
	err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateAnalysisQuery)
	return err
}

// monitorAnalysisUpdateStatus updates status and result of a given analysis.
func monitorAnalysisUpdateStatus(RID string) error {
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("monitorAnalysisUpdateStatus", "ANALYSIS", 2014, RID, err)
		return err
	}
	// analyze each cResult from each container to determine what is the value of analysis.Result
	finalResult := "passed"
	for _, container := range analysisResult.Containers {
		if container.CResult == "failed" || container.CResult == "error" {
			finalResult = "failed"
			break
		}
	}
	updateAnalysisQuery := bson.M{
		"$set": bson.M{
			"status":     "finished",
			"result":     finalResult,
			"finishedAt": time.Now(),
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateAnalysisQuery)
	if err != nil {
		log.Error("monitorAnalysisUpdateStatus", "ANALYSIS", 2007, err)
	}
	return err
}

// monitorAnalysisCheckStatus checks if an analysis has already finished and returns the correspoding boolean.
func monitorAnalysisCheckStatus(RID string) (bool, error) {
	securityTestDoneCounter = 0
	analysisFinished := false
	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("monitorAnalysisCheckStatus", "ANALYSIS", 2014, RID, err)
	}
	for _, container := range analysisResult.Containers {
		if container.CStatus != "finished" {
			analysisFinished = false
			break
		} else {
			// Enry must not be taken into account when verifying if the security tests have finished
			containerIsDoneButIsNotEnry := container.SecurityTest.Name != "enry" || ((container.SecurityTest.Name != "enry") && (container.CResult == ""))
			if containerIsDoneButIsNotEnry {
				securityTestDoneCounter++
				analysisFinished = true
			}
		}
	}
	// Makes sure all security tests found by Enry have finished
	if (len(analysisResult.SecurityTests) - 1) != securityTestDoneCounter {
		analysisFinished = false
	}
	return analysisFinished, err
}

func updateInfoAndResultBasedOnCID(cInfo, cResult, CID string) error {

	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
			"containers.$.cInfo":   cInfo,
		},
	}

	analysisQuery := map[string]interface{}{"containers.CID": CID}
	err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("updateCIDinfoAndResult", "ANALYSIS", 2007, err)
		return err
	}

	return nil
}

func updateHuskyCIResultsBasedOnRID(RID, securityTest string, securityTestResult interface{}) error {

	analysisQuery := map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("updateHuskyCIResultsBasedOnRID", "ANALYSIS", 2008, err)
		return err
	}

	switch securityTest {
	case "bandit":
		analysis.HuskyCIResults.PythonResults.HuskyCIBanditOutput = prepareHuskyCIBanditOutput(securityTestResult.(BanditOutput))
	case "retirejs":
		analysis.HuskyCIResults.JavaScriptResults.HuskyCIRetireJSOutput = prepareHuskyCIRetirejsOutput(securityTestResult.([]RetirejsOutput))
	case "brakeman":
		analysis.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput = prepareHuskyCIBrakemanResults(securityTestResult.(BrakemanOutput))
	case "gosec":
		analysis.HuskyCIResults.GoResults.HuskyCIGosecOutput = prepareHuskyCIGosecResults(securityTestResult.(GosecOutput))
	}

	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("updateHuskyCIResultsBasedOnRID", "ANALYSIS", 2007, err)
		return err
	}

	return nil
}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"time"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/securitytest"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
)

// securityTestDoneCounter keeps track of all the security tests that were already done,
// it always starts with 0 and goes up to len(securityTests) - 1.
// Enry is not taken into consideration, as all other security tests are called by it.
var securityTestDoneCounter int

// StartAnalysis starts the analysis given a RID and a repository.
func StartAnalysis(RID string, repository types.Repository) {

	// step 1: create a new analysis into MongoDB based on repository received
	if err := registerNewAnalysis(RID, repository); err != nil {
		return
	}
	log.Info("StartAnalysis", "ANALYSIS", 101, RID)

	// step 2: run enry as huskyCI initial step
	enryScan, err := securitytest.RunScanEnry(repository)
	if err != nil {
		return
	}

	// step 3: run generic and languages security tests based on enryScan result in parallel
	allScanResults := securitytest.RunAllScans(enryScan)

	// step 4: register all results found in MongoDB
	if err := registerFinishedAnalysis(RID, allScanResults); err != nil {
		return
	}
	log.Info("StartAnalysis", "ANALYSIS", 102, RID)

}

func registerNewAnalysis(RID string, repository types.Repository) error {

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

	if err := db.InsertDBAnalysis(newAnalysis); err != nil {
		log.Error("registerNewAnalysis", "ANALYSIS", 2011, err)
		return err
	}

	// log.Info("registerNewAnalysis", "ANALYSIS", 2012
	return nil
}

func registerFinishedAnalysis(RID string, allScanResults securitytest.AllScansResult) error {
	analysisQuery := map[string]interface{}{"RID": RID}
	updateAnalysisQuery := bson.M{
		"$set": bson.M{
			"status":         allScanResults.Status,
			"result":         allScanResults.FinalResult,
			"containers":     allScanResults.Containers,
			"huskyciresults": allScanResults.HuskyCIResults,
			"codes":          allScanResults.Codes,
			"finishedAt":     time.Now(),
		},
	}
	if err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateAnalysisQuery); err != nil {
		// log.Error("registerFinishedAnalysis", "ANALYSIS", 2011, err)
		return err
	}
	// log.Info("registerFinishedAnalysis", "ANALYSIS", 2012)
	return nil
}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"time"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/securitytest"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
)

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	RID string `bson:"RID" json:"RID"`
	// Repository      *repository.Repository        `bson:"repository" json:"repository"`
	Status     string    `bson:"status" json:"status"`
	Result     string    `bson:"result,omitempty" json:"result"`
	StartedAt  time.Time `bson:"startedAt" json:"startedAt"`
	FinishedAt time.Time `bson:"finishedAt" json:"finishedAt"`
	ErrorFound string    `bson:"errorFound,omitempty" json:"errorFound"`
	// Vulnerabilities []vulnerability.Vulnerability `bson:"vulnerabilities" json:"vulnerabilities"`
	// SecurityTests   []securitytest.SecurityTest   `bson:"securityTests" json:"securityTests"`
}

// New returns a new analysis struct based on a repository
// func New(repository *repository.Repository) *Analysis {
// 	return &Analysis{
// 		RID: uuid.New().String(),
// 		Repository: repository,
// 		StartedAt: time.Now(),
// 	}
// }

// Start runs a new analysis
func (a *Analysis) Start() error {

	// if err := a.Repository.Scan(); err != nil {
	// 	return err
	// }

	if err := a.checkCacheHit(); err != nil {
		return err
	}

	if err := a.setSecurityTests(); err != nil {
		return err
	}

	if err := a.startSecurityTests(); err != nil {
		return err
	}

	if err := a.registerInDatabase(); err != nil {
		return err
	}

	return nil
}

func (a *Analysis) checkCacheHit() error {

	var cacheHit bool

	if cacheHit {
		if err := a.registerInDatabase(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Analysis) setSecurityTests() error {

	// var securityTestsFound []securitytest.SecurityTest

	// securityTestsFound, err := securitytest.GetAllConfigsByLanguage(a.Repository.Languages)
	// if err != nil {
	// 	return err
	// }

	// a.SecurityTests = securityTestsFound

	return nil
}

func (a *Analysis) startSecurityTests() error {

	// for _, securityTest := range a.SecurityTests {
	// 	if err := securityTest.Container.Run(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (a *Analysis) registerInDatabase() error {
	return nil
}

const logActionStart = "StartAnalysis"
const logInfoAnalysis = "ANALYSIS"

// StartAnalysis starts the analysis given a RID and a repository.
func StartAnalysis(RID string, repository types.Repository) {

	// step 1: create a new analysis into MongoDB based on repository received
	if err := registerNewAnalysis(RID, repository); err != nil {
		return
	}
	log.Info(logActionStart, logInfoAnalysis, 101, RID)

	// step 2: run enry as huskyCI initial step
	enryScan := securitytest.SecTestScanInfo{}
	enryScan.SecurityTestName = "enry"
	allScansResults := securitytest.RunAllInfo{}

	defer func() {
		err := registerFinishedAnalysis(RID, &allScansResults)
		if err != nil {
			log.Error(logActionStart, logInfoAnalysis, 2011, err)
		}
	}()

	if err := enryScan.New(RID, repository.URL, repository.Branch, enryScan.SecurityTestName); err != nil {
		log.Error(logActionStart, logInfoAnalysis, 2011, err)
		return
	}
	if err := enryScan.Start(); err != nil {
		allScansResults.SetAnalysisError(err)
		return
	}

	// step 3: run generic and languages security tests based on enryScan result in parallel
	if err := allScansResults.Start(enryScan); err != nil {
		allScansResults.SetAnalysisError(err)
		return
	}

	log.Info("StartAnalysis", logInfoAnalysis, 102, RID)
}

func registerNewAnalysis(RID string, repository types.Repository) error {

	newAnalysis := types.Analysis{
		RID:       RID,
		URL:       repository.URL,
		Branch:    repository.Branch,
		Status:    "running",
		StartedAt: time.Now(),
	}

	if err := apiContext.APIConfiguration.DBInstance.InsertDBAnalysis(newAnalysis); err != nil {
		log.Error("registerNewAnalysis", logInfoAnalysis, 2011, err)
		return err
	}

	// log.Info("registerNewAnalysis", logInfoAnalysis, 2012
	return nil
}

func registerFinishedAnalysis(RID string, allScanResults *securitytest.RunAllInfo) error {
	analysisQuery := map[string]interface{}{"RID": RID}
	var errorString string
	if _, ok := allScanResults.ErrorFound.(error); ok {
		errorString = allScanResults.ErrorFound.Error()
	} else {
		errorString = ""
	}
	updateAnalysisQuery := bson.M{
		"status":         allScanResults.Status,
		"commitAuthors":  allScanResults.CommitAuthors,
		"result":         allScanResults.FinalResult,
		"containers":     allScanResults.Containers,
		"huskyciresults": allScanResults.HuskyCIResults,
		"codes":          allScanResults.Codes,
		"errorFound":     errorString,
		"finishedAt":     time.Now(),
	}

	if err := apiContext.APIConfiguration.DBInstance.UpdateOneDBAnalysisContainer(analysisQuery, updateAnalysisQuery); err != nil {
		log.Error("registerFinishedAnalysis", logInfoAnalysis, 2011, err)
		return err
	}
	return nil
}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
)

// EnryStartAnalysis checks the languages of a repository, update them into mongoDB, and starts corresponding new securityTests.
func EnryStartAnalysis(CID string, cOutput string, RID string) {

	// step 0.1: get analysis based on CID.
	analysisQuery := map[string]interface{}{"containers.CID": CID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2008, CID, err)
		return
	}

	// step 0.2: ERROR_CLONING or nil cOutput states that there were errors cloning a repository.
	if strings.Contains(cOutput, "ERROR_CLONING") || cOutput == "" {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "failed",
				"containers.$.cInfo":   errorOutput,
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("EnryStartAnalysis", "ENRY", 2007, err)
		}
		return
	}

	// step 1: get each language found in cOutput.
	mapLanguages := make(map[string][]interface{})
	err = json.Unmarshal([]byte(cOutput), &mapLanguages)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 1003, cOutput, err)
		return
	}
	repositoryLanguages := []types.Language{}
	newLanguage := types.Language{}
	for name, files := range mapLanguages {
		fs := []string{}
		for _, f := range files {
			if reflect.TypeOf(f).String() == "string" {
				fs = append(fs, f.(string))
			} else {
				log.Error("EnryStartAnalysis", "ENRY", 1004, err)
				return
			}
		}
		newLanguage = types.Language{
			Name:  name,
			Files: fs,
		}
		repositoryLanguages = append(repositoryLanguages, newLanguage)
	}

	// step 2: get all securityTests to be updated into RepositoryCollection and Analysiscollection.

	// step 2.1: querying MongoDB to gather up all securityTests that match (language=Generic and default=true).
	genericSecurityTestQuery := map[string]interface{}{"language": "Generic", "default": true}
	genericSecurityTests, err := db.FindAllDBSecurityTest(genericSecurityTestQuery)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2009, err)
		return
	}

	// step 2.2: querying MongoDB to gather up all securityTests that match (language=languageFound and default=true).
	newLanguageSecurityTests := []types.SecurityTest{}
	for _, language := range repositoryLanguages {
		languageSecurityTestQuery := map[string]interface{}{"language": language.Name, "default": true}
		languageSecurityTestResult, err := db.FindAllDBSecurityTest(languageSecurityTestQuery)
		if err == nil {
			newLanguageSecurityTests = append(newLanguageSecurityTests, languageSecurityTestResult...)
		}
	}

	allSecurityTests := append(genericSecurityTests, newLanguageSecurityTests...)

	// step 3: updating repository with all securityTests found.
	repositoryQuery := map[string]interface{}{"repositoryURL": analysis.URL, "repositoryBranch": analysis.Branch}
	updateRepositoryQuery := bson.M{
		"$set": bson.M{
			"securityTests": allSecurityTests,
			"languages":     repositoryLanguages,
		},
	}
	err = db.UpdateOneDBRepository(repositoryQuery, updateRepositoryQuery)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2010, err)
		return
	}

	// step 4: update analysis with the all securityTests found.
	analysis.SecurityTests = allSecurityTests
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2007, err)
		return
	}

	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cInfo": "Finished successfully.",
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2007, err)
	}

	// step 5: start all new securityTests.
	for _, securityTest := range newLanguageSecurityTests {
		// avoiding a loop here with this if condition.
		if securityTest.Name != "enry" {
			go DockerRun(RID, &analysis, securityTest)
		}
	}
}

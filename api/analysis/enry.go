// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// EnryCheckOutputFlow checks the languages of a repository, update them into mongoDB, and starts corresponding new securityTests.
func EnryCheckOutputFlow(CID string, cOutput string, RID string) {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	if errorClonning || cOutput == "" {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: get each language found in cOutput.
	repositoryLanguages, err := getLanguagesFromEnryOutput(cOutput)
	if err != nil {
		return
	}

	// step 3.1: get all generic securityTests from MongoDB
	genericSecurityTests, err := getAllSecurityTestsBasedOnLanguage("Generic")
	if err != nil {
		return
	}

	// step 3.2: get all securityTests based on each language found
	newLanguageSecurityTests := []types.SecurityTest{}
	for _, languageFound := range repositoryLanguages {
		languageFoundSecurityTests, err := getAllSecurityTestsBasedOnLanguage(languageFound.Language)
		if err != nil {
			return
		}
		newLanguageSecurityTests = append(newLanguageSecurityTests, languageFoundSecurityTests...)
	}

	// step 3.3: gather up generics securityTests + languages securityTests
	allSecurityTests := append(genericSecurityTests, newLanguageSecurityTests...)

	// step 4: update analysis with the all securityTests to be run in this repository
	analysisQuery := map[string]interface{}{"containers.CID": CID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2008, CID, err)
		return
	}
	analysis.SecurityTests = allSecurityTests
	analysis.Codes = repositoryLanguages
	if err := db.UpdateOneDBAnalysis(analysisQuery, analysis); err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 2007, err)
		return
	}

	// step 5: update enry cInfo and cResult
	if err := updateInfoAndResultBasedOnCID("Finished successfully.", "passed", CID); err != nil {
		return
	}

	// step 6: start all new securityTests.
	for _, securityTest := range newLanguageSecurityTests {
		// avoiding a loop here with this if condition.
		if securityTest.Name != "enry" {
			go DockerRun(RID, &analysis, securityTest)
		}
	}
}

func getLanguagesFromEnryOutput(cOutput string) ([]types.Code, error) {

	repositoryLanguages := []types.Code{}
	newLanguage := types.Code{}

	mapLanguages := make(map[string][]interface{})
	err := json.Unmarshal([]byte(cOutput), &mapLanguages)
	if err != nil {
		log.Error("EnryStartAnalysis", "ENRY", 1003, cOutput, err)
		return repositoryLanguages, err
	}

	for name, files := range mapLanguages {
		fs := []string{}
		for _, f := range files {
			if reflect.TypeOf(f).String() == "string" {
				fs = append(fs, f.(string))
			} else {
				log.Error("getLanguagesFromEnryOutput", "ENRY", 1004, err)
				return repositoryLanguages, errors.New("error mapping languages")
			}
		}
		newLanguage = types.Code{
			Language: name,
			Files:    fs,
		}
		repositoryLanguages = append(repositoryLanguages, newLanguage)
	}

	return repositoryLanguages, nil
}

func getAllSecurityTestsBasedOnLanguage(language string) ([]types.SecurityTest, error) {

	securityTests := []types.SecurityTest{}
	securityTestQuery := map[string]interface{}{"language": language, "default": true}

	securityTests, err := db.FindAllDBSecurityTest(securityTestQuery)
	if err != nil {
		log.Error("getAllSecurityTestsBasedOnLanguage", "ENRY", 2009, err)
		return securityTests, err
	}

	return securityTests, nil
}

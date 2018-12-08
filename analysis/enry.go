package analysis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/globocom/glbgelf"
	"github.com/globocom/huskyci/types"
	"gopkg.in/mgo.v2/bson"
)

// EnryStartAnalysis checks the languages of a repository, update them into mongoDB, and starts corresponding new securityTests.
func EnryStartAnalysis(CID string, cOutput string, RID string) {

	// step 0.1: get analysis based on CID.
	analysisQuery := map[string]interface{}{"containers.CID": CID}
	analysis, err := FindOneDBAnalysis(analysisQuery)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "EnryStartAnalysis",
			"info":   "ENRY"}, "ERROR", "Could not find analysis by this CID:", err)
		return
	}

	// step 0.2: ERROR_CLONING or nil cOutput states that there were errors cloning a repository.
	if strings.Contains(cOutput, "ERROR_CLONING") || cOutput == "" {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": errorOutput,
				"containers.$.cResult": "failed",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "EnryStartAnalysis",
				"info":   "ENRY"}, "ERROR", "Error updating AnalysisCollection (inside enry.go):", err)
		}
		return
	}

	// step 1: get each language found in cOutput.
	mapLanguages := make(map[string][]interface{})
	err = json.Unmarshal([]byte(cOutput), &mapLanguages)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "EnryStartAnalysis",
			"info":   "ENRY"}, "ERROR", "Unmarshall error (enry.go):", err)
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
				glbgelf.Logger.SendLog(map[string]interface{}{
					"action": "EnryStartAnalysis",
					"info":   "ENRY"}, "ERROR", "Error mapping languages.")
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
	genericSecurityTests, err := FindAllDBSecurityTest(genericSecurityTestQuery)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "EnryStartAnalysis",
			"info":   "ENRY"}, "ERROR", "Error finding securityTest (language=Generic and default=true):", err)
		return
	}

	// step 2.2: querying MongoDB to gather up all securityTests that match (language=languageFound and default=true).
	newLanguageSecurityTests := []types.SecurityTest{}
	for _, language := range repositoryLanguages {
		languageSecurityTestQuery := map[string]interface{}{"language": language.Name, "default": true}
		languageSecurityTestResult, err := FindOneDBSecurityTest(languageSecurityTestQuery)
		if err == nil {
			newLanguageSecurityTests = append(newLanguageSecurityTests, languageSecurityTestResult)
		} // else {} is OK to not find a securityTest by language.Name! for the future: log this error?
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
	err = UpdateOneDBRepository(repositoryQuery, updateRepositoryQuery)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "EnryStartAnalysis",
			"info":   "ENRY"}, "ERROR", "Could not update repository's securityTests:", err)
		return
	}

	// step 4: update analysis with the all securityTests found.
	analysis.SecurityTests = allSecurityTests
	err = UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "EnryStartAnalysis",
			"info":   "ENRY"}, "ERROR", "Error updating AnalysisCollection:", err)
		return
	}

	// step 5: start all new securityTests.
	for _, securityTest := range newLanguageSecurityTests {
		// avoiding a loop here with this if condition.
		if securityTest.Name != "enry" {
			go DockerRun(RID, &analysis, securityTest)
		}
	}
}

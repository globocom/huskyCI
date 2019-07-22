// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

//RetirejsOutput is the struct that holds issues, messages and errors found on a Retire scan.
type RetirejsOutput struct {
	RetirejsResult []RetirejsResult `json:"results"`
}

//RetirejsResult is a struct that holds the scanned results.
type RetirejsResult struct {
	Component       string                    `json:"component"`
	Version         string                    `json:"version"`
	Level           int                       `json:"level"`
	Vulnerabilities []RetireJSVulnerabilities `json:"vulnerabilities"`
}

//RetireJSVulnerabilities is a struct that holds the vulnerabilities found on a scan.
type RetireJSVulnerabilities struct {
	Info        []string                         `json:"info"`
	Severity    string                           `json:"severity"`
	Identifiers RetireJSVulnerabilityIdentifiers `json:"identifiers"`
}

//RetireJSVulnerabilityIdentifiers is a struct that holds identifiying information on a vulnerability found.
type RetireJSVulnerabilityIdentifiers struct {
	Summary string
}

//RetirejsStartAnalysis analyses the output from RetireJS and sets cResult basdes on it.
func RetirejsStartAnalysis(CID string, cOutput string, RID string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.1: error cloning repository!
	if strings.Contains(cOutput, "ERROR_CLONING") {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   errorOutput,
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("RetirejsStartAnalysis", "RETIREJS", 2007, err)
		}
		return
	}

	failedRunning := strings.Contains(cOutput, "ERROR_RUNNING_RETIREJS")
	if failedRunning {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   "Internal error running NPM Audit.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("RetirejsStartAnalysis", "RETIREJS", 2007, err)
			return
		}

		// step 4: get updated analysis based on its RID
		analysisQuery = map[string]interface{}{"RID": RID}
		analysis, err := db.FindOneDBAnalysis(analysisQuery)
		if err != nil {
			log.Error("RetirejsStartAnalysis", "RETIREJS", 2008, CID, err)
			return
		}

		// step 5: finally, update analysis with huskyCI results
		retireJSOutput := []RetirejsOutput{}
		analysis.HuskyCIResults.JavaScriptResults.HuskyCIRetireJSOutput = prepareHuskyCIRetirejsOutput(retireJSOutput, failedRunning)
		err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
		if err != nil {
			log.Error("RetirejsStartAnalysis", "RETIREJS", 2007, err)
			return
		}

		return
	}

	// step 1: Unmarshall cOutput into RetireOutput struct.
	retirejsOutput := []RetirejsOutput{}
	err := json.Unmarshal([]byte(cOutput), &retirejsOutput)
	if err != nil {
		log.Error("RetirejsStartAnalysis", "RETIREJS", 1014, cOutput, err)
		return
	}

	// step 1.1: Sets the container output to "No issues found" if RetirejsIssues returns an empty slice
	if len(retirejsOutput) == 0 {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "passed",
				"containers.$.cInfo":   "No issues found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("RetirejsStartAnalysis", "RETIREJS", 2007, err)
		}
		return
	}

	// step 2: find Vulnerabilities that have severity "medium" or "high".
	cResult = "passed"
	for _, output := range retirejsOutput {
		for _, result := range output.RetirejsResult {
			for _, vulnerability := range result.Vulnerabilities {
				if vulnerability.Severity == "high" || vulnerability.Severity == "medium" {
					cResult = "failed"
					break
				}
			}
		}
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	issueMessage := "No issues found."
	if cResult != "passed" {
		issueMessage = "Issues found."
	}
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
			"containers.$.cInfo":   issueMessage,
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("RetirejsStartAnalysis", "RETIREJS", 2007, err)
	}

	// step 4: get updated analysis based on its RID
	analysisQuery = map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("RetirejsStartAnalysis", "RETIREJS", 2008, CID, err)
		return
	}

	// step 5: finally, update analysis with huskyCI results
	analysis.HuskyCIResults.JavaScriptResults.HuskyCIRetireJSOutput = prepareHuskyCIRetirejsOutput(retirejsOutput, failedRunning)
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("RetirejsStartAnalysis", "RETIREJS", 2007, err)
		return
	}

}

// prepareHuskyCIRetirejsOutput will prepare Retirejs output to be added into JavaScriptResults struct
func prepareHuskyCIRetirejsOutput(retirejsOutput []RetirejsOutput, failedRunning bool) types.HuskyCIRetireJSOutput {

	var huskyCIretireJSResults types.HuskyCIRetireJSOutput
	var huskyCIretireJSResultsFinal types.HuskyCIRetireJSOutput

	if failedRunning {
		retirejsVuln := types.HuskyCIVulnerability{}
		retirejsVuln.Language = "JavaScript"
		retirejsVuln.SecurityTool = "RetireJS"
		retirejsVuln.Severity = "low"
		retirejsVuln.Details = "It looks like your project doesn't have package.json or yarn.lock. huskyCI was not able to run RetireJS properly."

		huskyCIretireJSResults.LowVulnsNpmRetireJS = append(huskyCIretireJSResults.LowVulnsNpmRetireJS, retirejsVuln)

		return huskyCIretireJSResults
	}

	for _, output := range retirejsOutput {
		for _, result := range output.RetirejsResult {
			for _, vulnerability := range result.Vulnerabilities {
				retirejsVuln := types.HuskyCIVulnerability{}
				retirejsVuln.Language = "JavaScript"
				retirejsVuln.SecurityTool = "RetireJS"
				retirejsVuln.Severity = vulnerability.Severity
				retirejsVuln.Code = result.Component
				retirejsVuln.Version = result.Version
				for _, info := range vulnerability.Info {
					retirejsVuln.Details = retirejsVuln.Details + info + "\n"
				}
				retirejsVuln.Details = retirejsVuln.Details + vulnerability.Identifiers.Summary

				switch retirejsVuln.Severity {
				case "low":
					huskyCIretireJSResults.LowVulnsNpmRetireJS = append(huskyCIretireJSResults.LowVulnsNpmRetireJS, retirejsVuln)
				case "medium":
					huskyCIretireJSResults.MediumVulnsRetireJS = append(huskyCIretireJSResults.MediumVulnsRetireJS, retirejsVuln)
				case "high":
					huskyCIretireJSResults.HighVulnsRetireJS = append(huskyCIretireJSResults.HighVulnsRetireJS, retirejsVuln)
				}
			}
		}
	}

	huskyCIretireJSResultsFinal.LowVulnsNpmRetireJS = util.CountRetireJSOccurrences(huskyCIretireJSResults.LowVulnsNpmRetireJS)
	huskyCIretireJSResultsFinal.MediumVulnsRetireJS = util.CountRetireJSOccurrences(huskyCIretireJSResults.MediumVulnsRetireJS)
	huskyCIretireJSResultsFinal.HighVulnsRetireJS = util.CountRetireJSOccurrences(huskyCIretireJSResults.HighVulnsRetireJS)

	return huskyCIretireJSResultsFinal
}

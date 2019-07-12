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
func RetirejsStartAnalysis(CID string, cOutput string) {

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

	if strings.Contains(cOutput, "ERROR_RUNNING_RETIREJS") {
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

	return

}

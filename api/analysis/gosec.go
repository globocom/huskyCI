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
	"gopkg.in/mgo.v2/bson"
)

// GosecOutput is the struct that holds all data from Gosec output.
type GosecOutput struct {
	GosecIssues []GosecIssue `json:"Issues"`
	GosecStats  GosecStats   `json:"Stats"`
}

// GosecIssue is the struct that holds all issues from Gosec output.
type GosecIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// GosecStats is the struct that holds all stats from Gosec output.
type GosecStats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	Nosec int `json:"nosec"`
	Found int `json:"found"`
}

// GosecStartAnalysis analyses the output from Gosec and sets a cResult based on it.
func GosecStartAnalysis(CID string, cOutput string, RID string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.1: nil cOutput states that no Issues were found.
	if cOutput == "" {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "passed",
				"containers.$.cInfo":   "No issues found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("GosecStartAnalysis", "GOSEC", 2007, "Step 0.1 ", err)
		}
		return
	}

	// step 0.2: error cloning repository!
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
			log.Error("GosecStartAnalysis", "GOSEC", 2007, "Step 0.2 ", err)
		}
		return
	}

	// step 1: Unmarshall cOutput into GosecOutput struct.
	gosecOutput := GosecOutput{}
	err := json.Unmarshal([]byte(cOutput), &gosecOutput)
	if err != nil {
		log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return
	}

	// step 2: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult = "warning"
	for _, issue := range gosecOutput.GosecIssues {
		if (issue.Severity == "HIGH" || issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			cResult = "failed"
			break
		}
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	issueMessage := "Warning found."
	if cResult != "warning" {
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
		log.Error("GosecStartAnalysis", "GOSEC", 2007, "Step 3 ", err)
		return
	}

	// step 4: get updated analysis based on its RID
	analysisQuery = map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("GosecStartAnalysis", "GOSEC", 2008, CID, err)
		return
	}

	// step 5: finally, update analysis with huskyCI results
	analysis.HuskyCIResults.GoResults.HuskyCIGosecOutput = prepareHuskyCIGosecResults(gosecOutput)
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("GosecStartAnalysis", "GOSEC", 2007, err)
		return
	}
}

// prepareHuskyCIGosecResults will prepare Gosec output to be added into goResults struct
func prepareHuskyCIGosecResults(gosecOutput GosecOutput) types.HuskyCIGosecOutput {

	var huskyCIgosecResults types.HuskyCIGosecOutput

	for _, issue := range gosecOutput.GosecIssues {
		gosecVuln := types.HuskyCIVulnerability{}
		gosecVuln.Language = "Go"
		gosecVuln.SecurityTool = "GoSec"
		gosecVuln.Severity = issue.Severity
		gosecVuln.Confidence = issue.Confidence
		gosecVuln.Details = issue.Details
		gosecVuln.File = issue.File
		gosecVuln.Line = issue.Line
		gosecVuln.Code = issue.Code

		switch gosecVuln.Severity {
		case "LOW":
			huskyCIgosecResults.LowVulnsGosec = append(huskyCIgosecResults.LowVulnsGosec, gosecVuln)
		case "MEDIUM":
			huskyCIgosecResults.MediumVulnsGosec = append(huskyCIgosecResults.MediumVulnsGosec, gosecVuln)
		case "HIGH":
			huskyCIgosecResults.HighVulnsGosec = append(huskyCIgosecResults.HighVulnsGosec, gosecVuln)
		}
	}

	return huskyCIgosecResults
}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
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

// GosecCheckOutputFlow analyses the output from Gosec and sets a cResult based on it.
func GosecCheckOutputFlow(CID string, cOutput string, RID string) {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	if errorClonning {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: nil cOutput states that no Issues were found.
	if cOutput == "" {
		if err := updateInfoAndResultBasedOnCID("No issues found.", "passed", CID); err != nil {
			return
		}
		return
	}

	// step 3: Unmarshall cOutput into GosecOutput struct.
	gosecOutput := GosecOutput{}
	err := json.Unmarshal([]byte(cOutput), &gosecOutput)
	if err != nil {
		log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return
	}

	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult := "warning"
	issueMessage := "Warning found."
	for _, issue := range gosecOutput.GosecIssues {
		if (issue.Severity == "HIGH" || issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			cResult = "failed"
			issueMessage = "Issues found."
			break
		}
	}
	if err := updateInfoAndResultBasedOnCID(issueMessage, cResult, CID); err != nil {
		return
	}

	// step 5: finally, update analysis with huskyCI results
	if err := updateHuskyCIResultsBasedOnRID(RID, "gosec", gosecOutput); err != nil {
		return
	}
}

// prepareHuskyCIGosecResults will prepare Gosec output to be added into goResults struct
func prepareHuskyCIGosecResults(gosecOutput GosecOutput) types.HuskyCISecurityTestOutput {

	var huskyCIgosecResults types.HuskyCISecurityTestOutput

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
			huskyCIgosecResults.LowVulns = append(huskyCIgosecResults.LowVulns, gosecVuln)
		case "MEDIUM":
			huskyCIgosecResults.MediumVulns = append(huskyCIgosecResults.MediumVulns, gosecVuln)
		case "HIGH":
			huskyCIgosecResults.HighVulns = append(huskyCIgosecResults.HighVulns, gosecVuln)
		}
	}

	return huskyCIgosecResults
}

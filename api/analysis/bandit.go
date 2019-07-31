// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// BanditOutput is the struct that holds all data from Bandit output.
type BanditOutput struct {
	Errors  json.RawMessage `json:"errors"`
	Results []Result        `json:"results"`
}

// Result is the struct that holds detailed information of issues from Bandit output.
type Result struct {
	Code            string `json:"code"`
	Filename        string `json:"filename"`
	IssueConfidence string `json:"issue_confidence"`
	IssueSeverity   string `json:"issue_severity"`
	IssueText       string `json:"issue_text"`
	LineNumber      int    `json:"line_number"`
	LineRange       []int  `json:"line_range"`
	TestID          string `json:"test_id"`
	TestName        string `json:"test_name"`
}

// BanditCheckOutputFlow analyses the output from Bandit and sets a cResult based on it.
func BanditCheckOutputFlow(CID string, cOutput string, RID string) {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	if errorClonning {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: get Bandit output to be checked
	var banditResult BanditOutput
	if err := json.Unmarshal([]byte(cOutput), &banditResult); err != nil {
		log.Error("BanditStartAnalysis", "BANDIT", 1006, cOutput, err)
		return
	}

	// step 3: sets the container output to "No issues found" if banditResult returns an empty slice
	if len(banditResult.Results) == 0 {
		if err := updateInfoAndResultBasedOnCID("No issues found.", "passed", CID); err != nil {
			return
		}
		return
	}

	// step 4: verify if there was any error in the analysis.
	if banditResult.Errors != nil {
		if err := updateInfoAndResultBasedOnCID("Internal error running Bandit.", "error", CID); err != nil {
			return
		}
	}

	// step 5: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult := "passed"
	issueMessage := "No issues found."
	for _, issue := range banditResult.Results {
		if (issue.IssueSeverity == "HIGH" || issue.IssueSeverity == "MEDIUM") && issue.IssueConfidence == "HIGH" {
			cResult = "failed"
			issueMessage = "Issues found."
			break
		}
	}
	if err := updateInfoAndResultBasedOnCID(issueMessage, cResult, CID); err != nil {
		return
	}

	// step 6: finally, update analysis with huskyCI results
	if err := updateHuskyCIResultsBasedOnRID(RID, "bandit", banditResult); err != nil {
		return
	}
}

// prepareHuskyCIBanditOutput will prepare Bandit output to be added into pythonResults struct
func prepareHuskyCIBanditOutput(banditOutput BanditOutput) types.HuskyCISecurityTestOutput {

	var huskyCIbanditResults types.HuskyCISecurityTestOutput

	for _, issue := range banditOutput.Results {
		banditVuln := types.HuskyCIVulnerability{}
		banditVuln.Language = "Python"
		banditVuln.SecurityTool = "Bandit"
		banditVuln.Severity = issue.IssueSeverity
		banditVuln.Confidence = issue.IssueConfidence
		banditVuln.Details = issue.IssueText
		banditVuln.File = issue.Filename
		banditVuln.Line = strconv.Itoa(issue.LineNumber)
		banditVuln.Code = issue.Code

		switch banditVuln.Severity {
		case "LOW":
			huskyCIbanditResults.LowVulns = append(huskyCIbanditResults.LowVulns, banditVuln)
		case "MEDIUM":
			huskyCIbanditResults.MediumVulns = append(huskyCIbanditResults.MediumVulns, banditVuln)
		case "HIGH":
			huskyCIbanditResults.HighVulns = append(huskyCIbanditResults.HighVulns, banditVuln)
		}
	}

	return huskyCIbanditResults
}

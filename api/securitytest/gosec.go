// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

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

func analyzeGosec(gosecScan *SecTestScanInfo) error {

	goSecOutput := GosecOutput{}
	gosecScan.FinalOutput = goSecOutput

	// nil cOutput states that no Issues were found.
	if gosecScan.Container.COutput == "" {
		gosecScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a GosecOutput struct.
	if err := json.Unmarshal([]byte(gosecScan.Container.COutput), &goSecOutput); err != nil {
		log.Error("analyzeGosec", "GOSEC", 1002, gosecScan.Container.COutput, err)
		gosecScan.ErrorFound = err
		gosecScan.prepareContainerAfterScan()
		return err
	}
	gosecScan.FinalOutput = goSecOutput

	// check results and prepare all vulnerabilities found
	gosecScan.prepareGosecVulns()
	gosecScan.prepareContainerAfterScan()
	return nil
}

func (gosecScan *SecTestScanInfo) prepareGosecVulns() {

	huskyCIgosecResults := types.HuskyCISecurityTestOutput{}
	gosecOutput := gosecScan.FinalOutput.(GosecOutput)

	for _, issue := range gosecOutput.GosecIssues {
		gosecVuln := types.HuskyCIVulnerability{}
		gosecVuln.Language = "Go"
		gosecVuln.SecurityTool = "GoSec"
		gosecVuln.Title = issue.Details
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

	for i := 0; i < gosecOutput.GosecStats.Nosec; i++ {
		gosecVuln := types.HuskyCIVulnerability{}
		huskyCIgosecResults.NoSecVulns = append(huskyCIgosecResults.NoSecVulns, gosecVuln)
	}

	gosecScan.Vulnerabilities = huskyCIgosecResults
}

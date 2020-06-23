// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"strconv"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// BanditOutput is the struct that holds all data from Bandit output.
type BanditOutput struct {
	Results []Result `json:"results"`
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

func analyzeBandit(banditScan *SecTestScanInfo) error {

	banditOutput := BanditOutput{}

	// Unmarshall rawOutput into finalOutput, that is a Bandit struct.
	if err := json.Unmarshal([]byte(banditScan.Container.COutput), &banditOutput); err != nil {
		log.Error("analyzeBandit", "BANDIT", 1006, banditScan.Container.COutput, err)
		banditScan.ErrorFound = err
		return err
	}
	banditScan.FinalOutput = banditOutput

	// an empty Results slice states that no Issues were found.
	if len(banditOutput.Results) == 0 {
		banditScan.prepareContainerAfterScan()
		return nil
	}
	// check results and prepare all vulnerabilities found
	banditScan.prepareBanditVulns()
	banditScan.prepareContainerAfterScan()
	return nil
}

func (banditScan *SecTestScanInfo) prepareBanditVulns() {

	huskyCIbanditResults := types.HuskyCISecurityTestOutput{}
	banditOutput := banditScan.FinalOutput.(BanditOutput)

	for _, issue := range banditOutput.Results {
		banditVuln := types.HuskyCIVulnerability{}
		banditVuln.Language = "Python"
		banditVuln.SecurityTool = "Bandit"
		noHuskyInLine := util.VerifyNoHusky(issue.Code, issue.LineNumber, banditVuln.SecurityTool)
		if noHuskyInLine {
			issue.IssueSeverity = "NOSEC"
		}
		banditVuln.Severity = issue.IssueSeverity
		banditVuln.Confidence = issue.IssueConfidence
		banditVuln.Title = issue.IssueText
		banditVuln.Details = issue.IssueText
		banditVuln.File = issue.Filename
		banditVuln.Line = strconv.Itoa(issue.LineNumber)
		banditVuln.Code = issue.Code

		switch banditVuln.Severity {
		case "NOSEC":
			huskyCIbanditResults.NoSecVulns = append(huskyCIbanditResults.NoSecVulns, banditVuln)
		case "LOW":
			huskyCIbanditResults.LowVulns = append(huskyCIbanditResults.LowVulns, banditVuln)
		case "MEDIUM":
			huskyCIbanditResults.MediumVulns = append(huskyCIbanditResults.MediumVulns, banditVuln)
		case "HIGH":
			huskyCIbanditResults.HighVulns = append(huskyCIbanditResults.HighVulns, banditVuln)
		}
	}

	banditScan.Vulnerabilities = huskyCIbanditResults
}

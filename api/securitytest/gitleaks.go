// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// GitleaksOutput is the struct that holds all data from Gitleaks output.
type GitleaksOutput []GitLeaksIssue

// GitLeaksIssue is the struct that holds all isssues from Gitleaks output.
type GitLeaksIssue struct {
	Line          string `json:"line"`
	Commit        string `json:"commit"`
	Offender      string `json:"offender"`
	Rule          string `json:"rule"`
	Info          string `json:"info"`
	CommitMessage string `json:"commitMsg"`
	Author        string `json:"author"`
	Email         string `json:"email"`
	File          string `json:"file"`
	Repository    string `json:"repo"`
	Date          string `json:"date"`
	Tags          string `json:"tags"`
	Severity      string `json:"severity"`
}

func analyseGitleaks(gitleaksScan *SecTestScanInfo) error {
	gitLeaksOutput := GitleaksOutput{}
	gitleaksScan.FinalOutput = gitLeaksOutput

	// nil cOutput states that no Issues were found.
	if gitleaksScan.Container.COutput == "" {
		gitleaksScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a GitleaksOutput struct.
	if err := json.Unmarshal([]byte(gitleaksScan.Container.COutput), &gitLeaksOutput); err != nil {
		log.Error("analyzeGitleaks", "GITLEAKS", 1038, gitleaksScan.Container.COutput, err)
		gitleaksScan.ErrorFound = err
		gitleaksScan.prepareContainerAfterScan()
		return err
	}
	gitleaksScan.FinalOutput = gitLeaksOutput

	// check results and prepare all vulnerabilities found
	gitleaksScan.prepareGitleaksVulns()
	gitleaksScan.prepareContainerAfterScan()
	return nil
}

func (gitleaksScan *SecTestScanInfo) prepareGitleaksVulns() {

	huskyCIgitleaksResults := types.HuskyCISecurityTestOutput{}
	gitleaksOutput := gitleaksScan.FinalOutput.(GitleaksOutput)

	for _, issue := range gitleaksOutput {
		gitleaksVuln := types.HuskyCIVulnerability{}
		gitleaksVuln.SecurityTool = "GitLeaks"
		gitleaksVuln.Severity = "medium"
		gitleaksVuln.Details = issue.Info + " @ " + issue.CommitMessage
		gitleaksVuln.File = issue.File
		gitleaksVuln.Code = issue.Line

		huskyCIgitleaksResults.MediumVulns = append(huskyCIgitleaksResults.MediumVulns, gitleaksVuln)
	}

	gitleaksScan.Vulnerabilities = huskyCIgitleaksResults
}

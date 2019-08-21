// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/db"
	huskydocker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// BanditScan holds all information needed for a Bandit scan.
type BanditScan struct {
	RID             string
	CID             string
	URL             string
	Branch          string
	Image           string
	Command         string
	RawOutput       string
	FinalOutput     BanditOutput
	Vulnerabilities types.HuskyCISecurityTestOutput
}

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

func newScanBandit(URL, branch, command string) BanditScan {
	return BanditScan{
		Image:   "huskyci/bandit",
		URL:     URL,
		Branch:  branch,
		Command: util.HandleCmd(URL, branch, command),
	}
}

func initBandit(enryScan EnryScan, allScansResult *AllScansResult) error {
	banditScan, banditContainer, err := runScanBandit(enryScan.URL, enryScan.Branch)
	if err != nil {
		return err
	}

	for _, highVuln := range banditScan.Vulnerabilities.HighVulns {
		allScansResult.HuskyCIResults.PythonResults.HuskyCIBanditOutput.HighVulns = append(allScansResult.HuskyCIResults.PythonResults.HuskyCIBanditOutput.HighVulns, highVuln)
	}
	for _, mediumVuln := range banditScan.Vulnerabilities.MediumVulns {
		allScansResult.HuskyCIResults.PythonResults.HuskyCIBanditOutput.MediumVulns = append(allScansResult.HuskyCIResults.PythonResults.HuskyCIBanditOutput.MediumVulns, mediumVuln)
	}
	for _, lowVuln := range banditScan.Vulnerabilities.LowVulns {
		allScansResult.HuskyCIResults.PythonResults.HuskyCIBanditOutput.LowVulns = append(allScansResult.HuskyCIResults.PythonResults.HuskyCIBanditOutput.LowVulns, lowVuln)
	}

	allScansResult.FinalResult = banditContainer.CResult
	allScansResult.Status = banditContainer.CStatus
	allScansResult.Containers = append(allScansResult.Containers, banditContainer)
	return nil
}

func runScanBandit(URL, branch string) (BanditScan, types.Container, error) {
	banditScan := BanditScan{}
	banditContainer, err := newContainerBandit()
	if err != nil {
		return banditScan, banditContainer, err
	}
	banditScan = newScanBandit(URL, branch, banditContainer.SecurityTest.Cmd)
	if err := banditScan.startBandit(); err != nil {
		return banditScan, banditContainer, err
	}

	banditScan.prepareContainerAfterScanBandit(&banditContainer)
	return banditScan, banditContainer, nil
}

func (banditScan *BanditScan) startBandit() error {
	if err := banditScan.dockerRunBandit(); err != nil {
		return err
	}
	if err := banditScan.analyzeBandit(); err != nil {
		return err
	}
	// log.Info("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
	return nil
}

func (banditScan *BanditScan) dockerRunBandit() error {
	CID, cOutput, err := huskydocker.DockerRun(banditScan.Image, banditScan.Command)
	if err != nil {
		// log.Error("DockerRun", "DOCKERRUN", 3013, err)
		return err
	}
	banditScan.CID = CID
	banditScan.RawOutput = cOutput
	return nil
}

func (banditScan *BanditScan) analyzeBandit() error {
	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(banditScan.RawOutput, "ERROR_CLONING")
	if errorClonning {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return errors.New("error clonning")
	}
	// step 2: nil cOutput states that no Issues were found.
	if banditScan.RawOutput == "" {
		return nil
	}
	// step 3: Unmarshall rawOutput into finalOutput, that is a GosecOutput struct.
	if err := json.Unmarshal([]byte(banditScan.RawOutput), &banditScan.FinalOutput); err != nil {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return err
	}
	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	banditScan.prepareBanditOutput(banditScan.FinalOutput)
	return nil
}

func (banditScan *BanditScan) prepareBanditOutput(banditOutput BanditOutput) {
	huskyCIbanditResults := types.HuskyCISecurityTestOutput{}

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
	banditScan.Vulnerabilities = huskyCIbanditResults
}

func (banditScan *BanditScan) prepareContainerAfterScanBandit(banditContainer *types.Container) {
	if len(banditScan.Vulnerabilities.MediumVulns) > 0 || len(banditScan.Vulnerabilities.HighVulns) > 0 {
		banditContainer.CInfo = "Issues found."
		banditContainer.CResult = "failed"
	} else if len(banditScan.Vulnerabilities.LowVulns) > 0 && (len(banditScan.Vulnerabilities.MediumVulns) == 0 || len(banditScan.Vulnerabilities.HighVulns) == 0) {
		banditContainer.CInfo = "Warnings found."
		banditContainer.CResult = "passed"
	}
	banditContainer.CStatus = "finished"
	banditContainer.CID = banditScan.CID
	banditContainer.COutput = banditScan.RawOutput
	banditContainer.FinishedAt = time.Now()
}

func newContainerBandit() (types.Container, error) {
	banditContainer := types.Container{}
	banditQuery := map[string]interface{}{"name": "bandit"}
	banditSecurityTest, err := db.FindOneDBSecurityTest(banditQuery)
	if err != nil {
		return banditContainer, err
	}
	return types.Container{
		SecurityTest: banditSecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/db"
	huskydocker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// GosecScan holds all information needed for a gosec scan.
type GosecScan struct {
	RID             string
	CID             string
	URL             string
	Branch          string
	Image           string
	Command         string
	RawOutput       string
	FinalOutput     GosecOutput
	Vulnerabilities types.HuskyCISecurityTestOutput
}

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

func newScanGosec(URL, branch, command string) GosecScan {
	return GosecScan{
		Image:   "huskyci/gosec",
		URL:     URL,
		Branch:  branch,
		Command: util.HandleCmd(URL, branch, command),
	}
}

func initGoSec(enryScan EnryScan, allScansResult *AllScansResult) error {
	gosecScan, gosecContainer, err := runScanGosec(enryScan.URL, enryScan.Branch)
	if err != nil {
		return err
	}

	for _, highVuln := range gosecScan.Vulnerabilities.HighVulns {
		allScansResult.HuskyCIResults.GoResults.HuskyCIGosecOutput.HighVulns = append(allScansResult.HuskyCIResults.GoResults.HuskyCIGosecOutput.HighVulns, highVuln)
	}
	for _, mediumVuln := range gosecScan.Vulnerabilities.MediumVulns {
		allScansResult.HuskyCIResults.GoResults.HuskyCIGosecOutput.MediumVulns = append(allScansResult.HuskyCIResults.GoResults.HuskyCIGosecOutput.MediumVulns, mediumVuln)
	}
	for _, lowVuln := range gosecScan.Vulnerabilities.LowVulns {
		allScansResult.HuskyCIResults.GoResults.HuskyCIGosecOutput.LowVulns = append(allScansResult.HuskyCIResults.GoResults.HuskyCIGosecOutput.LowVulns, lowVuln)
	}

	allScansResult.FinalResult = gosecContainer.CResult
	allScansResult.Status = gosecContainer.CStatus
	allScansResult.Containers = append(allScansResult.Containers, gosecContainer)
	return nil
}

func runScanGosec(URL, branch string) (GosecScan, types.Container, error) {
	gosecScan := GosecScan{}
	gosecContainer, err := newContainerGosec()
	if err != nil {
		return gosecScan, gosecContainer, err
	}
	gosecScan = newScanGosec(URL, branch, gosecContainer.SecurityTest.Cmd)
	if err := gosecScan.startGosec(); err != nil {
		return gosecScan, gosecContainer, err
	}

	gosecScan.prepareContainerAfterScanGosec(&gosecContainer)
	return gosecScan, gosecContainer, nil
}

func (gosecScan *GosecScan) startGosec() error {
	if err := gosecScan.dockerRunGosec(); err != nil {
		return err
	}
	if err := gosecScan.analyzeGosec(); err != nil {
		return err
	}
	// log.Info("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
	return nil
}

func (gosecScan *GosecScan) dockerRunGosec() error {
	CID, cOutput, err := huskydocker.DockerRun(gosecScan.Image, gosecScan.Command)
	if err != nil {
		// log.Error("DockerRun", "DOCKERRUN", 3013, err)
		return err
	}
	gosecScan.CID = CID
	gosecScan.RawOutput = cOutput
	return nil
}

func (gosecScan *GosecScan) analyzeGosec() error {
	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(gosecScan.RawOutput, "ERROR_CLONING")
	if errorClonning {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return errors.New("error clonning")
	}
	// step 2: nil cOutput states that no Issues were found.
	if gosecScan.RawOutput == "" {
		return nil
	}
	// step 3: Unmarshall rawOutput into finalOutput, that is a GosecOutput struct.
	if err := json.Unmarshal([]byte(gosecScan.RawOutput), &gosecScan.FinalOutput); err != nil {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return err
	}
	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	gosecScan.prepareGosecOutput(gosecScan.FinalOutput)
	return nil
}

func (gosecScan *GosecScan) prepareGosecOutput(gosecOutput GosecOutput) {
	huskyCIgosecResults := types.HuskyCISecurityTestOutput{}
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
	gosecScan.Vulnerabilities = huskyCIgosecResults
}

func (gosecScan *GosecScan) prepareContainerAfterScanGosec(gosecContainer *types.Container) {
	if len(gosecScan.Vulnerabilities.MediumVulns) > 0 || len(gosecScan.Vulnerabilities.HighVulns) > 0 {
		gosecContainer.CInfo = "Issues found."
		gosecContainer.CResult = "failed"
	} else if len(gosecScan.Vulnerabilities.LowVulns) > 0 && (len(gosecScan.Vulnerabilities.MediumVulns) == 0 || len(gosecScan.Vulnerabilities.HighVulns) == 0) {
		gosecContainer.CInfo = "Warnings found."
		gosecContainer.CResult = "passed"
	}
	gosecContainer.CStatus = "finished"
	gosecContainer.CID = gosecScan.CID
	gosecContainer.COutput = gosecScan.RawOutput
	gosecContainer.FinishedAt = time.Now()
}

func newContainerGosec() (types.Container, error) {
	gosecContainer := types.Container{}
	gosecQuery := map[string]interface{}{"name": "gosec"}
	gosecSecurityTest, err := db.FindOneDBSecurityTest(gosecQuery)
	if err != nil {
		return gosecContainer, err
	}
	return types.Container{
		SecurityTest: gosecSecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

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
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// NpmauditScan holds all information needed for a Npmaudit scan.
type NpmauditScan struct {
	RID             string
	CID             string
	URL             string
	Branch          string
	Image           string
	Command         string
	RawOutput       string
	ErrorFound      error
	FinalOutput     NpmAuditOutput
	Vulnerabilities types.HuskyCISecurityTestOutput
}

// NpmAuditOutput is the struct that stores all npm audit output
type NpmAuditOutput struct {
	Advisories    map[string]Vulnerability `json:"advisories"`
	Metadata      Metadata                 `json:"metadata"`
	FailedRunning bool
}

// Vulnerability is the granular output of a security info found
type Vulnerability struct {
	Findings           []Finding `json:"findings"`
	ID                 int       `json:"id"`
	ModuleName         string    `json:"module_name"`
	VulnerableVersions string    `json:"vulnerable_versions"`
	Severity           string    `json:"severity"`
	Overview           string    `json:"overview"`
}

// Finding holds the version of a given security issue found
type Finding struct {
	Version string `json:"version"`
}

// Metadata is the struct that holds vulnerabilities summary
type Metadata struct {
	Vulnerabilities VulnerabilitiesSummary `json:"vulnerabilities"`
}

// VulnerabilitiesSummary is the struct that has all types of possible vulnerabilities from npm audit
type VulnerabilitiesSummary struct {
	Info     int `json:"info"`
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

func newScanNpmaudit(URL, branch, command string) NpmauditScan {
	return NpmauditScan{
		Image:   "huskyci/npmaudit",
		URL:     URL,
		Branch:  branch,
		Command: util.HandleCmd(URL, branch, command),
	}
}

func initNpmaudit(enryScan EnryScan, allScansResult *AllScansResult) error {
	npmauditScan, npmauditContainer, err := runScanNpmaudit(enryScan.URL, enryScan.Branch)
	if err != nil {
		return err
	}

	for _, highVuln := range npmauditScan.Vulnerabilities.HighVulns {
		allScansResult.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns = append(allScansResult.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns, highVuln)
	}
	for _, mediumVuln := range npmauditScan.Vulnerabilities.MediumVulns {
		allScansResult.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns = append(allScansResult.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns, mediumVuln)
	}
	for _, lowVuln := range npmauditScan.Vulnerabilities.LowVulns {
		allScansResult.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns = append(allScansResult.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns, lowVuln)
	}

	allScansResult.FinalResult = npmauditContainer.CResult
	allScansResult.Status = npmauditContainer.CStatus
	allScansResult.Containers = append(allScansResult.Containers, npmauditContainer)
	return nil
}

func runScanNpmaudit(URL, branch string) (NpmauditScan, types.Container, error) {
	npmauditScan := NpmauditScan{}
	npmauditContainer, err := newContainerNpmaudit()
	if err != nil {
		log.Error("runScanNpmaudit", "NPMAUDIT", 1029, err)
		return npmauditScan, npmauditContainer, err
	}
	npmauditScan = newScanNpmaudit(URL, branch, npmauditContainer.SecurityTest.Cmd)
	if err := npmauditScan.startNpmaudit(); err != nil {
		return npmauditScan, npmauditContainer, err
	}

	npmauditScan.prepareContainerAfterScanNpmaudit(&npmauditContainer)
	return npmauditScan, npmauditContainer, nil
}

func (npmauditScan *NpmauditScan) startNpmaudit() error {
	if err := npmauditScan.dockerRunNpmaudit(); err != nil {
		npmauditScan.ErrorFound = err
		return err
	}
	if err := npmauditScan.analyzeNpmaudit(); err != nil {
		npmauditScan.ErrorFound = err
		return err
	}
	return nil
}

func (npmauditScan *NpmauditScan) dockerRunNpmaudit() error {
	CID, cOutput, err := huskydocker.DockerRun(npmauditScan.Image, npmauditScan.Command)
	if err != nil {
		return err
	}
	npmauditScan.CID = CID
	npmauditScan.RawOutput = cOutput
	return nil
}

func (npmauditScan *NpmauditScan) analyzeNpmaudit() error {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(npmauditScan.RawOutput, "ERROR_CLONING")
	failedRunning := strings.Contains(npmauditScan.RawOutput, "ERROR_RUNNING_NPMAUDIT")

	if errorClonning {
		errorMsg := errors.New("error clonning")
		log.Error("analyzeNpmaudit", "NPMAUDIT", 1031, npmauditScan.URL, npmauditScan.Branch, errorMsg)
		return errorMsg
	}

	if failedRunning {
		errorMsg := errors.New("internal error safety - ERROR_RUNNING_SAFETY")
		log.Error("analyzeNpmaudit", "NPMAUDIT", 1034, errorMsg)
		return errorMsg
	}

	// step 2: nil cOutput states that no Issues were found.
	if npmauditScan.RawOutput == "" {
		return nil
	}

	// step 3: Unmarshall rawOutput into finalOutput, that is a GosecOutput struct.
	if err := json.Unmarshal([]byte(npmauditScan.RawOutput), &npmauditScan.FinalOutput); err != nil {
		log.Error("analyzeNpmaudit", "NPMAUDIT", 1014, npmauditScan.RawOutput, err)
		return err
	}

	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	npmauditScan.prepareNpmauditOutput(npmauditScan.FinalOutput)
	return nil
}

func (npmauditScan *NpmauditScan) prepareNpmauditOutput(npmauditOutput NpmAuditOutput) {

	huskyCInpmauditResults := types.HuskyCISecurityTestOutput{}

	if npmauditScan.FinalOutput.FailedRunning {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Severity = "low"
		npmauditVuln.Details = "It looks like your project doesn't have package-lock.json. huskyCI was not able to run npm audit properly."

		huskyCInpmauditResults.LowVulns = append(huskyCInpmauditResults.LowVulns, npmauditVuln)
		return
	}

	for _, issue := range npmauditScan.FinalOutput.Advisories {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Details = issue.Overview
		npmauditVuln.VunerableBelow = issue.VulnerableVersions
		npmauditVuln.Code = issue.ModuleName
		for _, findings := range issue.Findings {
			npmauditVuln.Version = findings.Version
		}

		switch issue.Severity {
		case "info":
			npmauditVuln.Severity = "low"
			huskyCInpmauditResults.LowVulns = append(huskyCInpmauditResults.LowVulns, npmauditVuln)
		case "low":
			npmauditVuln.Severity = "low"
			huskyCInpmauditResults.LowVulns = append(huskyCInpmauditResults.LowVulns, npmauditVuln)
		case "moderate":
			npmauditVuln.Severity = "medium"
			huskyCInpmauditResults.MediumVulns = append(huskyCInpmauditResults.MediumVulns, npmauditVuln)
		case "high":
			npmauditVuln.Severity = "high"
			huskyCInpmauditResults.HighVulns = append(huskyCInpmauditResults.HighVulns, npmauditVuln)
		case "critical":
			npmauditVuln.Severity = "high"
			huskyCInpmauditResults.HighVulns = append(huskyCInpmauditResults.HighVulns, npmauditVuln)
		}

	}

	npmauditScan.Vulnerabilities = huskyCInpmauditResults
}

func (npmauditScan *NpmauditScan) prepareContainerAfterScanNpmaudit(npmauditContainer *types.Container) {
	if len(npmauditScan.Vulnerabilities.MediumVulns) > 0 || len(npmauditScan.Vulnerabilities.HighVulns) > 0 {
		npmauditContainer.CInfo = "Issues found."
		npmauditContainer.CResult = "failed"
	} else if len(npmauditScan.Vulnerabilities.LowVulns) > 0 && (len(npmauditScan.Vulnerabilities.MediumVulns) == 0 || len(npmauditScan.Vulnerabilities.HighVulns) == 0) {
		npmauditContainer.CInfo = "Warnings found."
		npmauditContainer.CResult = "passed"
	}
	npmauditContainer.CStatus = "finished"
	npmauditContainer.CID = npmauditScan.CID
	npmauditContainer.COutput = npmauditScan.RawOutput
	npmauditContainer.FinishedAt = time.Now()
}

func newContainerNpmaudit() (types.Container, error) {
	npmauditContainer := types.Container{}
	npmauditQuery := map[string]interface{}{"name": "npmaudit"}
	npmauditSecurityTest, err := db.FindOneDBSecurityTest(npmauditQuery)
	if err != nil {
		log.Error("newContainerNpmaudit", "NPMAUDIT", 2012, err)
		return npmauditContainer, err
	}
	return types.Container{
		SecurityTest: npmauditSecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

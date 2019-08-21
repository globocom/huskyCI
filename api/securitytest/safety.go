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

// SafetyScan holds all information needed for a Safety scan.
type SafetyScan struct {
	RID             string
	CID             string
	URL             string
	Branch          string
	Image           string
	Command         string
	RawOutput       string
	FinalOutput     SafetyOutput
	Vulnerabilities types.HuskyCISecurityTestOutput
}

// SafetyOutput is the struct that holds issues, messages and errors found on a Safety scan.
type SafetyOutput struct {
	SafetyIssues   []SafetyIssue `json:"issues"`
	ReqNotFound    bool
	WarningFound   bool
	OutputWarnings []string
}

// SafetyIssue is a struct that holds the results that were scanned and the file they came from.
type SafetyIssue struct {
	Dependency string `json:"dependency"`
	Below      string `json:"vulnerable_below"`
	Version    string `json:"installed_version"`
	Comment    string `json:"description"`
	ID         string `json:"id"`
}

func newScanSafety(URL, branch, command string) SafetyScan {
	return SafetyScan{
		Image:   "huskyci/safety",
		URL:     URL,
		Branch:  branch,
		Command: util.HandleCmd(URL, branch, "", command),
	}
}

func initSafety(enryScan EnryScan, allScansResult *AllScansResult) error {
	safetyScan, safetyContainer, err := runScanSafety(enryScan.URL, enryScan.Branch)
	if err != nil {
		return err
	}

	for _, highVuln := range safetyScan.Vulnerabilities.HighVulns {
		allScansResult.HuskyCIResults.PythonResults.HuskyCISafetyOutput.HighVulns = append(allScansResult.HuskyCIResults.PythonResults.HuskyCISafetyOutput.HighVulns, highVuln)
	}
	for _, mediumVuln := range safetyScan.Vulnerabilities.MediumVulns {
		allScansResult.HuskyCIResults.PythonResults.HuskyCISafetyOutput.MediumVulns = append(allScansResult.HuskyCIResults.PythonResults.HuskyCISafetyOutput.MediumVulns, mediumVuln)
	}
	for _, lowVuln := range safetyScan.Vulnerabilities.LowVulns {
		allScansResult.HuskyCIResults.PythonResults.HuskyCISafetyOutput.LowVulns = append(allScansResult.HuskyCIResults.PythonResults.HuskyCISafetyOutput.LowVulns, lowVuln)
	}

	allScansResult.FinalResult = safetyContainer.CResult
	allScansResult.Status = safetyContainer.CStatus
	allScansResult.Containers = append(allScansResult.Containers, safetyContainer)
	return nil
}

func runScanSafety(URL, branch string) (SafetyScan, types.Container, error) {
	safetyScan := SafetyScan{}
	safetyContainer, err := newContainerSafety()
	if err != nil {
		return safetyScan, safetyContainer, err
	}
	safetyScan = newScanSafety(URL, branch, safetyContainer.SecurityTest.Cmd)
	if err := safetyScan.startSafety(); err != nil {
		return safetyScan, safetyContainer, err
	}

	safetyScan.prepareContainerAfterScanSafety(&safetyContainer)
	return safetyScan, safetyContainer, nil
}

func (safetyScan *SafetyScan) startSafety() error {
	if err := safetyScan.dockerRunSafety(); err != nil {
		return err
	}
	if err := safetyScan.analyzeSafety(); err != nil {
		return err
	}
	// log.Info("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
	return nil
}

func (safetyScan *SafetyScan) dockerRunSafety() error {
	CID, cOutput, err := huskydocker.DockerRun(safetyScan.Image, safetyScan.Command)
	if err != nil {
		// log.Error("DockerRun", "DOCKERRUN", 3013, err)
		return err
	}
	safetyScan.CID = CID
	safetyScan.RawOutput = cOutput
	return nil
}

func (safetyScan *SafetyScan) analyzeSafety() error {

	errorCloning := strings.Contains(safetyScan.RawOutput, "ERROR_CLONING")
	failedRunning := strings.Contains(safetyScan.RawOutput, "ERROR_RUNNING_SAFETY")
	reqNotFound := strings.Contains(safetyScan.RawOutput, "ERROR_REQ_NOT_FOUND")
	warningFound := strings.Contains(safetyScan.RawOutput, "Warning: unpinned requirement ")

	if errorCloning {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return errors.New("error clonning")
	}

	if failedRunning {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return errors.New("failed running")
	}

	if reqNotFound {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return errors.New("failed running")
	}

	if warningFound {
		outputJSON := util.GetLastLine(safetyScan.RawOutput)
		safetyScan.FinalOutput.OutputWarnings = util.GetAllLinesButLast(safetyScan.RawOutput)
		safetyScan.RawOutput = outputJSON
	}

	cOutputSanitized := util.SanitizeSafetyJSON(safetyScan.RawOutput)
	safetyScan.RawOutput = cOutputSanitized

	// step 3: Unmarshall rawOutput into finalOutput, that is a GosecOutput struct.
	if err := json.Unmarshal([]byte(safetyScan.RawOutput), &safetyScan.FinalOutput); err != nil {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return err
	}
	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	safetyScan.prepareSafetyOutput(safetyScan.FinalOutput)
	return nil
}

func (safetyScan *SafetyScan) prepareSafetyOutput(safetyOutput SafetyOutput) {

	var huskyCIsafetyResults types.HuskyCISecurityTestOutput
	var onlyWarning bool

	if safetyOutput.ReqNotFound {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
		safetyVuln.Details = "requirements.txt not found"

		huskyCIsafetyResults.LowVulns = append(huskyCIsafetyResults.LowVulns, safetyVuln)

		safetyScan.Vulnerabilities = huskyCIsafetyResults
		return
	}

	if safetyOutput.WarningFound {

		if len(safetyOutput.SafetyIssues) == 0 {
			onlyWarning = true
		}

		for _, warning := range safetyOutput.OutputWarnings {
			safetyVuln := types.HuskyCIVulnerability{}
			safetyVuln.Language = "Python"
			safetyVuln.SecurityTool = "Safety"
			safetyVuln.Severity = "warning"
			safetyVuln.Details = util.AdjustWarningMessage(warning)

			huskyCIsafetyResults.LowVulns = append(huskyCIsafetyResults.LowVulns, safetyVuln)

		}
		if onlyWarning {
			safetyScan.Vulnerabilities = huskyCIsafetyResults
			return
		}
	}

	for _, issue := range safetyOutput.SafetyIssues {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "high"
		safetyVuln.Details = issue.Comment
		safetyVuln.Code = issue.Dependency + " " + issue.Version
		safetyVuln.VunerableBelow = issue.Below

		huskyCIsafetyResults.HighVulns = append(huskyCIsafetyResults.HighVulns, safetyVuln)
	}

	safetyScan.Vulnerabilities = huskyCIsafetyResults
}

func (safetyScan *SafetyScan) prepareContainerAfterScanSafety(safetyContainer *types.Container) {
	if len(safetyScan.Vulnerabilities.MediumVulns) > 0 || len(safetyScan.Vulnerabilities.HighVulns) > 0 {
		safetyContainer.CInfo = "Issues found."
		safetyContainer.CResult = "failed"
	} else if len(safetyScan.Vulnerabilities.LowVulns) > 0 && (len(safetyScan.Vulnerabilities.MediumVulns) == 0 || len(safetyScan.Vulnerabilities.HighVulns) == 0) {
		safetyContainer.CInfo = "Warnings found."
		safetyContainer.CResult = "passed"
	}
	safetyContainer.CStatus = "finished"
	safetyContainer.CID = safetyScan.CID
	safetyContainer.COutput = safetyScan.RawOutput
	safetyContainer.FinishedAt = time.Now()
}

func newContainerSafety() (types.Container, error) {
	safetyContainer := types.Container{}
	safetyQuery := map[string]interface{}{"name": "safety"}
	safetySecurityTest, err := db.FindOneDBSecurityTest(safetyQuery)
	if err != nil {
		return safetyContainer, err
	}
	return types.Container{
		SecurityTest: safetySecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

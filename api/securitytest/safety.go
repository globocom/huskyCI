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

// SafetyScan holds all information needed for a Safety scan.
type SafetyScan struct {
	RID             string
	CID             string
	URL             string
	Branch          string
	Image           string
	Command         string
	RawOutput       string
	ErrorFound      error
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
		Command: util.HandleCmd(URL, branch, command),
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
		log.Error("runScanSafety", "SAFETY", 1029, err)
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
		safetyScan.ErrorFound = err
		return err
	}
	if err := safetyScan.analyzeSafety(); err != nil {
		safetyScan.ErrorFound = err
		return err
	}
	return nil
}

func (safetyScan *SafetyScan) dockerRunSafety() error {
	CID, cOutput, err := huskydocker.DockerRun(safetyScan.Image, safetyScan.Command)
	if err != nil {
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

	// step 1: check for any errors when clonning repo
	if errorCloning {
		errorMsg := errors.New("error clonning")
		log.Error("analyzeSafety", "SAFETY", 1031, safetyScan.URL, safetyScan.Branch, errorMsg)
		return errorMsg
	}

	// step 2: check if there were any internal errors running safety
	if failedRunning {
		errorMsg := errors.New("internal error safety - ERROR_RUNNING_SAFETY")
		log.Error("analyzeSafety", "SAFETY", 1033, errorMsg)
		return errorMsg
	}

	// step 3: check if requirements.txt were found or not
	if reqNotFound {
		safetyScan.FinalOutput.ReqNotFound = true
	}

	// step 4: check if warning were found and handle its output
	if warningFound {
		safetyScan.FinalOutput.WarningFound = true
		outputJSON := util.GetLastLine(safetyScan.RawOutput)
		safetyScan.FinalOutput.OutputWarnings = util.GetAllLinesButLast(safetyScan.RawOutput)
		safetyScan.RawOutput = outputJSON
	}

	cOutputSanitized := util.SanitizeSafetyJSON(safetyScan.RawOutput)
	safetyScan.RawOutput = cOutputSanitized

	// step 5: Unmarshall rawOutput into finalOutput, that is a Safety struct.
	if err := json.Unmarshal([]byte(safetyScan.RawOutput), &safetyScan.FinalOutput); err != nil {
		log.Error("analyzeSafety", "SAFETY", 1018, safetyScan.RawOutput, err)
		return err
	}
	// step 6: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
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
		log.Error("newContainerSafety", "SAFETY", 2012, err)
		return safetyContainer, err
	}
	return types.Container{
		SecurityTest: safetySecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

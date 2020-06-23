// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

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

func analyzeSafety(safetyScan *SecTestScanInfo) error {

	failedRunning := strings.Contains(safetyScan.Container.COutput, "ERROR_RUNNING_SAFETY")
	reqNotFound := strings.Contains(safetyScan.Container.COutput, "ERROR_REQ_NOT_FOUND")
	warningFound := strings.Contains(safetyScan.Container.COutput, "Warning: unpinned requirement ")

	safetyOutput := SafetyOutput{}
	safetyScan.FinalOutput = safetyOutput

	// check if there were any internal errors running safety
	if failedRunning {
		errorMsg := errors.New("internal error safety - ERROR_RUNNING_SAFETY")
		log.Error("analyzeSafety", "SAFETY", 1033, errorMsg)
		safetyScan.ErrorFound = errorMsg
		safetyScan.prepareContainerAfterScan()
		return errorMsg
	}

	// check if requirements.txt were found or not
	if reqNotFound {
		safetyScan.ReqNotFound = true
		safetyScan.prepareSafetyVulns()
		safetyScan.prepareContainerAfterScan()
		return nil
	}

	// check if warning were found and handle its output
	if warningFound {
		safetyScan.WarningFound = true
		outputJSON := util.GetLastLine(safetyScan.Container.COutput)
		safetyOutput.OutputWarnings = util.GetAllLinesButLast(safetyScan.Container.COutput)
		safetyScan.Container.COutput = outputJSON
	}

	cOutputSanitized := util.SanitizeSafetyJSON(safetyScan.Container.COutput)
	safetyScan.Container.COutput = cOutputSanitized

	// Unmarshall rawOutput into finalOutput, that is a Safety struct.
	if err := json.Unmarshal([]byte(safetyScan.Container.COutput), &safetyOutput); err != nil {
		log.Error("analyzeSafety", "SAFETY", 1018, safetyScan.Container.COutput, err)
		safetyScan.ErrorFound = err
		safetyScan.prepareContainerAfterScan()
		return err
	}
	safetyScan.FinalOutput = safetyOutput

	// check results and prepare all vulnerabilities found
	safetyScan.prepareSafetyVulns()
	safetyScan.prepareContainerAfterScan()
	return nil
}

func (safetyScan *SecTestScanInfo) prepareSafetyVulns() {

	huskyCIsafetyResults := types.HuskyCISecurityTestOutput{}
	safetyOutput := safetyScan.FinalOutput.(SafetyOutput)
	onlyWarning := false

	if safetyScan.ReqNotFound {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "low"
		safetyVuln.Title = "No requirements.txt found."
		safetyVuln.Details = "It looks like your project doesn't have a requirements.txt file. huskyCI was not able to run safety properly."

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
			safetyVuln.Severity = "low"
			safetyVuln.Title = "Safety scan warning."
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
		safetyVuln.Title = fmt.Sprintf("Vulnerable Dependency: %s (%s)", issue.Dependency, issue.Below)
		safetyVuln.VunerableBelow = issue.Below

		huskyCIsafetyResults.HighVulns = append(huskyCIsafetyResults.HighVulns, safetyVuln)
	}

	safetyScan.Vulnerabilities = huskyCIsafetyResults
}

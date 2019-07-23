// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

//SafetyOutput is the struct that holds issues, messages and errors found on a Safety scan.
type SafetyOutput struct {
	SafetyIssues   []SafetyIssue `json:"issues"`
	ReqNotFound    bool
	WarningFound   bool
	OutputWarnings []string
}

//SafetyIssue is a struct that holds the results that were scanned and the file they came from.
type SafetyIssue struct {
	Dependency string `json:"dependency"`
	Below      string `json:"vulnerable_below"`
	Version    string `json:"installed_version"`
	Comment    string `json:"description"`
	ID         string `json:"id"`
}

//SafetyCheckOutputFlow analyses the output from Safety and sets cResult based on it.
func SafetyCheckOutputFlow(CID string, cOutput string, RID string) {

	reqNotFound := strings.Contains(cOutput, "ERROR_REQ_NOT_FOUND")
	failedRunning := strings.Contains(cOutput, "ERROR_RUNNING_SAFETY")
	warningFound := strings.Contains(cOutput, "Warning: unpinned requirement ")
	errorCloning := strings.Contains(cOutput, "ERROR_CLONING")

	// step 1: check for any errors when clonning repo
	if errorCloning {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: check for any errors when running securityTest
	if failedRunning {
		if err := updateInfoAndResultBasedOnCID("Internal error running Safety.", "error", CID); err != nil {
			return
		}
		return
	}

	// step 3: check if requirements.txt were found or not
	if reqNotFound {
		if err := updateInfoAndResultBasedOnCID("Requirements not found.", "warning", CID); err != nil {
			return
		}

		safetyOutput := SafetyOutput{ReqNotFound: true}
		if err := updateHuskyCIResultsBasedOnRID(RID, "safety", safetyOutput); err != nil {
			return
		}
		return
	}

	// step 4: check if warning were found and handle its output
	safetyOutput := SafetyOutput{}
	if warningFound {
		outputJSON := util.GetLastLine(cOutput)
		safetyOutput.OutputWarnings = util.GetAllLinesButLast(cOutput)
		cOutput = outputJSON
	}
	cOutputSanitized := util.SanitizeSafetyJSON(cOutput)

	// step 5: unmarshall safety output
	err := json.Unmarshal([]byte(cOutputSanitized), &safetyOutput)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 1018, cOutput, err)
		return
	}

	// step 6: check if issues, warnings, or both were found
	if len(safetyOutput.SafetyIssues) == 0 {

		// no issues but warning found!
		if warningFound {
			if err := updateInfoAndResultBasedOnCID("Warnings found.", "warning", CID); err != nil {
				return
			}

			safetyOutput := SafetyOutput{WarningFound: true}
			if err := updateHuskyCIResultsBasedOnRID(RID, "safety", safetyOutput); err != nil {
				return
			}

			return
		}

		// no issues and no warning
		if err := updateInfoAndResultBasedOnCID("No issues found.", "passed", CID); err != nil {
			return
		}

		return
	}

	// Issues found.
	if err := updateInfoAndResultBasedOnCID("Issues found.", "failed", CID); err != nil {
		return
	}

	// step 6: finally, update analysis with huskyCI results
	if err := updateHuskyCIResultsBasedOnRID(RID, "safety", safetyOutput); err != nil {
		return
	}

}

// prepareHuskyCISafetyResults will prepare Safety output to be added into PythonResults struct
func prepareHuskyCISafetyResults(safetyOutput SafetyOutput) types.HuskyCISafetyOutput {

	var huskyCIsafetyResults types.HuskyCISafetyOutput
	var onlyWarning bool

	if safetyOutput.ReqNotFound {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
		safetyVuln.Details = "requirements.txt not found"

		huskyCIsafetyResults.LowVulnsSafety = append(huskyCIsafetyResults.LowVulnsSafety, safetyVuln)

		return huskyCIsafetyResults
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

			huskyCIsafetyResults.LowVulnsSafety = append(huskyCIsafetyResults.LowVulnsSafety, safetyVuln)

		}
		if onlyWarning {
			return huskyCIsafetyResults
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

		huskyCIsafetyResults.HighVulnsSafety = append(huskyCIsafetyResults.HighVulnsSafety, safetyVuln)
	}

	return huskyCIsafetyResults
}

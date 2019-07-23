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

//RetirejsOutput is the struct that holds issues, messages and errors found on a Retire scan.
type RetirejsOutput struct {
	RetirejsResult []RetirejsResult `json:"results"`
}

//RetirejsResult is a struct that holds the scanned results.
type RetirejsResult struct {
	Component       string                    `json:"component"`
	Version         string                    `json:"version"`
	Level           int                       `json:"level"`
	Vulnerabilities []RetireJSVulnerabilities `json:"vulnerabilities"`
}

//RetireJSVulnerabilities is a struct that holds the vulnerabilities found on a scan.
type RetireJSVulnerabilities struct {
	Info        []string                         `json:"info"`
	Severity    string                           `json:"severity"`
	Identifiers RetireJSVulnerabilityIdentifiers `json:"identifiers"`
}

//RetireJSVulnerabilityIdentifiers is a struct that holds identifiying information on a vulnerability found.
type RetireJSVulnerabilityIdentifiers struct {
	Summary string
}

//RetirejsCheckOutputFlow analyses the output from RetireJS and sets cResult basdes on it.
func RetirejsCheckOutputFlow(CID string, cOutput string, RID string) {

	errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	failedRunning := strings.Contains(cOutput, "ERROR_RUNNING_RETIREJS")

	// step 1: check for any errors when clonning repo
	if errorClonning {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: check for any errors when running securityTest
	if failedRunning {
		if err := updateInfoAndResultBasedOnCID("Internal error running RetireJS.", "error", CID); err != nil {
			return
		}

		retireJSOutput := []RetirejsOutput{}
		if err := updateHuskyCIResultsBasedOnRID(RID, "retirejs", retireJSOutput); err != nil {
			return
		}

		return
	}

	// step 3: get retireJS output to be checked
	retirejsOutput := []RetirejsOutput{}
	err := json.Unmarshal([]byte(cOutput), &retirejsOutput)
	if err != nil {
		log.Error("RetirejsStartAnalysis", "RETIREJS", 1014, cOutput, err)
		return
	}

	// step 4: sets the container output to "No issues found" if RetirejsIssues returns an empty slice
	if len(retirejsOutput) == 0 {
		if err := updateInfoAndResultBasedOnCID("No issues found.", "passed", CID); err != nil {
			return
		}
		return
	}

	// step 5: find Vulnerabilities that have severity "medium" or "high"
	cResult := "passed"
	issueMessage := "No issues found."
	for _, output := range retirejsOutput {
		for _, result := range output.RetirejsResult {
			for _, vulnerability := range result.Vulnerabilities {
				if vulnerability.Severity == "high" || vulnerability.Severity == "medium" {
					cResult = "failed"
					issueMessage = "Issues found."
					break
				}
			}
		}
	}
	if err := updateInfoAndResultBasedOnCID(issueMessage, cResult, CID); err != nil {
		return
	}

	// step 6: finally, update analysis with huskyCI results
	if err := updateHuskyCIResultsBasedOnRID(RID, "retirejs", retirejsOutput); err != nil {
		return
	}

}

// prepareHuskyCIRetirejsOutput will prepare Retirejs output to be added into JavaScriptResults struct
func prepareHuskyCIRetirejsOutput(retirejsOutput []RetirejsOutput) types.HuskyCISecurityTestOutput {

	var huskyCIretireJSResults types.HuskyCISecurityTestOutput
	var huskyCIretireJSResultsFinal types.HuskyCISecurityTestOutput

	// failedRunning
	if retirejsOutput == nil {
		retirejsVuln := types.HuskyCIVulnerability{}
		retirejsVuln.Language = "JavaScript"
		retirejsVuln.SecurityTool = "RetireJS"
		retirejsVuln.Severity = "low"
		retirejsVuln.Details = "It looks like your project doesn't have package.json or yarn.lock. huskyCI was not able to run RetireJS properly."

		huskyCIretireJSResults.LowVulns = append(huskyCIretireJSResults.LowVulns, retirejsVuln)

		return huskyCIretireJSResults
	}

	for _, output := range retirejsOutput {
		for _, result := range output.RetirejsResult {
			for _, vulnerability := range result.Vulnerabilities {
				retirejsVuln := types.HuskyCIVulnerability{}
				retirejsVuln.Language = "JavaScript"
				retirejsVuln.SecurityTool = "RetireJS"
				retirejsVuln.Severity = vulnerability.Severity
				retirejsVuln.Code = result.Component
				retirejsVuln.Version = result.Version
				for _, info := range vulnerability.Info {
					retirejsVuln.Details = retirejsVuln.Details + info + "\n"
				}
				retirejsVuln.Details = retirejsVuln.Details + vulnerability.Identifiers.Summary

				switch retirejsVuln.Severity {
				case "low":
					huskyCIretireJSResults.LowVulns = append(huskyCIretireJSResults.LowVulns, retirejsVuln)
				case "medium":
					huskyCIretireJSResults.MediumVulns = append(huskyCIretireJSResults.MediumVulns, retirejsVuln)
				case "high":
					huskyCIretireJSResults.HighVulns = append(huskyCIretireJSResults.HighVulns, retirejsVuln)
				}
			}
		}
	}

	huskyCIretireJSResultsFinal.LowVulns = util.CountRetireJSOccurrences(huskyCIretireJSResults.LowVulns)
	huskyCIretireJSResultsFinal.MediumVulns = util.CountRetireJSOccurrences(huskyCIretireJSResults.MediumVulns)
	huskyCIretireJSResultsFinal.HighVulns = util.CountRetireJSOccurrences(huskyCIretireJSResults.HighVulns)

	return huskyCIretireJSResultsFinal
}

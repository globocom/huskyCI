// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"github.com/globocom/huskyCI/api/types"
)

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

// NpmAuditCheckOutputFlow analyses the output from Npm Audit and sets a cResult based on it.
func NpmAuditCheckOutputFlow(CID string, cOutput string, RID string) {

	// errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	// failedRunning := strings.Contains(cOutput, "ERROR_RUNNING_NPMAUDIT")

	// // step 1: check for any errors when clonning repo
	// if errorClonning {
	// 	if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
	// 		return
	// 	}
	// 	return
	// }

	// // step 2: check for any errors when running securityTest
	// if failedRunning {
	// 	if err := updateInfoAndResultBasedOnCID("Internal error running NPM Audit.", "error", CID); err != nil {
	// 		return
	// 	}

	// 	npmAuditOutput := NpmAuditOutput{FailedRunning: failedRunning}
	// 	if err := updateHuskyCIResultsBasedOnRID(RID, "npmaudit", npmAuditOutput); err != nil {
	// 		return
	// 	}

	// 	return
	// }

	// // step 3: Unmarshall cOutput into NpmAuditOutput struct.
	// npmAuditOutput := NpmAuditOutput{}
	// err := json.Unmarshal([]byte(cOutput), &npmAuditOutput)
	// if err != nil {
	// 	log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 1022, cOutput, err)
	// 	return
	// }

	// // step 4: find Issues that have severity "moderate" or "high.
	// cResult := "passed"
	// issueMessage := "No issues found."
	// for _, vulnerability := range npmAuditOutput.Advisories {
	// 	if vulnerability.Severity == "high" || vulnerability.Severity == "moderate" {
	// 		cResult = "failed"
	// 		issueMessage = "Issues found."
	// 		break
	// 	}
	// }
	// if err := updateInfoAndResultBasedOnCID(issueMessage, cResult, CID); err != nil {
	// 	return
	// }

	// // step 6: finally, update analysis with huskyCI results
	// if err := updateHuskyCIResultsBasedOnRID(RID, "npmaudit", npmAuditOutput); err != nil {
	// 	return
	// }

}

// prepareHuskyCINpmAuditResults will prepare NpmAudit output to be added into JavaScriptResults struct
func prepareHuskyCINpmAuditResults(npmAuditOutput NpmAuditOutput) types.HuskyCISecurityTestOutput {

	var huskyCInpmAuditResults types.HuskyCISecurityTestOutput

	if npmAuditOutput.FailedRunning {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Severity = "low"
		npmauditVuln.Details = "It looks like your project doesn't have package-lock.json. huskyCI was not able to run npm audit properly."

		huskyCInpmAuditResults.LowVulns = append(huskyCInpmAuditResults.LowVulns, npmauditVuln)

		return huskyCInpmAuditResults
	}

	for _, issue := range npmAuditOutput.Advisories {
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
			huskyCInpmAuditResults.LowVulns = append(huskyCInpmAuditResults.LowVulns, npmauditVuln)
		case "low":
			npmauditVuln.Severity = "low"
			huskyCInpmAuditResults.LowVulns = append(huskyCInpmAuditResults.LowVulns, npmauditVuln)
		case "moderate":
			npmauditVuln.Severity = "medium"
			huskyCInpmAuditResults.MediumVulns = append(huskyCInpmAuditResults.MediumVulns, npmauditVuln)
		case "high":
			npmauditVuln.Severity = "high"
			huskyCInpmAuditResults.HighVulns = append(huskyCInpmAuditResults.HighVulns, npmauditVuln)
		case "critical":
			npmauditVuln.Severity = "high"
			huskyCInpmAuditResults.HighVulns = append(huskyCInpmAuditResults.HighVulns, npmauditVuln)
		}

	}

	return huskyCInpmAuditResults
}

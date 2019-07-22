// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
)

// NpmAuditOutput is the struct that stores all npm audit output
type NpmAuditOutput struct {
	Advisories map[string]Vulnerability `json:"advisories"`
	Metadata   Metadata                 `json:"metadata"`
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

// NpmAuditStartAnalysis analyses the output from Npm Audit and sets a cResult based on it.
func NpmAuditStartAnalysis(CID string, cOutput string, RID string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.1: error cloning repository!
	if strings.Contains(cOutput, "ERROR_CLONING") {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   "Error clonning repository.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, "Step 0.1 ", err)
		}
		return
	}

	// step 0.2: repository doesn't have package.json
	failedRunning := strings.Contains(cOutput, "ERROR_RUNNING_NPMAUDIT")
	if failedRunning {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   "Internal error running NPM Audit.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, "Step 0.2 ", err)
			return
		}

		// get updated analysis based on its RID
		analysisQuery = map[string]interface{}{"RID": RID}
		analysis, err := db.FindOneDBAnalysis(analysisQuery)
		if err != nil {
			log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2008, CID, err)
			return
		}

		// finally, update analysis with huskyCI results
		npmAuditOutput := NpmAuditOutput{}
		analysis.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput = prepareHuskyCINpmAuditResults(npmAuditOutput, failedRunning)
		err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
		if err != nil {
			log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, err)
			return
		}

		return
	}

	// step 1: Unmarshall cOutput into NpmAuditOutput struct.
	npmAuditOutput := NpmAuditOutput{}
	err := json.Unmarshal([]byte(cOutput), &npmAuditOutput)
	if err != nil {
		log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 1022, cOutput, err)
		return
	}

	// step 2: find Issues that have severity "moderate" or "high.
	cResult = "passed"
	for _, vulnerability := range npmAuditOutput.Advisories {
		if vulnerability.Severity == "high" || vulnerability.Severity == "moderate" {
			cResult = "failed"
			break
		}
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	issueMessage := "No issues found."
	if cResult == "failed" {
		issueMessage = "Issues found."
	}
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
			"containers.$.cInfo":   issueMessage,
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, "Step 3 ", err)
		return
	}

	// step 4: get updated analysis based on its RID
	analysisQuery = map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2008, CID, err)
		return
	}

	// step 5: finally, update analysis with huskyCI results
	analysis.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput = prepareHuskyCINpmAuditResults(npmAuditOutput, failedRunning)
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, err)
		return
	}

}

// prepareHuskyCINpmAuditResults will prepare NpmAudit output to be added into JavaScriptResults struct
func prepareHuskyCINpmAuditResults(npmAuditOutput NpmAuditOutput, failedRunning bool) types.HuskyCINpmAuditOutput {

	var huskyCInpmAuditResults types.HuskyCINpmAuditOutput

	if failedRunning {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Severity = "low"
		npmauditVuln.Details = "It looks like your project doesn't have package-lock.json. huskyCI was not able to run npm audit properly."

		huskyCInpmAuditResults.LowVulnsNpmAudit = append(huskyCInpmAuditResults.LowVulnsNpmAudit, npmauditVuln)

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
			huskyCInpmAuditResults.LowVulnsNpmAudit = append(huskyCInpmAuditResults.LowVulnsNpmAudit, npmauditVuln)
		case "low":
			npmauditVuln.Severity = "low"
			huskyCInpmAuditResults.LowVulnsNpmAudit = append(huskyCInpmAuditResults.LowVulnsNpmAudit, npmauditVuln)
		case "moderate":
			npmauditVuln.Severity = "medium"
			huskyCInpmAuditResults.MediumVulnsNpmAudit = append(huskyCInpmAuditResults.MediumVulnsNpmAudit, npmauditVuln)
		case "high":
			npmauditVuln.Severity = "high"
			huskyCInpmAuditResults.HighVulnsNpmAudit = append(huskyCInpmAuditResults.HighVulnsNpmAudit, npmauditVuln)
		case "critical":
			npmauditVuln.Severity = "high"
			huskyCInpmAuditResults.HighVulnsNpmAudit = append(huskyCInpmAuditResults.HighVulnsNpmAudit, npmauditVuln)
		}

	}

	return huskyCInpmAuditResults
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// YarnAuditOutput is the struct that stores all yarn audit output
type YarnAuditOutput struct {
	Advisories       []YarnIssue `json:"advisories"`
	Metadata         Metadata    `json:"metadata"`
	YarnLockNotFound bool
	YarnErrorRunning bool
}

// YarnIssue is the granular output of a security info about yarn found
type YarnIssue struct {
	Findings           []YarnFinding `json:"findings"`
	ID                 int           `json:"id"`
	ModuleName         string        `json:"module_name"`
	VulnerableVersions string        `json:"vulnerable_versions"`
	Severity           string        `json:"severity"`
	Overview           string        `json:"overview"`
	Title              string        `json:"title"`
}

// YarnFinding holds the version of a given yarn security issue found
type YarnFinding struct {
	Version string `json:"version"`
}

// YarnMetadata is the struct that holds vulnerabilities summary
type YarnMetadata struct {
	Vulnerabilities YarnVulnerabilitiesSummary `json:"vulnerabilities"`
}

// YarnVulnerabilitiesSummary is the struct that has all types of possible vulnerabilities from yarn audit
type YarnVulnerabilitiesSummary struct {
	Info     int `json:"info"`
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

func analyzeYarnaudit(yarnAuditScan *SecTestScanInfo) error {

	yarnAuditOutput := YarnAuditOutput{}
	yarnAuditScan.FinalOutput = yarnAuditOutput

	// if package-lock was not found, a warning will be genrated as a low vuln
	yarnLockNotFound := strings.Contains(yarnAuditScan.Container.COutput, "ERROR_YARN_LOCK_NOT_FOUND")
	if yarnLockNotFound {
		yarnAuditScan.YarnLockNotFound = true
		yarnAuditScan.prepareYarnAuditVulns()
		yarnAuditScan.prepareContainerAfterScan()
		return nil
	}

	// if yarn audit fails to run, a warning will be genrated as a low vuln
	YarnErrorRunning := strings.Contains(yarnAuditScan.Container.COutput, "ERROR_RUNNING_YARN_AUDIT")
	if YarnErrorRunning {
		yarnAuditScan.YarnErrorRunning = true
		yarnAuditScan.prepareYarnAuditVulns()
		yarnAuditScan.prepareContainerAfterScan()
		return nil
	}

	// nil cOutput states that no Issues were found.
	if yarnAuditScan.Container.COutput == "" {
		yarnAuditScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a YarnAuditOutput struct.
	if err := json.Unmarshal([]byte(yarnAuditScan.Container.COutput), &yarnAuditOutput); err != nil {
		log.Error("analyzeYarnaudit", "YARNAUDIT", 1036, yarnAuditScan.Container.COutput, err)
		return err
	}
	yarnAuditScan.FinalOutput = yarnAuditOutput

	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	yarnAuditScan.prepareYarnAuditVulns()
	yarnAuditScan.prepareContainerAfterScan()
	return nil
}

func (yarnAuditScan *SecTestScanInfo) prepareYarnAuditVulns() {

	huskyCIyarnauditResults := types.HuskyCISecurityTestOutput{}
	yarnAuditOutput := yarnAuditScan.FinalOutput.(YarnAuditOutput)

	if yarnAuditScan.YarnLockNotFound {
		yarnauditVuln := types.HuskyCIVulnerability{}
		yarnauditVuln.Language = "JavaScript"
		yarnauditVuln.SecurityTool = "YarnAudit"
		yarnauditVuln.Severity = "low"
		yarnauditVuln.Title = "No yarn.lock found."
		yarnauditVuln.Details = "It looks like your project doesn't have a yarn.lock file. If you use Yarn to handle your dependencies, it would be a good idea to commit it so huskyCI can check for vulnerabilities."

		yarnAuditScan.Vulnerabilities.LowVulns = append(yarnAuditScan.Vulnerabilities.LowVulns, yarnauditVuln)
		return
	}

	if yarnAuditScan.YarnErrorRunning {
		yarnauditVuln := types.HuskyCIVulnerability{}
		yarnauditVuln.Language = "JavaScript"
		yarnauditVuln.SecurityTool = "YarnAudit"
		yarnauditVuln.Severity = "low"
		yarnauditVuln.Title = "Error while running yarn audit scan."
		yarnauditVuln.Details = "Yarn returned an error"

		yarnAuditScan.Vulnerabilities.LowVulns = append(yarnAuditScan.Vulnerabilities.LowVulns, yarnauditVuln)
		return
	}

	for _, issue := range yarnAuditOutput.Advisories {
		yarnauditVuln := types.HuskyCIVulnerability{}
		yarnauditVuln.Language = "JavaScript"
		yarnauditVuln.SecurityTool = "YarnAudit"
		yarnauditVuln.Details = issue.Overview
		yarnauditVuln.Title = fmt.Sprintf("Vulnerable Dependency: %s %s (%s)", issue.ModuleName, issue.VulnerableVersions, issue.Title)
		yarnauditVuln.VunerableBelow = issue.VulnerableVersions
		yarnauditVuln.Code = issue.ModuleName
		yarnauditVuln.Occurrences = 1
		for _, findings := range issue.Findings {
			yarnauditVuln.Version = findings.Version
		}

		switch issue.Severity {
		case "info", "low":
			yarnauditVuln.Severity = "low"
			if !vulnListContains(huskyCIyarnauditResults.LowVulns, yarnauditVuln) {
				huskyCIyarnauditResults.LowVulns = append(huskyCIyarnauditResults.LowVulns, yarnauditVuln)
			}
		case "moderate":
			yarnauditVuln.Severity = "medium"
			if !vulnListContains(huskyCIyarnauditResults.MediumVulns, yarnauditVuln) {
				huskyCIyarnauditResults.MediumVulns = append(huskyCIyarnauditResults.MediumVulns, yarnauditVuln)
			}
		case "high", "critical":
			yarnauditVuln.Severity = "high"
			if !vulnListContains(huskyCIyarnauditResults.HighVulns, yarnauditVuln) {
				huskyCIyarnauditResults.HighVulns = append(huskyCIyarnauditResults.HighVulns, yarnauditVuln)
			}
		}

	}

	yarnAuditScan.Vulnerabilities = huskyCIyarnauditResults
}

// vulnListContains increments the occurrence counter in case a vulnerability is found again
func vulnListContains(vulnList []types.HuskyCIVulnerability, vuln types.HuskyCIVulnerability) bool {
	for i := range vulnList {
		if vulnList[i].Details == vuln.Details && vulnList[i].Code == vuln.Code {
			vulnList[i].Occurrences = vulnList[i].Occurrences + 1
			return true
		}
	}
	return false
}

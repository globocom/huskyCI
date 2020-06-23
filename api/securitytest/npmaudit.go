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

// NpmAuditOutput is the struct that stores all npm audit output
type NpmAuditOutput struct {
	Advisories      map[string]Vulnerability `json:"advisories"`
	Metadata        Metadata                 `json:"metadata"`
	PackageNotFound bool
}

// Vulnerability is the granular output of a security info found
type Vulnerability struct {
	Findings           []Finding `json:"findings"`
	ID                 int       `json:"id"`
	ModuleName         string    `json:"module_name"`
	VulnerableVersions string    `json:"vulnerable_versions"`
	Severity           string    `json:"severity"`
	Overview           string    `json:"overview"`
	Title              string    `json:"title"`
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

func analyzeNpmaudit(npmAuditScan *SecTestScanInfo) error {

	npmAuditOutput := NpmAuditOutput{}
	npmAuditScan.FinalOutput = npmAuditOutput

	// if package-lock was not found, a warning will be genrated as a low vuln
	packageNotFound := strings.Contains(npmAuditScan.Container.COutput, "ERROR_PACKAGE_LOCK_NOT_FOUND")
	if packageNotFound {
		npmAuditScan.PackageNotFound = true
		npmAuditScan.prepareNpmAuditVulns()
		npmAuditScan.prepareContainerAfterScan()
		return nil
	}

	// nil cOutput states that no Issues were found.
	if npmAuditScan.Container.COutput == "" {
		npmAuditScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a NpmAuditOutput struct.
	if err := json.Unmarshal([]byte(npmAuditScan.Container.COutput), &npmAuditOutput); err != nil {
		log.Error("analyzeNpmaudit", "NPMAUDIT", 1014, npmAuditScan.Container.COutput, err)
		return err
	}
	npmAuditScan.FinalOutput = npmAuditOutput

	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	npmAuditScan.prepareNpmAuditVulns()
	npmAuditScan.prepareContainerAfterScan()
	return nil
}

func (npmAuditScan *SecTestScanInfo) prepareNpmAuditVulns() {

	huskyCInpmauditResults := types.HuskyCISecurityTestOutput{}
	npmAuditOutput := npmAuditScan.FinalOutput.(NpmAuditOutput)

	if npmAuditScan.PackageNotFound {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Severity = "low"
		npmauditVuln.Title = "No package-lock.json found."
		npmauditVuln.Details = "It looks like your project doesn't have a package-lock.json file. If you use NPM to handle your dependencies, it would be a good idea to commit it so huskyCI can check for vulnerabilities."

		npmAuditScan.Vulnerabilities.LowVulns = append(npmAuditScan.Vulnerabilities.LowVulns, npmauditVuln)
		return
	}

	for _, issue := range npmAuditOutput.Advisories {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Title = fmt.Sprintf("Vulnerable Dependency: %s %s (%s)", issue.ModuleName, issue.VulnerableVersions, issue.Title)
		npmauditVuln.Details = issue.Overview
		npmauditVuln.VunerableBelow = issue.VulnerableVersions
		npmauditVuln.Code = issue.ModuleName
		for _, findings := range issue.Findings {
			npmauditVuln.Version = findings.Version
		}

		switch issue.Severity {
		case "info", "low":
			npmauditVuln.Severity = "low"
			huskyCInpmauditResults.LowVulns = append(huskyCInpmauditResults.LowVulns, npmauditVuln)
		case "moderate":
			npmauditVuln.Severity = "medium"
			huskyCInpmauditResults.MediumVulns = append(huskyCInpmauditResults.MediumVulns, npmauditVuln)
		case "high", "critical":
			npmauditVuln.Severity = "high"
			huskyCInpmauditResults.HighVulns = append(huskyCInpmauditResults.HighVulns, npmauditVuln)
		}

	}

	npmAuditScan.Vulnerabilities = huskyCInpmauditResults
}

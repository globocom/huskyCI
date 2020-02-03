// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
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

func (s *SecurityTest) analyzeYarnaudit() error {

	// An empty container output states that no Issues were found.
	if s.Container.Output == "" {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	yarnAuditOutput := YarnAuditOutput{}

	// Unmarshall  container output into a YarnAuditOutput struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &yarnAuditOutput); err != nil {
		log.Error("analyzeYarnaudit", "YARNAUDIT", 1036, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1036]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareYarnAuditVulns(yarnAuditOutput)

	return nil
}

func (s *SecurityTest) prepareYarnAuditVulns(yarnAuditOutput YarnAuditOutput) {

	results := yarnAuditOutput.Advisories

	for _, issue := range results {

		yarnauditVuln := vulnerability.New()

		yarnauditVuln.Language = "JavaScript"
		yarnauditVuln.SecurityTest = "YarnAudit"
		yarnauditVuln.Details = issue.Overview
		yarnauditVuln.VunerableBelow = issue.VulnerableVersions
		yarnauditVuln.Code = issue.ModuleName
		for _, findings := range issue.Findings {
			yarnauditVuln.Version = findings.Version
		}

		switch issue.Severity {
		case "info", "low":
			yarnauditVuln.Severity = "low"
		case "moderate":
			yarnauditVuln.Severity = "medium"
		case "high", "critical":
			yarnauditVuln.Severity = "high"
		}

		s.Vulnerabilities = append(s.Vulnerabilities, *yarnauditVuln)

	}

}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
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

func (s *SecurityTest) analyzeNpmaudit() error {

	// An empty container output states that no Issues were found.
	if s.Container.Output == "" {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	npmAuditOutput := NpmAuditOutput{}

	// Unmarshall container output into a NpmAuditOutput struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &npmAuditOutput); err != nil {
		log.Error("analyzeNpmaudit", "NPMAUDIT", 1014, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1014]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareNpmAuditVulns(npmAuditOutput)

	return nil
}

func (s *SecurityTest) prepareNpmAuditVulns(npmAuditOutput NpmAuditOutput) {

	results := npmAuditOutput.Advisories

	for _, issue := range results {

		npmauditVuln := vulnerability.New()

		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTest = "NpmAudit"
		npmauditVuln.Details = issue.Overview
		npmauditVuln.VunerableBelow = issue.VulnerableVersions
		npmauditVuln.Code = issue.ModuleName
		for _, findings := range issue.Findings {
			npmauditVuln.Version = findings.Version
		}

		switch issue.Severity {
		case "info", "low":
			npmauditVuln.Severity = "low"
		case "moderate":
			npmauditVuln.Severity = "medium"
		case "high", "critical":
			npmauditVuln.Severity = "high"
		}

		s.Vulnerabilities = append(s.Vulnerabilities, *npmauditVuln)

	}

}

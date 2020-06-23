// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

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

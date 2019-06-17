// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//ResultsStruct is a struct that holds the array of results scanned.
type ResultsStruct struct {
	Results []RetirejsResult	`json:"results"`
}

//RetirejsResult is a struct that holds each scanned result.
type RetirejsResult struct {
	Component       string                    `json:"component"`
	Version         string                    `json:"version"`
	Level           int                       `json:"level"`
	Vulnerabilities []RetireJSVulnerabilities `json:"vulnerabilities"`
}

//RetireJSVulnerabilities is a struct that holds the vulnerabilities found on a scan.
type RetireJSVulnerabilities struct {
	Info        []string        					`json:"info"`
	Severity    string                             	`json:"severity"`
	Identifiers RetireJSVulnerabilityIdentifiers 	`json:"identifiers"`
}

//RetireJSVulnerabilityIdentifiers is a struct that holds identifiying information on a vulnerability found.
type RetireJSVulnerabilityIdentifiers struct {
	Summary string	`json:"summary"`
}

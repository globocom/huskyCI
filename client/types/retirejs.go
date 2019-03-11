// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "encoding/json"

//RetirejsOutput is the struct that holds issues, messages and errors found on a Retire scan.
type RetirejsOutput struct {
	RetirejsIssues []RetirejsIssue `json:"data"`
	Messages       json.RawMessage `json:"messages"`
	Errors         json.RawMessage `json:"errors"`
}

//RetirejsIssue is a struct that holds the results that were scanned and the file they came from.
type RetirejsIssue struct {
	File            string           `json:"file"`
	RetirejsResults []RetirejsResult `json:"results"`
}

//RetirejsResult is a struct that holds the vulnerabilities found on a component being used by the code being analysed.
type RetirejsResult struct {
	Version                 string                  `json:"version"`
	Component               string                  `json:"component"`
	Detection               string                  `json:"detection"`
	RetirejsVulnerabilities []RetirejsVulnerability `json:"vulnerabilities"`
}

//RetirejsVulnerability is a struct that holds info on what vulnerabilies were found.
type RetirejsVulnerability struct {
	Info                []string           `json:"info"`
	Below               string             `json:"below"`
	Severity            string             `json:"severity"`
	RetirejsIdentifiers RetirejsIdentifier `json:"identifiers"`
}

//RetirejsIdentifier is a struct that holds details on the vulnerabilities found.
type RetirejsIdentifier struct {
	IssueFound string   `json:"issue"`
	Summary    string   `json:"summary"`
	CVE        []string `json:"CVE"`
}

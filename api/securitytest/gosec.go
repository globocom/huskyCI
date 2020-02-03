// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
)

// GosecOutput is the struct that holds all data from Gosec output.
type GosecOutput struct {
	GosecIssues []GosecIssue `json:"Issues"`
	GosecStats  GosecStats   `json:"Stats"`
}

// GosecIssue is the struct that holds all issues from Gosec output.
type GosecIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// GosecStats is the struct that holds all stats from Gosec output.
type GosecStats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	Nosec int `json:"nosec"`
	Found int `json:"found"`
}

func (s *SecurityTest) analyzeGosec() error {

	// An empty container output states that no Issues were found.
	if s.Container.Output == "" {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	goSecOutput := GosecOutput{}

	//  Unmarshall container output into a Gosec struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &goSecOutput); err != nil {
		log.Error("analyzeGosec", "GOSEC", 1002, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1002]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareGosecVulns(goSecOutput)

	return nil
}

func (s *SecurityTest) prepareGosecVulns(gosecOutput GosecOutput) {

	results := gosecOutput.GosecIssues
	stats := gosecOutput.GosecStats

	for _, issue := range results {

		gosecVuln := vulnerability.New()

		gosecVuln.Language = "Go"
		gosecVuln.SecurityTest = "GoSec"
		gosecVuln.Severity = issue.Severity
		gosecVuln.Confidence = issue.Confidence
		gosecVuln.Details = issue.Details
		gosecVuln.File = issue.File
		gosecVuln.Line = issue.Line
		gosecVuln.Code = issue.Code

		s.Vulnerabilities = append(s.Vulnerabilities, *gosecVuln)
	}

	for i := 0; i <= stats.Nosec; i++ {

		gosecVuln := vulnerability.New()

		gosecVuln.Language = "Go"
		gosecVuln.SecurityTest = "GoSec"
		gosecVuln.Severity = "NOSEC"

		s.Vulnerabilities = append(s.Vulnerabilities, *gosecVuln)
	}

}

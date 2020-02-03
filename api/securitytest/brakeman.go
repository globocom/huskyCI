// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"strconv"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
)

// BrakemanOutput is the struct that holds issues and stats found on a Brakeman scan.
type BrakemanOutput struct {
	Results []BrakemanResult `json:"warnings"`
}

// BrakemanResult is the struct that holds all detailed information of a vulnerability found.
type BrakemanResult struct {
	Type       string `json:"warning_type"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Details    string `json:"link"`
	Confidence string `json:"confidence"`
}

func (s *SecurityTest) analyzeBrakeman() error {

	// An empty container output states that no Issues were found.
	if s.Container.Output == "" {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	brakemanOutput := BrakemanOutput{}

	// Unmarshall container output into a BrakemanOutput struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &brakemanOutput); err != nil {
		log.Error("analyzeBrakeman", "BRAKEMAN", 1005, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1005]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareBrakemanVulns(brakemanOutput)

	return nil
}

func (s *SecurityTest) prepareBrakemanVulns(brakemanOutput BrakemanOutput) {

	results := brakemanOutput.Results

	for _, issue := range results {

		brakemanVuln := vulnerability.New()

		brakemanVuln.Language = "Ruby"
		brakemanVuln.SecurityTest = "Brakeman"
		brakemanVuln.Confidence = issue.Confidence
		brakemanVuln.Details = issue.Details + issue.Message
		brakemanVuln.File = issue.File
		brakemanVuln.Line = strconv.Itoa(issue.Line)
		brakemanVuln.Code = issue.Code
		brakemanVuln.Type = issue.Type

		s.Vulnerabilities = append(s.Vulnerabilities, *brakemanVuln)

	}

}

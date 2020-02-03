// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/util"
	"github.com/globocom/huskyCI/api/vulnerability"
)

// BanditOutput is the struct that holds all data from Bandit output.
type BanditOutput struct {
	Results []BanditResult `json:"results"`
}

// BanditResult is the struct that holds detailed information of issues from Bandit output.
type BanditResult struct {
	Code            string `json:"code"`
	Filename        string `json:"filename"`
	IssueConfidence string `json:"issue_confidence"`
	IssueSeverity   string `json:"issue_severity"`
	IssueText       string `json:"issue_text"`
	LineNumber      int    `json:"line_number"`
	LineRange       []int  `json:"line_range"`
	TestID          string `json:"test_id"`
	TestName        string `json:"test_name"`
}

func (s *SecurityTest) analyzeBandit() error {

	banditOutput := BanditOutput{}

	// Unmarshall container output into a BanditOutput struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &banditOutput); err != nil {
		log.Error("analyzeBandit", "BANDIT", 1006, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1006]
		s.ErrorFound = err.Error()
		return err
	}

	// An empty Results slice states that no Issues were found.
	if len(banditOutput.Results) == 0 {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	s.prepareBanditVulns(banditOutput.Results)

	return nil
}

func (s *SecurityTest) prepareBanditVulns(results []BanditResult) {

	for _, issue := range results {

		banditVuln := vulnerability.New()

		isFalsePositive := checkFalsePositiveBandit(issue.Code, issue.LineNumber)
		if isFalsePositive {
			banditVuln.Nosec = true
			banditVuln.Severity = "NOSEC"
		} else {
			banditVuln.Severity = issue.IssueSeverity
		}

		banditVuln.Language = "Python"
		banditVuln.SecurityTest = "Bandit"
		banditVuln.Confidence = issue.IssueConfidence
		banditVuln.File = issue.Filename
		banditVuln.Line = strconv.Itoa(issue.LineNumber)
		banditVuln.Code = issue.Code
		banditVuln.Details = issue.IssueText

		s.Vulnerabilities = append(s.Vulnerabilities, *banditVuln)

	}

}

func checkFalsePositiveBandit(code string, lineNumber int) bool {
	lineNumberLength := util.CountDigits(lineNumber)
	splitCode := strings.Split(code, "\n")
	for _, codeLine := range splitCode {
		if len(codeLine) > 0 {
			codeLineNumber := codeLine[:lineNumberLength]
			if strings.Contains(codeLine, "#nohusky") && (codeLineNumber == strconv.Itoa(lineNumber)) {
				return true
			}
		}
	}
	return false
}

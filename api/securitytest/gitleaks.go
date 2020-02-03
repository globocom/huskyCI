// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
)

// GitleaksOutput is the struct that holds all data from Gitleaks output.
type GitleaksOutput struct {
	Results []GitLeaksResult `json:"results"`
}

// GitLeaksResult is the struct that holds all isssues from Gitleaks output.
type GitLeaksResult struct {
	Line          string `json:"line"`
	Commit        string `json:"commit"`
	Offender      string `json:"offender"`
	Rule          string `json:"rule"`
	Info          string `json:"info"`
	CommitMessage string `json:"commitMsg"`
	Author        string `json:"author"`
	Email         string `json:"email"`
	File          string `json:"file"`
	Repository    string `json:"repo"`
	Date          string `json:"date"`
	Tags          string `json:"tags"`
	Severity      string `json:"severity"`
}

func (s *SecurityTest) analyzeGitleaks() error {

	// An empty container output states that no Issues were found.
	if s.Container.Output == "" {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	gitLeaksOutput := GitleaksOutput{}

	// gitleaks took too long to finish
	if s.WarningFound != "" {
		s.prepareGitleaksVulns(gitLeaksOutput.Results)
		return nil
	}

	// Unmarshall container output into a GitleaksOutput struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &gitLeaksOutput); err != nil {
		log.Error("analyzeGitleaks", "GITLEAKS", 1038, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1038]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareGitleaksVulns(gitLeaksOutput.Results)

	return nil
}

func (s *SecurityTest) prepareGitleaksVulns(results []GitLeaksResult) {

	if s.WarningFound != "" {

		gitleaksVuln := vulnerability.New()

		gitleaksVuln.Language = "Generic"
		gitleaksVuln.SecurityTest = "GitLeaks"
		gitleaksVuln.Details = s.WarningFound

		s.Vulnerabilities = append(s.Vulnerabilities, *gitleaksVuln)

		return
	}

	for _, issue := range results {

		// dependencies issues will not checked at this moment by huskyCI
		if strings.Contains(issue.File, "vendor/") || strings.Contains(issue.File, "node_modules/") {
			continue
		}

		gitleaksVuln := vulnerability.New()

		gitleaksVuln.Language = "Generic"
		gitleaksVuln.SecurityTest = "GitLeaks"
		gitleaksVuln.File = issue.File
		gitleaksVuln.Line = issue.Line

		switch issue.Rule {
		case "PKCS8", "RSA", "SSH", "PGP", "EC":
			gitleaksVuln.Severity = "HIGH"
		case "AWS Secret Key", "Facebook Secret Key", "Facebook access token", "Twitter Secret Key", "LinkedIn Secret Key", "Google OAuth access token", "Google Cloud Platform API key", "Heroku API key", "MailChimp API key", "Mailgun API key", "PayPal Braintree access token", "Picatic API key", "Stripe API key", "Twilio API key":
			gitleaksVuln.Severity = "MEDIUM"
		default:
			gitleaksVuln.Severity = "LOW"
		}

		gitleaksVuln.Details = issue.Rule + " @ [" + issue.Commit + "]"

		s.Vulnerabilities = append(s.Vulnerabilities, *gitleaksVuln)

	}

}

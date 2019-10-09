// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// GitleaksOutput is the struct that holds all data from Gitleaks output.
type GitleaksOutput []GitLeaksIssue

// GitLeaksIssue is the struct that holds all isssues from Gitleaks output.
type GitLeaksIssue struct {
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

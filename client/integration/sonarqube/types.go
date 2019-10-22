// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sonarqube

// HuskyCISonarOutput is the struct that holds the Sonar output
type HuskyCISonarOutput struct {
	Issues []SonarIssue `json:"issues"`
}

// SonarIssue is the struct that holds a single Sonar issue
type SonarIssue struct {
	EngineID           string          `json:"engineId"`
	RuleID             string          `json:"ruleId"`
	PrimaryLocation    SonarLocation   `json:"primaryLocation"`
	Type               string          `json:"type"`
	Severity           string          `json:"severity"`
	EffortMinutes      int             `json:"effortMinutes,omitempty"`
	SecondaryLocations []SonarLocation `json:"secondaryLocations,omitempty"`
}

// SonarLocation is the struct that holds a vulnerability location within code
type SonarLocation struct {
	Message   string         `json:"message"`
	FilePath  string         `json:"filePath"`
	TextRange SonarTextRange `json:"textRange"`
}

// SonarTextRange is the struct that holds addtional location fields
type SonarTextRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

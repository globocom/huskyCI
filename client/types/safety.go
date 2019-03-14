// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//SafetyOutput is the struct that holds issues, messages and errors found on a Safety scan.
type SafetyOutput struct {
	SafetyIssues []SafetyIssue `json:"issues"`
}

//SafetyIssue is a struct that holds the results that were scanned and the file they came from.
type SafetyIssue struct {
	Dependency string `json:"dependency"`
	Below      string `json:"vulnerable_below"`
	Version    string `json:"installed_version"`
	Comment    string `json:"description"`
	ID         string `json:"id"`
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// BrakemanOutput is the struct that holds issues and stats found on a Brakeman scan.
type BrakemanOutput struct {
	Warnings []WarningItem `json:"warnings"`
	IgnoredWarnings []WarningItem `json:"ignored_warnings"`
}

// WarningItem is the struct that holds all detailed information of a vulnerability found.
type WarningItem struct {
	Type       string `json:"warning_type"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Details    string `json:"link"`
	Confidence string `json:"confidence"`
}

func analyzeBrakeman(brakemanScan *SecTestScanInfo) error {

	brakemanOutput := BrakemanOutput{}

	// nil cOutput states that no Issues were found.
	if brakemanScan.Container.COutput == "" {
		brakemanScan.prepareContainerAfterScan()
		return nil
	}
	// Unmarshall rawOutput into finalOutput, that is a Brakeman struct.
	if err := json.Unmarshal([]byte(brakemanScan.Container.COutput), &brakemanOutput); err != nil {
		log.Error("analyzeBrakeman", "BRAKEMAN", 1005, brakemanScan.Container.COutput, err)
		brakemanScan.ErrorFound = err
		return err
	}
	brakemanScan.FinalOutput = brakemanOutput

	// check results and prepare all vulnerabilities found
	brakemanScan.prepareBrakemanVulns()
	brakemanScan.prepareContainerAfterScan()
	return nil
}

func (brakemanScan *SecTestScanInfo) prepareBrakemanVulns() {

	huskyCIbrakemanResults := types.HuskyCISecurityTestOutput{}
	brakemanOutput := brakemanScan.FinalOutput.(BrakemanOutput)

	for _, warning := range brakemanOutput.Warnings {
		brakemanVuln := types.HuskyCIVulnerability{}
		brakemanVuln.Language = "Ruby"
		brakemanVuln.SecurityTool = "Brakeman"
		brakemanVuln.Confidence = warning.Confidence
		brakemanVuln.Title = fmt.Sprintf("Vulnerable Dependency: %s %s", warning.Type, warning.Message)
		brakemanVuln.Details = warning.Details
		brakemanVuln.File = warning.File
		brakemanVuln.Line = strconv.Itoa(warning.Line)
		brakemanVuln.Code = warning.Code
		brakemanVuln.Type = warning.Type

		switch brakemanVuln.Confidence {
		case "High":
			huskyCIbrakemanResults.HighVulns = append(huskyCIbrakemanResults.HighVulns, brakemanVuln)
		case "Medium":
			huskyCIbrakemanResults.MediumVulns = append(huskyCIbrakemanResults.MediumVulns, brakemanVuln)
		case "Low":
			huskyCIbrakemanResults.LowVulns = append(huskyCIbrakemanResults.LowVulns, brakemanVuln)
		}
	}
	for _, ignoredWarning := range brakemanOutput.IgnoredWarnings {
		brakemanVuln := types.HuskyCIVulnerability{}
		brakemanVuln.Language = "Ruby"
		brakemanVuln.SecurityTool = "Brakeman"
		brakemanVuln.Confidence = ignoredWarning.Confidence
		brakemanVuln.Title = fmt.Sprintf("Vulnerable Dependency: %s %s", ignoredWarning.Type, ignoredWarning.Message)
		brakemanVuln.Severity = "NOSEC"
		brakemanVuln.Details = ignoredWarning.Details
		brakemanVuln.File = ignoredWarning.File
		brakemanVuln.Line = strconv.Itoa(ignoredWarning.Line)
		brakemanVuln.Code = ignoredWarning.Code
		brakemanVuln.Type = ignoredWarning.Type
		huskyCIbrakemanResults.NoSecVulns = append(huskyCIbrakemanResults.NoSecVulns, brakemanVuln)
	}

	brakemanScan.Vulnerabilities = huskyCIbrakemanResults
}

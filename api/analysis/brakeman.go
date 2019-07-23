// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// BrakemanOutput is the struct that holds issues and stats found on a Brakeman scan.
type BrakemanOutput struct {
	Warnings []WarningItem `json:"warnings"`
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

// BrakemanCheckOutputFlow analyses the output from Brakeman and sets a cResult based on it.
func BrakemanCheckOutputFlow(CID string, cOutput string, RID string) {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	if errorClonning {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: nil cOutput states that no Issues were found.
	if cOutput == "" {
		if err := updateInfoAndResultBasedOnCID("No issues found.", "passed", CID); err != nil {
			return
		}
		return
	}

	// step 3: Unmarshall cOutput into BrakemanOutput struct.
	brakemanOutput := BrakemanOutput{}
	err := json.Unmarshal([]byte(cOutput), &brakemanOutput)
	if err != nil {
		log.Error("BrakemanStartAnalysis", "BRAKEMAN", 1005, cOutput, err)
		return
	}

	// step 4: An empty errors slice also means that no vulnerabilities were found
	if len(brakemanOutput.Warnings) == 0 {
		if err := updateInfoAndResultBasedOnCID("No issues found.", "passed", CID); err != nil {
			return
		}
		return
	}

	// step 5: find Issues that have confidence "High" or "Medium".
	cResult := "passed"
	issueMessage := "No issues found."
	for _, warning := range brakemanOutput.Warnings {
		if warning.Confidence == "High" || warning.Confidence == "Medium" {
			cResult = "failed"
			issueMessage = "Issues found."
			break
		}
	}
	if err := updateInfoAndResultBasedOnCID(issueMessage, cResult, CID); err != nil {
		return
	}

	// step 6: finally, update analysis with huskyCI results
	if err := updateHuskyCIResultsBasedOnRID(RID, "brakeman", brakemanOutput); err != nil {
		return
	}
}

// prepareHuskyCIBrakemanResults will prepare Brakeman output to be added into RubyResults struct
func prepareHuskyCIBrakemanResults(brakemanOutput BrakemanOutput) types.HuskyCIBrakemanOutput {

	var huskyCIbrakemanResults types.HuskyCIBrakemanOutput

	for _, warning := range brakemanOutput.Warnings {
		brakemanVuln := types.HuskyCIVulnerability{}
		brakemanVuln.Language = "Ruby"
		brakemanVuln.SecurityTool = "Brakeman"
		brakemanVuln.Confidence = warning.Confidence
		brakemanVuln.Details = warning.Details + warning.Message
		brakemanVuln.File = warning.File
		brakemanVuln.Line = strconv.Itoa(warning.Line)
		brakemanVuln.Code = warning.Code
		brakemanVuln.Type = warning.Type

		switch brakemanVuln.Confidence {
		case "High":
			huskyCIbrakemanResults.LowVulnsBrakeman = append(huskyCIbrakemanResults.LowVulnsBrakeman, brakemanVuln)
		case "Medium":
			huskyCIbrakemanResults.MediumVulnsBrakeman = append(huskyCIbrakemanResults.MediumVulnsBrakeman, brakemanVuln)
		case "Low":
			huskyCIbrakemanResults.HighVulnsBrakeman = append(huskyCIbrakemanResults.HighVulnsBrakeman, brakemanVuln)
		}
	}

	return huskyCIbrakemanResults
}

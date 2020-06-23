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

// TFSecOutput is the struct that holds all data from TFSec output.
type TFSecOutput struct {
	Warnings json.RawMessage `json:"warnings"`
	Results  []TFSecResult   `json:"results"`
}

// TFSecResult is the struct that holds detailed information of results from TFSec output.
type TFSecResult struct {
	RuleID      string   `json:"rule_id"`
	Link        string   `json:"link"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
}

// Location is the struct that holds detailed information of location from each result
type Location struct {
	Filename  string `json:"filename"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

func analyzeTFSec(tfsecScan *SecTestScanInfo) error {

	tfsecOutput := TFSecOutput{}

	// Unmarshall rawOutput into finalOutput, that is a TFSec struct.
	if err := json.Unmarshal([]byte(tfsecScan.Container.COutput), &tfsecOutput); err != nil {
		log.Error("analyzeTFSec", "TFSEC", 1040, tfsecScan.Container.COutput, err)
		tfsecScan.ErrorFound = err
		return err
	}
	tfsecScan.FinalOutput = tfsecOutput

	// an empty Results slice states that no Issues were found.
	if tfsecOutput.Results == nil {
		tfsecScan.prepareContainerAfterScan()
		return nil
	}

	// check results and prepare all vulnerabilities found
	tfsecScan.prepareTFSecVulns()
	tfsecScan.prepareContainerAfterScan()
	return nil
}

func (tfsecScan *SecTestScanInfo) prepareTFSecVulns() {

	huskyCItfsecResults := types.HuskyCISecurityTestOutput{}
	tfsecOutput := tfsecScan.FinalOutput.(TFSecOutput)

	for _, result := range tfsecOutput.Results {
		tfsecVuln := types.HuskyCIVulnerability{}
		tfsecVuln.Language = "HCL"
		tfsecVuln.SecurityTool = "TFSec"
		tfsecVuln.Severity = result.Severity
		tfsecVuln.Title = result.Description
		tfsecVuln.Details = result.RuleID + " @ [" + result.Description + "]"
		startLine := strconv.Itoa(result.Location.StartLine)
		endLine := strconv.Itoa(result.Location.EndLine)
		tfsecVuln.Line = startLine
		tfsecVuln.Code = fmt.Sprintf("Code beetween Line %s and Line %s.", startLine, endLine)
		tfsecVuln.File = result.Location.Filename

		switch tfsecVuln.Severity {
		case "INFO":
			tfsecVuln.Severity = "Low"
			huskyCItfsecResults.LowVulns = append(huskyCItfsecResults.LowVulns, tfsecVuln)
		case "WARNING":
			tfsecVuln.Severity = "Medium"
			huskyCItfsecResults.MediumVulns = append(huskyCItfsecResults.MediumVulns, tfsecVuln)
		case "ERROR":
			tfsecVuln.Severity = "High"
			huskyCItfsecResults.HighVulns = append(huskyCItfsecResults.HighVulns, tfsecVuln)
		}
	}

	tfsecScan.Vulnerabilities = huskyCItfsecResults
}

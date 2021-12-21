// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// SecurityCodeScanOutput is the struct that holds all data from SecurityCodeScan output.
type SecurityCodeScanOutput struct {
	Schema  string                 `json:"$schema"`
	Version string                 `json:"version"`
	Runs    []SecurityCodeScanRuns `json:"runs"`
}

// SecurityCodeScanRuns is the struct that holds detailed information of Runs from SecurityCodeScan.
type SecurityCodeScanRuns struct {
	Results []SecurityCodeScanResult `json:"results"`
}

// SecurityCodeScanResult is the struct that holds detailed information of Result from SecurityCodeScan.
type SecurityCodeScanResult struct {
	RuleID  string `json:"ruleId"`
	Level   string `json:"level"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
	Locations []SecurityCodeScanLocation `json:"locations"`
}

// SecurityCodeScanLocation is the struct that holds detailed information of locations from SecurityCodeScan.
type SecurityCodeScanLocation struct {
	PhysicalLocation struct {
		ArtifactLocation struct {
			URI string `json:"uri"`
		} `json:"artifactLocation"`
	} `json:"physicalLocation"`
	Region struct {
		StartLine   int `json:"startLine"`
		StartColumn int `json:"startColumn"`
		EndLine     int `json:"endLine"`
		EndColumn   int `json:"endColumn"`
	} `json:"region"`
}

func analyzeSecurityCodeScan(securitycodescanScan *SecTestScanInfo) error {

	securitycodescanOutput := SecurityCodeScanOutput{}
	securitycodescanScan.FinalOutput = securitycodescanOutput

	// if security-scan fails, a warning will be genrated as a low vuln
	errorRunning := strings.Contains(securitycodescanScan.Container.COutput, "ERROR_SECURITY_CODE_SCAN_RUNNING")
	if errorRunning {
		securitycodescanScan.SecurityCodeScanErrorRunning = true
		securitycodescanScan.prepareSecurityCodeScanVulns()
		securitycodescanScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a SecurityCodeScanOutput struct.
	if err := json.Unmarshal([]byte(securitycodescanScan.Container.COutput), &securitycodescanOutput); err != nil {
		log.Error("analyzeSecurityCodeScan", "SecurityCodeScan", 1041, securitycodescanScan.Container.COutput, err)
		securitycodescanScan.ErrorFound = err
		return err
	}
	securitycodescanScan.FinalOutput = securitycodescanOutput

	// an empty Results slice states that no Issues were found.
	if (len(securitycodescanOutput.Runs) > 0) && (len(securitycodescanOutput.Runs[0].Results) == 0) {
		securitycodescanScan.prepareContainerAfterScan()
		return nil
	}

	// check results and prepare all vulnerabilities found
	securitycodescanScan.prepareSecurityCodeScanVulns()
	securitycodescanScan.prepareContainerAfterScan()
	return nil
}

func (s *SecTestScanInfo) prepareSecurityCodeScanVulns() {

	huskyCISecurityCodeScanResults := types.HuskyCISecurityTestOutput{}
	securityCodeScanOutput := s.FinalOutput.(SecurityCodeScanOutput)

	if s.SecurityCodeScanErrorRunning {
		securityCodeScanVuln := types.HuskyCIVulnerability{}
		securityCodeScanVuln.Language = "C#"
		securityCodeScanVuln.SecurityTool = "Security Code Scan"
		securityCodeScanVuln.Severity = "low"
		securityCodeScanVuln.Title = "Error running Security Code Scan Tool."
		securityCodeScanVuln.Details = "It looks like huskyCI could not run 'security-scan' on your project. No .sln file was found on your project or an unknown file extension was loaded on your .sln file."

		s.Vulnerabilities.LowVulns = append(s.Vulnerabilities.LowVulns, securityCodeScanVuln)
		return
	}

	for _, result := range securityCodeScanOutput.Runs[0].Results {
		securityCodeScanVuln := types.HuskyCIVulnerability{}
		securityCodeScanVuln.Language = "C#"
		securityCodeScanVuln.SecurityTool = "Security Code Scan"
		securityCodeScanVuln.Severity = result.Level
		securityCodeScanVuln.Title = result.RuleID
		securityCodeScanVuln.Details = result.Message.Text
		if len(result.Locations) > 0 {
			startLine := strconv.Itoa(result.Locations[0].Region.StartLine)
			endLine := strconv.Itoa(result.Locations[0].Region.EndLine)
			securityCodeScanVuln.Line = startLine
			securityCodeScanVuln.Code = fmt.Sprintf("Code beetween Line %s and Line %s.", startLine, endLine)
			pathSlice := strings.Split(result.Locations[0].PhysicalLocation.ArtifactLocation.URI, "code/")
			if len(pathSlice) > 1 {
				securityCodeScanVuln.File = pathSlice[1]
			} else {
				securityCodeScanVuln.File = result.Locations[0].PhysicalLocation.ArtifactLocation.URI
			}
		}

		switch securityCodeScanVuln.Severity {
		case "recommendation":
			securityCodeScanVuln.Severity = "Low"
			huskyCISecurityCodeScanResults.LowVulns = append(huskyCISecurityCodeScanResults.LowVulns, securityCodeScanVuln)
		case "warning":
			securityCodeScanVuln.Severity = "Medium"
			huskyCISecurityCodeScanResults.MediumVulns = append(huskyCISecurityCodeScanResults.MediumVulns, securityCodeScanVuln)
		case "error":
			securityCodeScanVuln.Severity = "High"
			huskyCISecurityCodeScanResults.HighVulns = append(huskyCISecurityCodeScanResults.HighVulns, securityCodeScanVuln)
		}
	}

	s.Vulnerabilities = huskyCISecurityCodeScanResults
}

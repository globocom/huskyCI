// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// SpotBugsOutput is the struct that holds all data from SpotBugs output.
type SpotBugsOutput []SpotBugsIssue

// SpotBugsIssue is the struct that holds all issues from SpotBugs output.
type SpotBugsIssue struct {
	Type         string `json:"type"`
	Priority     string `json:"priority"`
	Rank         string `json:"rank"`
	Abbreviation string `json:"abbrev"`
	Category     string `json:"category"`
	// Files []SpotBugFile
}

// type SpotBugFile struct{

// }

func analyzeSpotBugs(spotbugsScan *SecTestScanInfo) error {

	goSecOutput := SpotBugsOutput{}
	spotbugsScan.FinalOutput = goSecOutput

	// nil cOutput states that no Issues were found.
	if spotbugsScan.Container.COutput == "" {
		spotbugsScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a SpotBugsOutput struct.
	if err := json.Unmarshal([]byte(spotbugsScan.Container.COutput), &goSecOutput); err != nil {
		log.Error("analyzeSpotBugs", "SPOTBUGS", 1039, spotbugsScan.Container.COutput, err)
		spotbugsScan.ErrorFound = err
		spotbugsScan.prepareContainerAfterScan()
		return err
	}
	spotbugsScan.FinalOutput = goSecOutput

	// check results and prepare all vulnerabilities found
	spotbugsScan.prepareSpotBugsVulns()
	spotbugsScan.prepareContainerAfterScan()
	return nil
}

func (spotbugsScan *SecTestScanInfo) prepareSpotBugsVulns() {

	huskyCIspotbugsResults := types.HuskyCISecurityTestOutput{}
	spotbugsOutput := spotbugsScan.FinalOutput.(SpotBugsOutput)

	for _, issue := range spotbugsOutput {
		spotbugsVuln := types.HuskyCIVulnerability{}
		spotbugsVuln.Language = "Java"
		spotbugsVuln.SecurityTool = "SpotBugs"
		spotbugsVuln.Severity = issue.Priority
		// spotbugsVuln.Confidence = issue.Confidence
		// spotbugsVuln.Details = issue.Details
		// spotbugsVuln.File = issue.File
		// spotbugsVuln.Line = issue.Line
		// spotbugsVuln.Code = issue.Code

		switch spotbugsVuln.Severity {
		case "LOW":
			huskyCIspotbugsResults.LowVulns = append(huskyCIspotbugsResults.LowVulns, spotbugsVuln)
		case "MEDIUM":
			huskyCIspotbugsResults.MediumVulns = append(huskyCIspotbugsResults.MediumVulns, spotbugsVuln)
		case "HIGH":
			huskyCIspotbugsResults.HighVulns = append(huskyCIspotbugsResults.HighVulns, spotbugsVuln)
		}
	}

	spotbugsScan.Vulnerabilities = huskyCIspotbugsResults
}

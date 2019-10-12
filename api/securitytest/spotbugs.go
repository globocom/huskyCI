// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/xml"
	"strconv"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// SpotBugsOutput is the struct that holds all data from SpotBugs output.
type SpotBugsOutput struct {
	XMLName       xml.Name        `xml:"BugCollection"`
	Project       Project         `xml:"Project"`
	SpotBugsIssue []SpotBugsIssue `xml:"BugInstance"`
	Errors        Error           `xml:"Errors"`
}

// Project is the struct that holds data about project
type Project struct {
	XMLName xml.Name `xml:"Project"`
	Name    string   `xml:"projectName,attr"`
	Jar     string   `xml:"Jar"`
	Plugin  string   `xml:"Plugin"`
}

// SpotBugsIssue is the struct that holds all issues from SpotBugs output.
type SpotBugsIssue struct {
	XMLName      xml.Name     `xml:"BugInstance"`
	Type         string       `xml:"type,attr"`
	Priority     string       `xml:"priority,attr"`
	Rank         string       `xml:"rank,attr"`
	Abbreviation string       `xml:"abbrev,attr"`
	Category     string       `xml:"category,attr"`
	SourceLine   []SourceLine `xml:"SourceLine"`
}

// Error is the struct that holds errors that happened in analysis
type Error struct {
	XMLName        xml.Name `xml:"Errors"`
	Errors         string   `xml:"errors,attr"`
	MissingClasses string   `xml:"missingClasses,attr"`
}

// SourceLine is the struct that holds details about issue location
type SourceLine struct {
	XMLName       xml.Name `xml:"SourceLine"`
	ClassName     string   `xml:"classname,attr"`
	Start         string   `xml:"start,attr"`
	End           string   `xml:"end,attr"`
	StartByteCode string   `xml:"startBytecode,attr"`
	EndByteCode   string   `xml:"endBytecode,attr"`
	SourceFile    string   `xml:"sourcefile,attr"`
	SourcePath    string   `xml:"sourcepath,attr"`
}

func analyzeSpotBugs(spotbugsScan *SecTestScanInfo) error {

	spotBugsOutput := SpotBugsOutput{}
	spotbugsScan.FinalOutput = spotBugsOutput

	// nil cOutput states that no Issues were found.
	if spotbugsScan.Container.COutput == "" {
		spotbugsScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a SpotBugsOutput struct.
	spotBugsOutput, err := parseXMLtoJSON([]byte(spotbugsScan.Container.COutput))
	if err != nil {
		log.Error("analyzeSpotBugs", "SPOTBUGS", 1039, spotbugsScan.Container.COutput, err)
		spotbugsScan.ErrorFound = err
		spotbugsScan.prepareContainerAfterScan()
		return err
	}

	spotbugsScan.FinalOutput = spotBugsOutput

	// check results and prepare all vulnerabilities found
	spotbugsScan.prepareSpotBugsVulns()
	spotbugsScan.prepareContainerAfterScan()
	return nil
}
func parseXMLtoJSON(byteValue []byte) (SpotBugsOutput, error) {
	var bugs SpotBugsOutput
	if err := xml.Unmarshal(byteValue, &bugs); err != nil {
		return bugs, err
	}
	// if numOfErrors, err := strconv.Atoi(bugs.Errors.Errors); err == nil {
	// 	if numOfErrors > 0 {
	// 		fmt.Println("Errors happened")
	// 	}
	// }
	// if numOfMissingClasses, err := strconv.Atoi(bugs.Errors.MissingClasses); err == nil {
	// 	if numOfMissingClasses > 0 {
	// 		fmt.Println("missing classes")
	// 	}
	// }
	return bugs, nil
}

func (spotbugsScan *SecTestScanInfo) prepareSpotBugsVulns() {
	var spotbugsOutput SpotBugsOutput

	huskyCIspotbugsResults := types.HuskyCISecurityTestOutput{}
	spotbugsOutput = spotbugsScan.FinalOutput.(SpotBugsOutput)

	for i := 0; i < len(spotbugsOutput.SpotBugsIssue); i++ {
		for j := 0; j < len(spotbugsOutput.SpotBugsIssue[i].SourceLine); j++ {
			spotbugsVuln := types.HuskyCIVulnerability{}
			spotbugsVuln.Language = "Java"
			spotbugsVuln.SecurityTool = "SpotBugs"
			spotbugsVuln.Details = spotbugsOutput.SpotBugsIssue[i].Type
			spotbugsVuln.File = spotbugsOutput.SpotBugsIssue[i].SourceLine[j].SourcePath
			spotbugsVuln.Line = spotbugsOutput.SpotBugsIssue[i].SourceLine[j].Start

			switch spotbugsOutput.SpotBugsIssue[i].Priority {
			case "1":
				spotbugsVuln.Confidence = "HIGH"
			case "2":
				spotbugsVuln.Confidence = "MEDIUM"
			default:
				spotbugsVuln.Confidence = "LOW"
			}

			switch rank, _ := strconv.Atoi(spotbugsOutput.SpotBugsIssue[i].Rank); {
			case rank < 10:
				huskyCIspotbugsResults.HighVulns = append(huskyCIspotbugsResults.HighVulns, spotbugsVuln)
			case rank < 15:
				huskyCIspotbugsResults.MediumVulns = append(huskyCIspotbugsResults.MediumVulns, spotbugsVuln)
			default:
				huskyCIspotbugsResults.LowVulns = append(huskyCIspotbugsResults.LowVulns, spotbugsVuln)
			}
		}
	}

	spotbugsScan.Vulnerabilities = huskyCIspotbugsResults
}

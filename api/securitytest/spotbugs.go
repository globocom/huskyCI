// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

const (
	mediumSeverityValue = 15
	highSeverityValue   = 10
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
	ShortMessage string       `xml:"ShortMessage"`
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

	// check if GRADLE failed running
	errorGradleFound := strings.Contains(spotbugsScan.Container.COutput, "ERROR_RUNNING_GRADLE_BUILD")
	if errorGradleFound {
		spotbugsScan.ErrorFound = errors.New("error running gradle")
		spotbugsScan.prepareSpotBugsVulns()
		spotbugsScan.prepareContainerAfterScan()
		return nil
	}

	// check if Maven failed running
	errorMavenFound := strings.Contains(spotbugsScan.Container.COutput, "ERROR_RUNNING_MAVEN_BUILD")
	if errorMavenFound {
		spotbugsScan.ErrorFound = errors.New("error running maven")
		spotbugsScan.prepareSpotBugsVulns()
		spotbugsScan.prepareContainerAfterScan()
		return nil
	}

	// check if unsuported java project was found
	errorUnsuported := strings.Contains(spotbugsScan.Container.COutput, "ERROR_UNSUPPORTED_JAVA_PROJECT")
	if errorUnsuported {
		spotbugsScan.ErrorFound = errors.New("error unsuported java project")
		spotbugsScan.prepareSpotBugsVulns()
		spotbugsScan.prepareContainerAfterScan()
		return nil
	}

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
	numOfErrors, err := strconv.Atoi(bugs.Errors.Errors)
	if err != nil {
		return bugs, err
	}
	numOfMissingClasses, err := strconv.Atoi(bugs.Errors.MissingClasses)
	if err != nil {
		return bugs, err
	}
	if (len(bugs.SpotBugsIssue) == 0) && (numOfErrors > 0 || numOfMissingClasses > 0) {
		return bugs, fmt.Errorf("spotbugs has risen errors because of [%d] missing classes and [%d] errors while analysing", numOfMissingClasses, numOfErrors)
	}
	return bugs, nil
}

func (spotbugsScan *SecTestScanInfo) prepareSpotBugsVulns() {
	var spotbugsOutput SpotBugsOutput

	huskyCIspotbugsResults := types.HuskyCISecurityTestOutput{}
	spotbugsOutput = spotbugsScan.FinalOutput.(SpotBugsOutput)

	if spotbugsScan.ErrorFound != nil {
		spotbugsVuln := types.HuskyCIVulnerability{}
		spotbugsVuln.Language = "Java"
		spotbugsVuln.SecurityTool = "SpotBugs"
		spotbugsVuln.Title = "Error while running SpotBugs scan."
		spotbugsVuln.Details = fmt.Sprintf("An error occured running huskyCI scan on your Java project: %s", spotbugsScan.ErrorFound.Error())
		spotbugsVuln.Severity = "LOW"
		spotbugsVuln.Confidence = "HIGH"

		spotbugsScan.Vulnerabilities.LowVulns = append(huskyCIspotbugsResults.LowVulns, spotbugsVuln)
		return
	}

	for i := 0; i < len(spotbugsOutput.SpotBugsIssue); i++ {
		for j := 0; j < len(spotbugsOutput.SpotBugsIssue[i].SourceLine); j++ {
			spotbugsVuln := types.HuskyCIVulnerability{}
			spotbugsVuln.Language = "Java"
			spotbugsVuln.SecurityTool = "SpotBugs"
			spotbugsVuln.Type = spotbugsOutput.SpotBugsIssue[i].Abbreviation
			spotbugsVuln.Details = spotbugsOutput.SpotBugsIssue[i].Type
			startLine := spotbugsOutput.SpotBugsIssue[i].SourceLine[j].Start
			endLine := spotbugsOutput.SpotBugsIssue[i].SourceLine[j].End
			spotbugsVuln.Code = fmt.Sprintf("Code beetween Line %s and Line %s.", startLine, endLine)
			spotbugsVuln.Line = startLine
			spotbugsVuln.File = spotbugsOutput.SpotBugsIssue[i].SourceLine[j].SourcePath
			spotbugsVuln.Title = spotbugsVuln.Details

			switch spotbugsOutput.SpotBugsIssue[i].Priority {
			case "1":
				spotbugsVuln.Confidence = "HIGH"
			case "2":
				spotbugsVuln.Confidence = "MEDIUM"
			default:
				spotbugsVuln.Confidence = "LOW"
			}

			rank, err := strconv.Atoi(spotbugsOutput.SpotBugsIssue[i].Rank)
			if err != nil {
				log.Warning("analyzeSpotBugs", "SPOTBUGS", 1039, "exception while reading rank from a spotbugs issue", err)
				continue
			}

			switch {
			case rank < highSeverityValue:
				spotbugsVuln.Severity = "HIGH"
				huskyCIspotbugsResults.HighVulns = append(huskyCIspotbugsResults.HighVulns, spotbugsVuln)
			case rank < mediumSeverityValue:
				spotbugsVuln.Severity = "MEDIUM"
				huskyCIspotbugsResults.MediumVulns = append(huskyCIspotbugsResults.MediumVulns, spotbugsVuln)
			default:
				spotbugsVuln.Severity = "LOW"
				huskyCIspotbugsResults.LowVulns = append(huskyCIspotbugsResults.LowVulns, spotbugsVuln)
			}
		}
	}

	spotbugsScan.Vulnerabilities = huskyCIspotbugsResults
}

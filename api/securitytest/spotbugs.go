// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
)

const (
	mediumSeverityValue = 15
	highSeverityValue   = 10
)

// SpotBugsOutput is the struct that holds all data from SpotBugs output.
type SpotBugsOutput struct {
	XMLName        xml.Name        `xml:"BugCollection"`
	Project        Project         `xml:"Project"`
	SpotBugsIssues []SpotBugsIssue `xml:"BugInstance"`
	Errors         Error           `xml:"Errors"`
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

func (s *SecurityTest) analyzeSpotBugs() error {

	// An empty container output states that no Issues were found.
	if s.Container.Output == "" {
		s.Result = "passed"
		s.Info = "No issues found."
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a SpotBugsOutput struct.
	spotBugsOutput, err := parseXMLtoJSON([]byte(s.Container.Output))
	if err != nil {
		log.Error("analyzeSpotBugs", "SPOTBUGS", 1039, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1039]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareSpotBugsVulns(spotBugsOutput)
	return nil
}

func (s *SecurityTest) prepareSpotBugsVulns(spotBugsOutput SpotBugsOutput) {

	results := spotBugsOutput.SpotBugsIssues

	for _, issue := range results {
		for _, sourceLine := range issue.SourceLine {

			spotbugsVuln := vulnerability.New()

			spotbugsVuln.Language = "Java"
			spotbugsVuln.SecurityTest = "SpotBugs"
			spotbugsVuln.Type = issue.Abbreviation
			spotbugsVuln.Details = issue.Type
			spotbugsVuln.Code = fmt.Sprintf("Code beetween Line %s and Line %s.", sourceLine.Start, sourceLine.End)
			spotbugsVuln.Line = sourceLine.Start
			spotbugsVuln.File = sourceLine.SourcePath

			switch issue.Priority {
			case "1":
				spotbugsVuln.Confidence = "HIGH"
			case "2":
				spotbugsVuln.Confidence = "MEDIUM"
			default:
				spotbugsVuln.Confidence = "LOW"
			}

			rank, err := strconv.Atoi(issue.Rank)
			if err != nil {
				log.Warning("analyzeSpotBugs", "SPOTBUGS", 1039, "exception while reading rank from a spotbugs issue", err)
				continue
			}

			switch {
			case rank < highSeverityValue:
				spotbugsVuln.Severity = "HIGH"
			case rank < mediumSeverityValue:
				spotbugsVuln.Severity = "MEDIUM"
			default:
				spotbugsVuln.Severity = "LOW"
			}

			s.Vulnerabilities = append(s.Vulnerabilities, *spotbugsVuln)

		}

	}

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
	if (len(bugs.SpotBugsIssues) == 0) && (numOfErrors > 0 || numOfMissingClasses > 0) {
		return bugs, fmt.Errorf("spotbugs has risen errors because of [%d] missing classes and [%d] errors while analysing", numOfMissingClasses, numOfErrors)
	}
	return bugs, nil
}

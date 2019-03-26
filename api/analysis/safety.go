// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"gopkg.in/mgo.v2/bson"
)

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

//SafetyStartAnalysis analyses the output from Safety and sets cResult based on it.
func SafetyStartAnalysis(CID string, cOutput string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	warningFound := strings.Contains(cOutput, "Warning: unpinned requirement ")
	if warningFound {
		tmpcOutput := StringToLastLine(cOutput)
		cOutput = tmpcOutput
	}

	// step 0.1: error cloning repository!
	if strings.Contains(cOutput, "ERROR_CLONING") {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": errorOutput,
				"containers.$.cResult": "failed",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
		}
		return
	}

	// step 1: Unmarshall cOutput into safetyOutput struct.
	safetyOutput := SafetyOutput{}
	err := json.Unmarshal([]byte(cOutput), &safetyOutput)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 1018, cOutput, err)
		return
	}

	// step 1.1: Sets the container output to "No issues found" if SafetyIssues returns an empty slice
	if len(safetyOutput.SafetyIssues) == 0 && !warningFound {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": "No issues found.",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
		}
		return
	}

	// step 2: finds Vulnerabilities
	cResult = "passed"
	if len(safetyOutput.SafetyIssues) != 0 {
		cResult = "failed"
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
		},
	}
	err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
		return
	}
}

// StringToLastLine receives a string with multiple lines and returns it's last
func StringToLastLine(s string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines[len(lines)-1]
}

// GetAllLinesButLast receives a string with multiple lines and returns all but the last line.
func GetAllLinesButLast(s string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	lines = lines[:len(lines)-1]
	return lines
}

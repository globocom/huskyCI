// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyci/log"
	"gopkg.in/mgo.v2/bson"
)

//SafetyOutput is the struct that holds issues, messages and errors found on a Safety scan.
type SafetyOutput struct {
	SafetyIssues []SafetyIssue
}

//SafetyIssue is a struct that holds the results that were scanned and the file they came from.
type SafetyIssue struct {
	Dependency string
	Below      string
	Version    string
	Comment    string
	ID         int
}

//SafetyStartAnalysis analyses the output from Safety and sets cResult based on it.
func SafetyStartAnalysis(CID string, cOutput string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

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
		log.Error("SafetyStartAnalysis", "SAFETY", 1014, cOutput, err)
		return
	}

	// step 1.1: Sets the container output to "No issues found" if SafetyIssues returns an empty slice
	if len(safetyOutput.SafetyIssues) == 0 {
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

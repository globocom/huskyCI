// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"gopkg.in/mgo.v2/bson"
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

// BrakemanStartAnalysis analyses the output from Brakeman and sets a cResult based on it.
func BrakemanStartAnalysis(CID string, cOutput string, RID string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.1: nil cOutput states that no Issues were found.
	if cOutput == "" {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "passed",
				"containers.$.cInfo":   "No issues found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("BrakemanStartAnalysis", "BRAKEMAN", 2007, "Step 0.1 ", err)
		}
		return
	}

	// step 0.2: error cloning repository!
	if strings.Contains(cOutput, "ERROR_CLONING") {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   errorOutput,
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("BrakemanStartAnalysis", "BRAKEMAN", 2007, "Step 0.2 ", err)
		}
		return
	}

	// step 1: Unmarshall cOutput into BrakemanOutput struct.
	brakemanOutput := BrakemanOutput{}
	err := json.Unmarshal([]byte(cOutput), &brakemanOutput)
	if err != nil {
		log.Error("BrakemanStartAnalysis", "BRAKEMAN", 1005, cOutput, err)
		return
	}

	// step 1.1: An empty errors slice means no vulnerabilities were found
	if len(brakemanOutput.Warnings) == 0 {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "passed",
				"containers.$.cInfo":   "No issues found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("BrakemanStartAnalysis", "BRAKEMAN", 2007, "Step 1.1 ", err)
		}
		return
	}

	// step 2: find Issues that have confidence "High" or "Medium".
	cResult = "passed"
	for _, warning := range brakemanOutput.Warnings {
		if warning.Confidence == "High" || warning.Confidence == "Medium" {
			cResult = "failed"
			break
		}
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	issueMessage := "No issues found."
	if cResult != "passed" {
		issueMessage = "Issues found."
	}
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
			"containers.$.cInfo":   issueMessage,
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("BrakemanStartAnalysis", "BRAKEMAN", 2007, "Step 3 ", err)
		return
	}

	// step 4: get updated analysis based on its RID
	analysisQuery = map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("BrakemanStartAnalysis", "BRAKEMAN", 2008, CID, err)
		return
	}

	// step 5: finally, update analysis with huskyCI results
	analysis.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput = prepareHuskyCIBrakemanResults(brakemanOutput)
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("BrakemanStartAnalysis", "BRAKEMAN", 2007, err)
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

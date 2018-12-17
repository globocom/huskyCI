// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/glbgelf"
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
func BrakemanStartAnalysis(CID string, cOutput string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.2: error cloning repository!
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
			if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "BrakemanStartAnalysis",
				"info":   "BRAKEMAN"}, "ERROR", "Error updating AnalysisCollection (inside brakeman.go):", err); errLog != nil {
				fmt.Println("glbgelf error: ", errLog)
			}
		}
		return
	}

	// step 1: Unmarshall cOutput into BrakemanOutput struct.
	brakemanOutput := BrakemanOutput{}
	err := json.Unmarshal([]byte(cOutput), &brakemanOutput)
	if err != nil {
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "BrakemanStartAnalysis",
			"info":   "BRAKEMAN"}, "ERROR", "Unmarshall error (brakeman.go):", err); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		return
	}

	// step 1.1: An empty errors slice means no vulnerabilities were found
	if len(brakemanOutput.Warnings) == 0 {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": "No issues found.",
				"containers.$.cResult": "passed",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "BrakemanStartAnalysis",
				"info":   "BRAKEMAN"}, "ERROR", "Error updating AnalysisCollection (inside brakeman.go):", err); errLog != nil {
				fmt.Println("glbgelf error: ", errLog)
			}
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
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
		},
	}
	err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "BrakemanStartAnalysis",
			"info":   "BRAKEMAN"}, "ERROR", "Error updating AnalysisCollection (inside brakeman.go):", err); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		return
	}
}

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

// BanditOutput is the structs that holds the json output form bandit analysis.
type BanditOutput struct {
	Errors  json.RawMessage `json:"errors"`
	Results []Result        `json:"results"`
}

// Result is the struct that holds detailed information of issues found in bandit analysis.
type Result struct {
	Code            string `json:"code"`
	Filename        string `json:"filename"`
	IssueConfidence string `json:"issue_confidence"`
	IssueSeverity   string `json:"issue_severity"`
	IssueText       string `json:"issue_text"`
	LineNumber      int    `json:"line_number"`
	LineRange       []int  `json:"line_range"`
	TestID          string `json:"test_id"`
	TestName        string `json:"test_name"`
}

// BanditStartAnalysis analyses the output from Bandit and sets a cResult based on it.
func BanditStartAnalysis(CID string, cOutput string) {

	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// error cloning repository!
	if strings.Contains(cOutput, "ERROR_CLONING") {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": errorOutput,
				"containers.$.cResult": "error",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("BanditStartAnalysis", "BANDIT", 2007, "Step 1", err)
		}
		return
	}

	var banditResult BanditOutput
	if err := json.Unmarshal([]byte(cOutput), &banditResult); err != nil {
		log.Error("BanditStartAnalysis", "BANDIT", 1006, cOutput, err)
		return
	}

	// Sets the container output to "No issues found" if banditResult returns an empty slice
	if len(banditResult.Results) == 0 {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": "No issues found.",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "BanditStartAnalysis",
				"info":   "BANDIT"}, "ERROR", "Error updating AnalysisCollection (inside bandit.go):", err); errLog != nil {
				fmt.Println("glbgelf error: ", errLog)
			}
		}
		return
	}

	// verify if there was any error in the analysis.
	if banditResult.Errors != nil {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("BanditStartAnalysis", "BANDIT", 2007, "Step 2", err)
		}
	}

	// find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult := "passed"
	for _, issue := range banditResult.Results {
		if (issue.IssueSeverity == "HIGH" || issue.IssueSeverity == "MEDIUM") && issue.IssueConfidence == "HIGH" {
			cResult = "failed"
			break
		}
	}

	// update the status of analysis.
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": cResult,
		},
	}
	if err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery); err != nil {
		log.Error("BanditStartAnalysis", "BANDIT", 2007, "Step 3", err)
		return
	}
}

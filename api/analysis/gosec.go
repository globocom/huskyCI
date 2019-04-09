// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"gopkg.in/mgo.v2/bson"
)

// GosecOutput is the struct that holds issues and stats found on a Gosec scan.
type GosecOutput struct {
	GosecIssues []GosecIssue
	GosecStats  GosecStats
}

// GosecIssue is the struct that holds all detailed information of a vulnerability found.
type GosecIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// GosecStats is the struct that holds the stats found on a Gosec scan.
type GosecStats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	NoSec int `json:"nosec"`
	Found int `json:"found"`
}

// GosecStartAnalysis analyses the output from Gosec and sets a cResult based on it.
func GosecStartAnalysis(CID string, cOutput string) {

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
			log.Error("GosecStartAnalysis", "GOSEC", 2007, "Step 0.1 ", err)
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
			log.Error("GosecStartAnalysis", "GOSEC", 2007, "Step 0.2 ", err)
		}
		return
	}

	// step 1: Unmarshall cOutput into GosecOutput struct.
	gosecOutput := GosecOutput{}
	err := json.Unmarshal([]byte(cOutput), &gosecOutput)
	if err != nil {
		log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return
	}

	// step 2: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult = "warning"
	for _, issue := range gosecOutput.GosecIssues {
		if (issue.Severity == "HIGH" || issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			cResult = "failed"
			break
		}
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	issueMessage := "Warning found."
	if cResult != "warning" {
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
		log.Error("GosecStartAnalysis", "GOSEC", 2007, "Step 3 ", err)
		return
	}
}

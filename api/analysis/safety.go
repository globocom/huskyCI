// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/util"
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

	var outputJSON string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	errorSafety := strings.Contains(cOutput, "ERROR_RUNNING_SAFETY")
	if errorSafety {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "warning", // will not fail CI now
				"containers.$.cInfo":   "Internal error running Safety.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Warning("SafetyStartAnalysis", "SAFETY", 2007, err)
		}
		return
	}

	requirementsNotFound := strings.Contains(cOutput, "ERROR_REQ_NOT_FOUND")
	if requirementsNotFound {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "warning", // will not fail CI now
				"containers.$.cInfo":   "Requirements not found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Warning("SafetyStartAnalysis", "SAFETY", 2007, err)
		}
		return
	}

	errorCloning := strings.Contains(cOutput, "ERROR_CLONING")
	if errorCloning {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   "Error clonning repository.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
		}
		return
	}

	warningFound := strings.Contains(cOutput, "Warning: unpinned requirement ")
	if warningFound {
		outputJSON = util.GetLastLine(cOutput)
		cOutput = outputJSON
	}

	cOutputSanitized := util.SanitizeSafetyJSON(cOutput)

	safetyOutput := SafetyOutput{}
	err := json.Unmarshal([]byte(cOutputSanitized), &safetyOutput)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 1018, cOutput, err)
		return
	}

	if len(safetyOutput.SafetyIssues) == 0 {

		if warningFound {
			updateContainerAnalysisQuery := bson.M{
				"$set": bson.M{
					"containers.$.cResult": "warning",
					"containers.$.cInfo":   "Warning found",
				},
			}
			err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
			if err != nil {
				log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
			}
			return
		}

		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "passed",
				"containers.$.cInfo":   "No issues found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
		}
		return
	}

	// Issues found. client will have to handle with warnings and issues.
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cResult": "failed",
			"containers.$.cInfo":   "Issues found.",
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
	}
	return

}

// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
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
func SafetyStartAnalysis(CID string, cOutput string, RID string) {

	var outputJSON string
	var outputWarnings []string

	analysisQuery := map[string]interface{}{"containers.CID": CID}

	reqNotFound := strings.Contains(cOutput, "ERROR_REQ_NOT_FOUND")
	failedRunning := strings.Contains(cOutput, "ERROR_RUNNING_SAFETY")
	warningFound := strings.Contains(cOutput, "Warning: unpinned requirement ")
	errorCloning := strings.Contains(cOutput, "ERROR_CLONING")

	safetyOutput := SafetyOutput{}

	if failedRunning {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "warning", // will not fail CI now
				"containers.$.cInfo":   "Internal error running Safety.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Warning("SafetyStartAnalysis", "SAFETY", 2007, err)
			return
		}

		if err := updateSafetyAnalysisWithResults(RID, CID, safetyOutput, failedRunning, reqNotFound, warningFound, outputWarnings); err != nil {
			return
		}

		return
	}

	if reqNotFound {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "warning", // will not fail CI now
				"containers.$.cInfo":   "Requirements not found.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Warning("SafetyStartAnalysis", "SAFETY", 2007, err)
			return
		}

		if err := updateSafetyAnalysisWithResults(RID, CID, safetyOutput, failedRunning, reqNotFound, warningFound, outputWarnings); err != nil {
			return
		}

		return
	}

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
			return
		}
		return
	}

	if warningFound {
		outputJSON = util.GetLastLine(cOutput)
		outputWarnings = util.GetAllLinesButLast(cOutput)
		cOutput = outputJSON
	}

	cOutputSanitized := util.SanitizeSafetyJSON(cOutput)

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
				return
			}

			if err := updateSafetyAnalysisWithResults(RID, CID, safetyOutput, failedRunning, reqNotFound, warningFound, outputWarnings); err != nil {
				return
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

		if err := updateSafetyAnalysisWithResults(RID, CID, safetyOutput, failedRunning, reqNotFound, warningFound, outputWarnings); err != nil {
			return
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

	if err := updateSafetyAnalysisWithResults(RID, CID, safetyOutput, failedRunning, reqNotFound, warningFound, outputWarnings); err != nil {
		return
	}

}

func updateSafetyAnalysisWithResults(RID, CID string, safetyOutput SafetyOutput, failedRunning, reqNotFound, warningFound bool, outputWarnings []string) error {

	// get updated analysis based on its RID
	analysisQuery := map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 2008, CID, err)
		return err
	}

	// update analysis with huskyCI results
	analysis.HuskyCIResults.PythonResults.HuskyCISafetyOutput = prepareHuskyCISafetyOutput(safetyOutput, failedRunning, reqNotFound, warningFound, outputWarnings)
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("SafetyStartAnalysis", "SAFETY", 2007, err)
		return err
	}

	return nil
}

// prepareHuskyCISafetyOutput will prepare Safety output to be added into PythonResults struct
func prepareHuskyCISafetyOutput(safetyOutput SafetyOutput, failedRunning, reqNotFound, warningFound bool, outputWarnings []string) types.HuskyCISafetyOutput {

	var huskyCIsafetyResults types.HuskyCISafetyOutput
	var onlyWarning bool

	if failedRunning {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
		safetyVuln.Details = "Internal error running Safety."

		huskyCIsafetyResults.LowVulnsSafety = append(huskyCIsafetyResults.LowVulnsSafety, safetyVuln)

		return huskyCIsafetyResults
	}

	if reqNotFound {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
		safetyVuln.Details = "requirements.txt not found"

		huskyCIsafetyResults.LowVulnsSafety = append(huskyCIsafetyResults.LowVulnsSafety, safetyVuln)

		return huskyCIsafetyResults
	}

	if warningFound {

		if len(safetyOutput.SafetyIssues) == 0 {
			onlyWarning = true
		}

		for _, warning := range outputWarnings {
			safetyVuln := types.HuskyCIVulnerability{}
			safetyVuln.Language = "Python"
			safetyVuln.SecurityTool = "Safety"
			safetyVuln.Severity = "warning"
			safetyVuln.Details = util.AdjustWarningMessage(warning)

			huskyCIsafetyResults.LowVulnsSafety = append(huskyCIsafetyResults.LowVulnsSafety, safetyVuln)

		}
		if onlyWarning {
			return huskyCIsafetyResults
		}
	}

	for _, issue := range safetyOutput.SafetyIssues {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "high"
		safetyVuln.Details = issue.Comment
		safetyVuln.Code = issue.Dependency + " " + issue.Version
		safetyVuln.VunerableBelow = issue.Below

		huskyCIsafetyResults.HighVulnsSafety = append(huskyCIsafetyResults.HighVulnsSafety, safetyVuln)
	}

	return huskyCIsafetyResults
}

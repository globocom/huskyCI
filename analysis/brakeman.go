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
	Errors []ErrorItem `json:"errors"`
}

// ErrorItem is the struct that holds all detailed information of a vulnerability found.
type ErrorItem struct {
	Error    string `json:"error"`
	Location string `json:"location"`
}

// BrakemanStartAnalysis analyses the output from Brakeman and sets a cResult based on it.
func BrakemanStartAnalysis(CID string, cOutput string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.1: nil cOutput states that no Issues were found.
	if cOutput == "" {
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

	// step 2: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult = "passed"
	if len(brakemanOutput.Errors) > 0 {
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
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "BrakemanStartAnalysis",
			"info":   "BRAKEMAN"}, "ERROR", "Error updating AnalysisCollection (inside brakeman.go):", err); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		return
	}
}

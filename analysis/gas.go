package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// GasOutput is the struct that holds issues and stats found on a Gas scan.
type GasOutput struct {
	Issues []Issue
	Stats  Stats
}

// Issue is the struct that holds all detailed information of a vulnerability found.
type Issue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// Stats is the struct that holds the stats found on a Gas scan.
type Stats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	NoSec int `json:"nosec"`
	Found int `json:"found"`
}

// GasStartAnalysis analyses the output from Gas and sets a cResult based on it.
func GasStartAnalysis(CID string, cOutput string) {

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
			fmt.Println("Error updating AnalysisCollection (inside gas.go):", err)
		}
		return
	}

	// step 0.2: error cloning repository!
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
			fmt.Println("Error updating AnalysisCollection (inside gas.go):", err)
		}
		return
	}

	// step 1: Unmarshall cOutput into GasOutput struct.
	gasOutput := GasOutput{}
	err := json.Unmarshal([]byte(cOutput), &gasOutput)
	if err != nil {
		fmt.Println("Unmarshall error (gas.go):", err)
		fmt.Println(cOutput)
		return
	}

	// step 2: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	cResult = "passed"
	for _, issue := range gasOutput.Issues {
		if (issue.Severity == "HIGH" || issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
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
		fmt.Println("Error updating AnalysisCollection (inside gas.go):", err)
		return
	}
}

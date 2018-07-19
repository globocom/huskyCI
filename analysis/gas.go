package analysis

import (
	"encoding/json"
	"fmt"
	"regexp"
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

// GasStartAnalysis checks
func GasStartAnalysis(CID string, cOutput string, RID string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0: nil cOutput states that no Issues were found.
	if cOutput == "" {
		// what if some error occurred? use $? to check this?
		cResult = "passed"
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cOutput": "No issues found.",
				"containers.$.cResult": cResult,
			},
		}
		err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			fmt.Println("Error updating AnalysisCollection (inside gas.go):", err)
		}
		return
	}

	// step 1: clean cOutput (removing strange chars from docker logs API and control characters)
	reg, err := regexp.Compile(`\{.*\:\{.*\:.*\}\}`)
	if err != nil {
		fmt.Println("Error regexp:", err)
	}
	cleanedOutput := reg.FindString(cOutput)
	cleanedOutput = stripCtlAndExtFromBytes(cleanedOutput)

	// step 2: Unmarshall cleanedOutput
	gasOutput := GasOutput{}
	err = json.Unmarshal([]byte(cleanedOutput), &gasOutput)
	if err != nil {
		fmt.Println("Unmarshall error:", err)
		fmt.Println(cleanedOutput)
		return
	}

	// step 3: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH"
	cResult = "passed"
	for _, issue := range gasOutput.Issues {
		if (issue.Severity == "HIGH" || issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			cResult = "failed"
		}
	}

	// step 4: update cResult and cOutput
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cOutput": cleanedOutput,
			"containers.$.cResult": cResult,
		},
	}
	err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		fmt.Println("Error updating AnalysisCollection (inside gas.go):", err)
	}
}

// stripCtlAndExtFromBytes removes control chars from a specific string
func stripCtlAndExtFromBytes(str string) string {
	// the following chars was also breaking Unmarshall.
	str = strings.Replace(str, `\ `, `\\`, -1)
	str = strings.Replace(str, `@`, ``, -1)
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

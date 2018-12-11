package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

//RetirejsOutput is the struct that holds issues, messages and errors found on a Retire scan.
type RetirejsOutput struct {
	RetirejsIssues []RetirejsIssue `json:"data"`
	Messages       json.RawMessage `json:"messages"`
	Errors         json.RawMessage `json:"errors"`
}

//RetirejsIssue is a struct that holds the results that were scanned and the file they came from.
type RetirejsIssue struct {
	File            string           `json:"file"`
	RetirejsResults []RetirejsResult `json:"results"`
}

//RetirejsResult is a struct that holds the vulnerabilities found on a component being used by the code being analysed.
type RetirejsResult struct {
	Version                 string                  `json:"version"`
	Component               string                  `json:"component"`
	Detection               string                  `json:"detection"`
	RetirejsVulnerabilities []RetirejsVulnerability `json:"vulnerabilities"`
}

//RetirejsVulnerability is a struct that holds info on what vulnerabilies were found.
type RetirejsVulnerability struct {
	Info                string               `json:"info"`
	Below               string               `json:"below"`
	Severity            string               `json:"severity"`
	RetirejsIdentifiers []RetirejsIdentifier `json:"identifiers"`
}

//RetirejsIdentifier is a struct that holds details on the vulnerabilities found.
type RetirejsIdentifier struct {
	IssueFound string   `json:"issue"`
	Summary    string   `json:"summary"`
	CVE        []string `json:"CVE"`
}

//RetirejsStartAnalysis analyses the output from RetireJS and sets cResult basdes on it.
func RetirejsStartAnalysis(CID string, cOutput string) {

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
			fmt.Println("Error updating AnalysisCollection (inside retirejs.go):", err)
		}
		return
	}

	// step 1: Unmarshall cOutput into RetireOutput struct.
	retirejsOutput := RetirejsOutput{}
	err := json.Unmarshal([]byte(cOutput), &retirejsOutput)
	if err != nil {
		fmt.Println("Unmarshall error (retirejs.go):", err)
		fmt.Println(cOutput)
		return
	}

	// step 2: find Vulnerabilities that have severity "medium" or "high".
	cResult = "passed"
	for _, issue := range retirejsOutput.RetirejsIssues {
		for _, result := range issue.RetirejsResults {
			for _, vulnerability := range result.RetirejsVulnerabilities {
				if vulnerability.Severity == "high" || vulnerability.Severity == "medium" {
					cResult = "failed"
					break
				}
			}
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
		fmt.Println("Error updating AnalysisCollection (inside retirejs.go):", err)
		return
	}
}

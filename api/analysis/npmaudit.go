// Copyright 2019 Globo.com authors. All rights reserved.
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

type NpmAuditOutput struct {
	Advisories map[string]Vulnerability `json:"advisories"`
	Metadata   Metadata                 `json:"metadata"`
}

type Vulnerability struct {
	Findings           []Finding `json:"findings"`
	ID                 int       `json:"id"`
	ModuleName         string    `json:"module_name"`
	VulnerableVersions string    `json:"vulnerable_versions"`
	Severity           string    `json:"severity"`
	Overview           string    `json:"overview"`
}

type Finding struct {
	Version string `json:"version"`
}

type Metadata struct {
	Vulnerabilities VulnerabilitiesSummary `json:"vulnerabilities"`
}

type VulnerabilitiesSummary struct {
	Info     int `json:"info"`
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

func NpmAuditStartAnalysis(CID string, cOutput string) {

	var cResult string
	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// step 0.1: error cloning repository!
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
			log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, "Step 0.1 ", err)
		}
		return
	}

	// step 0.2: repository doesn't have package.json
	if strings.Contains(cOutput, "ERROR_RUNNING_NPMAUDIT") {
		errorOutput := fmt.Sprintf("Container error: %s", cOutput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   errorOutput,
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, "Step 0.2 ", err)
		}
		return
	}

	// step 1: Unmarshall cOutput into NpmAuditOutput struct.
	npmAuditOutput := NpmAuditOutput{}
	err := json.Unmarshal([]byte(cOutput), &npmAuditOutput)
	if err != nil {
		log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 1022, cOutput, err)
		return
	}

	// step 2: find Issues that have severity "moderate" or "high.
	cResult = "passed"
	for _, vulnerability := range npmAuditOutput.Advisories {
		if vulnerability.Severity == "high" || vulnerability.Severity == "moderate" {
			cResult = "failed"
			break
		}
	}

	// step 3: update analysis' cResult into AnalyisCollection.
	issueMessage := "No issues found."
	if cResult == "failed" {
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
		log.Error("NpmAuditStartAnalysis", "NPMAUDIT", 2007, "Step 3 ", err)
		return
	}
}

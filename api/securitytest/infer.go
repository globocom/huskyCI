// Copyright 2021 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// InferOutput holds all data from Infer output.
type InferOutput []InferResult

// InferResult holds the detailed information from Infer output.
type InferResult struct {
	Type     string `json:"bug_type"`
	Message  string `json:"qualifier"`
	File     string `json:"file"`
	Line     string `json:"line"`
	Severity string `json:"severity"`
	Title    string `json:"bug_type_hum"`
}

func analyzeInfer(scanInfo *SecTestScanInfo) error {
	var inferOutput InferOutput

	if err := json.Unmarshal([]byte(scanInfo.Container.COutput), &inferOutput); err != nil {
		log.Error("analyzeInfer", "INFER", 1041, scanInfo.Container.COutput, err)
		scanInfo.ErrorFound = err
		return err
	}
	scanInfo.FinalOutput = inferOutput

	// if len is equal to zero no issues were found
	if len(inferOutput) == 0 {
		scanInfo.prepareContainerAfterScan()
		return nil
	}

	scanInfo.prepareInferVulns()
	scanInfo.prepareContainerAfterScan()
	return nil
}

func (inferScan *SecTestScanInfo) prepareInferVulns() {
	huskyCIInferResults := types.HuskyCISecurityTestOutput{}
	inferOutput := inferScan.FinalOutput.(InferOutput)

	for _, result := range inferOutput {
		inferVuln := types.HuskyCIVulnerability{
			Language:     "Java",
			SecurityTool: "Infer",
			Severity:     result.Severity,
			File:         result.File,
			Line:         result.Line,
			Details:      result.Message,
			Type:         result.Type,
			Title:        result.Title,
		}

		switch inferVuln.Severity {
		case "INFO":
			huskyCIInferResults.LowVulns = append(huskyCIInferResults.LowVulns, inferVuln)
		case "WARNING":
			huskyCIInferResults.MediumVulns = append(huskyCIInferResults.MediumVulns, inferVuln)
		case "ERROR":
			huskyCIInferResults.HighVulns = append(huskyCIInferResults.HighVulns, inferVuln)
		}
	}
}

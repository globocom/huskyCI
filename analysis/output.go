// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/globocom/husky-client/types"
)

// CheckContainerOutput will validate the output of a given container.
func CheckContainerOutput(container types.Container) {

	switch container.SecurityTest.Name {
	case "enry":

	case "gas":
		PrintGasOutput(container.COutput)
	default:
		fmt.Println("[HUSKYCI][ERROR] securityTest name not recognized:", container.SecurityTest.Name)
		os.Exit(1)
	}
}

// PrintGasOutput will print the output of Gas.
func PrintGasOutput(containerOutput string) {

	if containerOutput == "No issues found." {
		color.Green("[HUSKYCI][*] :) ")
		os.Exit(0)
	}

	foundVuln := false
	foundInfo := false
	gasOutput := types.GasOutput{}
	err := json.Unmarshal([]byte(containerOutput), &gasOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal gasOutput!", containerOutput)
		os.Exit(1)
	}

	for _, issue := range gasOutput.Issues {
		if (issue.Severity == "HIGH") && (issue.Confidence == "HIGH") {
			foundVuln = true
			fmt.Println("[HUSKYCI][!]")
			color.Red("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Red("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Red("[HUSKYCI][!] Details: %s", issue.Details)
			color.Red("[HUSKYCI][!] File: %s", issue.File)
			color.Red("[HUSKYCI][!] Line: %s", issue.Line)
			color.Red("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range gasOutput.Issues {
		if (issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			foundVuln = true
			fmt.Println("[HUSKYCI][!]")
			color.Yellow("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Yellow("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.Details)
			color.Yellow("[HUSKYCI][!] File: %s", issue.File)
			color.Yellow("[HUSKYCI][!] Line: %s", issue.Line)
			color.Yellow("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range gasOutput.Issues {
		if issue.Severity == "LOW" {
			foundInfo = true
			fmt.Println("[HUSKYCI][!]")
			color.Blue("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Blue("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Blue("[HUSKYCI][!] Details: %s", issue.Details)
			color.Blue("[HUSKYCI][!] File: %s", issue.File)
			color.Blue("[HUSKYCI][!] Line: %s", issue.Line)
			color.Blue("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	if foundVuln {
		color.Red("[HUSKYCI][X] :( ")
		os.Exit(1)
	}
	if foundInfo {
		fmt.Println("[HUSKYCI][*] :|")
		os.Exit(0)
	}
	color.Green("[HUSKYCI][*] :) ")
}

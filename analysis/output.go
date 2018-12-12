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
	case "gosec":
		PrintGosecOutput(container.COutput)
	case "bandit":
		PrintBanditOutput(container.COutput)
	default:
		fmt.Println("[HUSKYCI][ERROR] securityTest name not recognized:", container.SecurityTest.Name)
		os.Exit(1)
	}
}

// PrintGosecOutput will print the Gosec output.
func PrintGosecOutput(containerOutput string) {

	if containerOutput == "No issues found." {
		color.Green("[HUSKYCI][*] Gosec :)\n\n")
		return
	}

	foundVuln := false
	foundInfo := false
	gosecOutput := types.GosecOutput{}
	err := json.Unmarshal([]byte(containerOutput), &gosecOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal gosecOutput!", containerOutput)
		os.Exit(1)
	}

	for _, issue := range gosecOutput.GosecIssues {
		if (issue.Severity == "HIGH") && (issue.Confidence == "HIGH") {
			foundVuln = true
			color.Red("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Red("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Red("[HUSKYCI][!] Details: %s", issue.Details)
			color.Red("[HUSKYCI][!] File: %s", issue.File)
			color.Red("[HUSKYCI][!] Line: %d", issue.Line)
			color.Red("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range gosecOutput.GosecIssues {
		if (issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			foundVuln = true
			color.Yellow("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Yellow("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.Details)
			color.Yellow("[HUSKYCI][!] File: %s", issue.File)
			color.Yellow("[HUSKYCI][!] Line: %d", issue.Line)
			color.Yellow("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range gosecOutput.GosecIssues {
		if issue.Severity == "LOW" {
			foundInfo = true
			color.Blue("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Blue("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Blue("[HUSKYCI][!] Details: %s", issue.Details)
			color.Blue("[HUSKYCI][!] File: %s", issue.File)
			color.Blue("[HUSKYCI][!] Line: %d", issue.Line)
			color.Blue("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	if foundVuln {
		color.Red("[HUSKYCI][X] :(\n\n")
		types.FoundVuln = true
	} else if foundInfo {
		fmt.Printf("[HUSKYCI][*] Gosec :|\n\n")
	} else {
		color.Green("[HUSKYCI][*] Gosec :)\n\n")
	}

}

// PrintBanditOutput will print Bandit output.
func PrintBanditOutput(containerOutput string) {

	if containerOutput == "No issues found." {
		color.Green("[HUSKYCI][*] Bandit :)\n\n")
		return
	}

	foundVuln := false
	foundInfo := false
	banditOutput := types.BanditOutput{}
	err := json.Unmarshal([]byte(containerOutput), &banditOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal banditOutput!", containerOutput)
		os.Exit(1)
	}

	for _, issue := range banditOutput.Results {
		if (issue.IssueSeverity == "HIGH") && (issue.IssueConfidence == "HIGH") {
			foundVuln = true
			color.Red("[HUSKYCI][!] Severity: %s", issue.IssueSeverity)
			color.Red("[HUSKYCI][!] Confidence: %s", issue.IssueConfidence)
			color.Red("[HUSKYCI][!] Details: %s", issue.IssueText)
			color.Red("[HUSKYCI][!] File: %s", issue.Filename)
			color.Red("[HUSKYCI][!] Line: %d", issue.LineNumber)
			color.Red("[HUSKYCI][!] Code:\n%s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range banditOutput.Results {
		if (issue.IssueSeverity == "MEDIUM") && (issue.IssueConfidence == "HIGH") {
			foundVuln = true
			color.Yellow("[HUSKYCI][!] Severity: %s", issue.IssueSeverity)
			color.Yellow("[HUSKYCI][!] Confidence: %s", issue.IssueConfidence)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.IssueText)
			color.Yellow("[HUSKYCI][!] File: %s", issue.Filename)
			color.Yellow("[HUSKYCI][!] Line: %d", issue.LineNumber)
			color.Yellow("[HUSKYCI][!] Code:\n%s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range banditOutput.Results {
		if issue.IssueSeverity == "LOW" {
			foundInfo = true
			color.Blue("[HUSKYCI][!] Severity: %s", issue.IssueSeverity)
			color.Blue("[HUSKYCI][!] Confidence: %s", issue.IssueConfidence)
			color.Blue("[HUSKYCI][!] Details: %s", issue.IssueText)
			color.Blue("[HUSKYCI][!] File: %s", issue.Filename)
			color.Blue("[HUSKYCI][!] Line: %d", issue.LineNumber)
			color.Blue("[HUSKYCI][!] Code:\n%s", issue.Code)
			fmt.Println()
		}
	}

	if foundVuln {
		color.Red("[HUSKYCI][X] Bandit :(\n\n")
		types.FoundVuln = true
	} else if foundInfo {
		fmt.Printf("[HUSKYCI][*] Bandit :|\n\n")
	} else {
		color.Green("[HUSKYCI][*] Bandit :)\n\n")
	}

}

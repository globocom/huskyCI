// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/globocom/huskyCI/client/types"
	"github.com/globocom/huskyCI/client/util"
)

// CheckMongoDBContainerOutput will validate the output of a given container.
func CheckMongoDBContainerOutput(container types.Container) {

	switch container.SecurityTest.Name {
	case "enry":
	case "gosec":
		PrintGosecOutput(container.COutput, container.CInfo)
	case "bandit":
		PrintBanditOutput(container.COutput, container.CInfo)
	case "retirejs":
		PrintRetirejsOutput(container.COutput, container.CInfo)
	case "brakeman":
		PrintBrakemanOutput(container.COutput, container.CInfo)
	case "safety":
		PrintSafetyOutput(container.COutput, container.CInfo)
	default:
		fmt.Println("[HUSKYCI][ERROR] securityTest name not recognized:", container.SecurityTest.Name)
		os.Exit(1)
	}
}

// PrintGosecOutput will print the Gosec output.
func PrintGosecOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		color.Green("[HUSKYCI][*] Gosec :)\n\n")
		return
	}

	foundVuln := false
	foundInfo := false
	gosecOutput := types.GosecOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &gosecOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal gosecOutput!", mongoDBcontainerOutput)
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
		color.Red("[HUSKYCI][X] Gosec :(\n\n")
		types.FoundVuln = true
	} else if foundInfo {
		fmt.Printf("[HUSKYCI][*] Gosec :|\n\n")
	} else {
		color.Green("[HUSKYCI][*] Gosec :)\n\n")
	}

}

// PrintBanditOutput will print Bandit output.
func PrintBanditOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		color.Green("[HUSKYCI][*] Bandit :)\n\n")
		return
	}

	foundVuln := false
	foundInfo := false
	banditOutput := types.BanditOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &banditOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal banditOutput!", mongoDBcontainerOutput)
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

// PrintRetirejsOutput will print Retirejs output.
func PrintRetirejsOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		color.Green("[HUSKYCI][*] RetireJS :)\n\n")
		return
	}

	foundVuln := false
	foundInfo := false
	retirejsOutput := types.RetirejsOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &retirejsOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal retirejsOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, issue := range retirejsOutput.RetirejsIssues {
		for _, result := range issue.RetirejsResults {
			for _, vulnerability := range result.RetirejsVulnerabilities {
				if vulnerability.Severity == "high" {
					foundVuln = true
					color.Red("[HUSKYCI][!] Severity: %s", vulnerability.Severity)
					color.Red("[HUSKYCI][!] Details: %s", vulnerability.Info)
					color.Red("[HUSKYCI][!] File: %s", issue.File)
					color.Red("[HUSKYCI][!] Component: %s", result.Component)
					color.Red("[HUSKYCI][!] Version: %s", result.Version)
					color.Red("[HUSKYCI][!] Vulnerable Below: %s", vulnerability.Below)
					fmt.Println()
				}
			}
		}
	}

	for _, issue := range retirejsOutput.RetirejsIssues {
		for _, result := range issue.RetirejsResults {
			for _, vulnerability := range result.RetirejsVulnerabilities {
				if vulnerability.Severity == "medium" {
					foundVuln = true
					color.Yellow("[HUSKYCI][!] Severity: %s", vulnerability.Severity)
					color.Yellow("[HUSKYCI][!] Details: %s", vulnerability.Info)
					color.Yellow("[HUSKYCI][!] File: %s", issue.File)
					color.Yellow("[HUSKYCI][!] Component: %s", result.Component)
					color.Yellow("[HUSKYCI][!] Version: %s", result.Version)
					color.Yellow("[HUSKYCI][!] Vulnerable Below: %s", vulnerability.Below)
					fmt.Println()
				}
			}
		}
	}

	for _, issue := range retirejsOutput.RetirejsIssues {
		for _, result := range issue.RetirejsResults {
			for _, vulnerability := range result.RetirejsVulnerabilities {
				if vulnerability.Severity == "low" {
					foundInfo = true
					color.Blue("[HUSKYCI][!] Severity: %s", vulnerability.Severity)
					color.Blue("[HUSKYCI][!] Details: %s", vulnerability.Info)
					color.Blue("[HUSKYCI][!] File: %s", issue.File)
					color.Blue("[HUSKYCI][!] Component: %s", result.Component)
					color.Blue("[HUSKYCI][!] Version: %s", result.Version)
					color.Blue("[HUSKYCI][!] Vulnerable Below: %s", vulnerability.Below)
					fmt.Println()
				}
			}
		}
	}

	if foundVuln {
		color.Red("[HUSKYCI][X] RetireJS :(\n\n")
		types.FoundVuln = true
	} else if foundInfo {
		fmt.Printf("[HUSKYCI][*] RetireJS :|\n\n")
	} else {
		color.Green("[HUSKYCI][*] RetireJS :)\n\n")
	}

}

// PrintBrakemanOutput will print Brakeman output.
func PrintBrakemanOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {
	if mongoDBcontainerInfo == "No issues found." {
		color.Green("[HUSKYCI][*] Brakeman :)\n\n")
		return
	}

	foundVuln := false
	foundInfo := false
	brakemanOutput := types.BrakemanOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &brakemanOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal brakemanOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, warning := range brakemanOutput.Warnings {
		if warning.Confidence == "High" {
			foundVuln = true
			color.Red("[HUSKYCI][!] Confidence: %s", warning.Confidence)
			color.Red("[HUSKYCI][!] Type: %s", warning.Type)
			color.Red("[HUSKYCI][!] Details: %s", warning.Details)
			color.Red("[HUSKYCI][!] Info: %s", warning.Message)
			color.Red("[HUSKYCI][!] File: %s", warning.File)
			color.Red("[HUSKYCI][!] line: %d", warning.Line)
			color.Red("[HUSKYCI][!] Code: %s", warning.Code)
			fmt.Println()
		}

		if warning.Confidence == "Medium" {
			foundVuln = true
			color.Yellow("[HUSKYCI][!] Confidence: %s", warning.Confidence)
			color.Yellow("[HUSKYCI][!] Type: %s", warning.Type)
			color.Yellow("[HUSKYCI][!] Details: %s", warning.Details)
			color.Yellow("[HUSKYCI][!] Info: %s", warning.Message)
			color.Yellow("[HUSKYCI][!] File: %s", warning.File)
			color.Yellow("[HUSKYCI][!] line: %d", warning.Line)
			color.Yellow("[HUSKYCI][!] Code: %s", warning.Code)
			fmt.Println()
		}

		if warning.Confidence == "Low" {
			foundInfo = true
			color.Blue("[HUSKYCI][!] Confidence: %s", warning.Confidence)
			color.Blue("[HUSKYCI][!] Type: %s", warning.Type)
			color.Blue("[HUSKYCI][!] Details: %s", warning.Details)
			color.Blue("[HUSKYCI][!] Info: %s", warning.Message)
			color.Blue("[HUSKYCI][!] File: %s", warning.File)
			color.Blue("[HUSKYCI][!] line: %d", warning.Line)
			color.Blue("[HUSKYCI][!] Code: %s", warning.Code)
			fmt.Println()
		}
	}

	if foundVuln {
		color.Red("[HUSKYCI][X] Brakeman :(\n\n")
		types.FoundVuln = true
	} else if foundInfo {
		fmt.Printf("[HUSKYCI][*] Brakeman :|\n\n")
	} else {
		color.Green("[HUSKYCI][*] Brakeman :)\n\n")
	}

}

// PrintSafetyOutput will print Safety output.
func PrintSafetyOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		color.Green("[HUSKYCI][*] Safety :)\n\n")
		return
	}

	if mongoDBcontainerInfo == "Requirements not found or this project uses latest dependencies." {
		fmt.Printf("[HUSKYCI][*] requirements.txt not found or this project uses latest dependencies.\n")
		fmt.Printf("[HUSKYCI][*] Safety :|\n\n")
		return
	}

	// Safety returns warnings and the json output in the same string, which need to be split
	var cOutputSanitized string
	safetyOutput := types.SafetyOutput{}
	warningFound := strings.Contains(mongoDBcontainerOutput, "Warning: unpinned requirement ")
	if !warningFound {
		// only issue found
		cOutputSanitized = util.SanitizeSafetyJSON(mongoDBcontainerOutput)
		err := json.Unmarshal([]byte(cOutputSanitized), &safetyOutput)
		if err != nil {
			fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal safetyOutput: ", err)
			os.Exit(1)
		}
	} else {
		// issues and warnings found
		onlyWarning := false
		outputJSON := util.GetLastLine(mongoDBcontainerOutput)
		outputWarnings := util.GetAllLinesButLast(mongoDBcontainerOutput)
		cOutputSanitized = util.SanitizeSafetyJSON(outputJSON)
		err := json.Unmarshal([]byte(cOutputSanitized), &safetyOutput)
		if err != nil {
			fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal safetyOutput: ", err)
			os.Exit(1)
		}
		if len(safetyOutput.SafetyIssues) == 0 {
			onlyWarning = true
		}
		for _, warning := range outputWarnings {
			color.Yellow("[HUSKYCI][!]: %s", warning)
		}
		if onlyWarning {
			fmt.Printf("[HUSKYCI][*] Safety :|\n\n")
			return
		}
	}

	for _, issue := range safetyOutput.SafetyIssues {
		color.Red("[HUSKYCI][!] Vulnerable Dependency: %s", issue.Dependency)
		color.Red("[HUSKYCI][!] Vulnerable Below: %s", issue.Below)
		color.Red("[HUSKYCI][!] Current Version: %s", issue.Version)
		color.Red("[HUSKYCI][!] Comment: %s", issue.Comment)
		fmt.Println()
	}

	color.Red("[HUSKYCI][X] Safety :(\n\n")
	types.FoundVuln = true
}

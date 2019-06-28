// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/globocom/huskyCI/client/types"
	"github.com/globocom/huskyCI/client/util"
)

var outputJSON types.JSONOutput
var goResults types.GoResults
var pythonResults types.PythonResults
var javaScriptResults types.JavaScriptResults
var rubyResults types.RubyResults

// CheckMongoDBContainerOutput will validate the output of a given container.
func CheckMongoDBContainerOutput(container types.Container) {

	switch container.SecurityTest.Name {
	case "enry":
	case "gosec":
		PrepareGosecOutput(container.COutput, container.CInfo)
		outputJSON.GoResults = goResults
	case "bandit":
		PrepareBanditOutput(container.COutput, container.CInfo)
		outputJSON.PythonResults.BanditOutput = pythonResults.BanditOutput
	case "retirejs":
		PrepareRetirejsOutput(container.COutput, container.CInfo)
		outputJSON.JavaScriptResults.RetirejsResult = javaScriptResults.RetirejsResult
	case "brakeman":
		PrepareBrakemanOutput(container.COutput, container.CInfo)
		outputJSON.RubyResults.BrakemanOutput = rubyResults.BrakemanOutput
	case "safety":
		PrepareSafetyOutput(container.COutput, container.CInfo)
		outputJSON.PythonResults.SafetyOutput = pythonResults.SafetyOutput
	default:
		fmt.Println("[HUSKYCI][ERROR] securityTest name not recognized:", container.SecurityTest.Name)
		os.Exit(1)
	}
}

// PrepareGosecOutput will prepare Gosec output.
func PrepareGosecOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	foundVuln := false
	gosecOutput := types.GosecOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &gosecOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal gosecOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, issue := range gosecOutput.GosecIssues {
		gosecVuln := types.HuskyCIVulnerability{}
		gosecVuln.SecurityTool = "gosec"
		gosecVuln.Severity = issue.Severity
		gosecVuln.Confidence = issue.Confidence
		gosecVuln.Details = issue.Details
		gosecVuln.File = issue.File
		gosecVuln.Line = issue.Line
		gosecVuln.Code = issue.Code

		goResults.GosecOutput = append(goResults.GosecOutput, gosecVuln)
	}

	if foundVuln {
		types.FoundVuln = true
	}
}

// PrepareBanditOutput will prepare Bandit output.
func PrepareBanditOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	foundVuln := false
	banditOutput := types.BanditOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &banditOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal banditOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, issue := range banditOutput.Results {
		banditVuln := types.HuskyCIVulnerability{}
		banditVuln.SecurityTool = "bandit"
		banditVuln.Severity = issue.IssueSeverity
		banditVuln.Confidence = issue.IssueConfidence
		banditVuln.Details = issue.IssueText
		banditVuln.File = issue.Filename
		banditVuln.Line = strconv.Itoa(issue.LineNumber)
		banditVuln.Code = issue.Code

		pythonResults.BanditOutput = append(pythonResults.BanditOutput, banditVuln)
	}

	if foundVuln {
		types.FoundVuln = true
	}
}

// PrepareRetirejsOutput will prepare Retirejs output.
func PrepareRetirejsOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	if strings.Contains(mongoDBcontainerInfo, "ERROR_RUNNING_RETIREJS") {
		retirejsVuln := types.HuskyCIVulnerability{}
		retirejsVuln.SecurityTool = "retirejs"
		retirejsVuln.Severity = "info"
		retirejsVuln.Confidence = "high"
		retirejsVuln.Details = "It looks like your project doesn't have package.json or yarn.lock. huskyCI was not able to run RetireJS properly."

		javaScriptResults.RetirejsResult = append(javaScriptResults.RetirejsResult, retirejsVuln)

		return
	}

	foundVuln := false
	retirejsOutput := []types.ResultsStruct{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &retirejsOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal retirejsOutput!", err)
		os.Exit(1)
	}

	for _, output := range retirejsOutput {
		for _, result := range output.Results {
			for _, vulnerability := range result.Vulnerabilities {
				retirejsVuln := types.HuskyCIVulnerability{}
				retirejsVuln.SecurityTool = "retirejs"
				retirejsVuln.Severity = vulnerability.Severity
				retirejsVuln.Code = result.Component
				retirejsVuln.Version = result.Version
				for _, info := range vulnerability.Info {
					retirejsVuln.Details = retirejsVuln.Details + info
				}
				retirejsVuln.Details = retirejsVuln.Details + vulnerability.Identifiers.Summary

				javaScriptResults.RetirejsResult = append(javaScriptResults.RetirejsResult, retirejsVuln)
			}
		}
	}

	if foundVuln {
		types.FoundVuln = true
	}
}

// PrepareBrakemanOutput will prepare Brakeman output.
func PrepareBrakemanOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	foundVuln := false
	brakemanOutput := types.BrakemanOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &brakemanOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal brakemanOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, warning := range brakemanOutput.Warnings {
		brakemanVuln := types.HuskyCIVulnerability{}
		brakemanVuln.SecurityTool = "brakeman"
		brakemanVuln.Confidence = warning.Confidence
		brakemanVuln.Details = warning.Details + warning.Message
		brakemanVuln.File = warning.File
		brakemanVuln.Line = strconv.Itoa(warning.Line)
		brakemanVuln.Code = warning.Code
		brakemanVuln.Type = warning.Type

		rubyResults.BrakemanOutput = append(rubyResults.BrakemanOutput, brakemanVuln)
	}

	if foundVuln {
		types.FoundVuln = true
	}
}

// PrepareSafetyOutput will prepare Safety output.
func PrepareSafetyOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	if mongoDBcontainerInfo == "Requirements not found or this project uses latest dependencies." {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.SecurityTool = "safety"
		safetyVuln.Severity = "info"
		safetyVuln.Confidence = "high"
		safetyVuln.Details = "requirements.txt not found or this project uses latest dependencies"

		pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)

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
			safetyVuln := types.HuskyCIVulnerability{}
			safetyVuln.SecurityTool = "safety"
			safetyVuln.Severity = "warning"
			safetyVuln.Details = warning

			pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)
		}
		if onlyWarning {
			return
		}
	}

	for _, issue := range safetyOutput.SafetyIssues {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.SecurityTool = "safety"
		safetyVuln.Severity = "high"
		safetyVuln.Details = issue.Comment
		safetyVuln.Code = issue.Version
		safetyVuln.VunerableBelow = issue.Below

		pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)
	}

	types.FoundVuln = true
}

// PrintJSONOutput prints the analysis output in a JSON format
func PrintJSONOutput() error {
	jsonReady, err := json.Marshal(outputJSON)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonReady))
	return nil
}

// PrinthuskyCIOutput prints the analysis output in huskyCI's format
func PrinthuskyCIOutput() {
	for _, issue := range outputJSON.GoResults.GosecOutput {
		if (issue.Severity == "HIGH") && (issue.Confidence == "HIGH") {
			color.Red("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Red("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Red("[HUSKYCI][!] Details: %s", issue.Details)
			color.Red("[HUSKYCI][!] File: %s", issue.File)
			color.Red("[HUSKYCI][!] Line: %s", issue.Line)
			color.Red("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		} else if (issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			color.Yellow("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Yellow("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.Details)
			color.Yellow("[HUSKYCI][!] File: %s", issue.File)
			color.Yellow("[HUSKYCI][!] Line: %s", issue.Line)
			color.Yellow("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		} else if issue.Severity == "LOW" {
			color.Blue("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Blue("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Blue("[HUSKYCI][!] Details: %s", issue.Details)
			color.Blue("[HUSKYCI][!] File: %s", issue.File)
			color.Blue("[HUSKYCI][!] Line: %s", issue.Line)
			color.Blue("[HUSKYCI][!] Code: %s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range outputJSON.PythonResults.BanditOutput {
		if (issue.Severity == "HIGH") && (issue.Confidence == "HIGH") {
			color.Red("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Red("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Red("[HUSKYCI][!] Details: %s", issue.Details)
			color.Red("[HUSKYCI][!] File: %s", issue.File)
			color.Red("[HUSKYCI][!] Line: %s", issue.Line)
			color.Red("[HUSKYCI][!] Code:\n%s", issue.Code)
			fmt.Println()
		} else if (issue.Severity == "MEDIUM") && (issue.Confidence == "HIGH") {
			color.Yellow("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Yellow("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.Details)
			color.Yellow("[HUSKYCI][!] File: %s", issue.File)
			color.Yellow("[HUSKYCI][!] Line: %s", issue.Line)
			color.Yellow("[HUSKYCI][!] Code:\n%s", issue.Code)
			fmt.Println()
		} else if issue.Severity == "LOW" {
			color.Blue("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Blue("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Blue("[HUSKYCI][!] Details: %s", issue.Details)
			color.Blue("[HUSKYCI][!] File: %s", issue.File)
			color.Blue("[HUSKYCI][!] Line: %s", issue.Line)
			color.Blue("[HUSKYCI][!] Code:\n%s", issue.Code)
			fmt.Println()
		}
	}

	for _, issue := range outputJSON.PythonResults.SafetyOutput {
		color.Red("[HUSKYCI][!] Severity: %s", issue.Severity)
		color.Red("[HUSKYCI][!] Details: %s", issue.Details)
		color.Red("[HUSKYCI][!] Code:\n%s", issue.Code)
		color.Red("[HUSKYCI][!] Vulnerable Below:\n%s", issue.VunerableBelow)
		fmt.Println()
	}

	for _, issue := range outputJSON.RubyResults.BrakemanOutput {
		if issue.Confidence == "High" {
			color.Red("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Red("[HUSKYCI][!] Details: %s", issue.Details)
			color.Red("[HUSKYCI][!] File: %s", issue.File)
			color.Red("[HUSKYCI][!] Line: %s", issue.Line)
			color.Red("[HUSKYCI][!] Code: %s", issue.Code)
			color.Red("[HUSKYCI][!] Type: %s", issue.Type)
			fmt.Println()
		} else if issue.Confidence == "MEDIUM" {
			color.Yellow("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.Details)
			color.Yellow("[HUSKYCI][!] File: %s", issue.File)
			color.Yellow("[HUSKYCI][!] Line: %s", issue.Line)
			color.Yellow("[HUSKYCI][!] Code: %s", issue.Code)
			color.Yellow("[HUSKYCI][!] Type: %s", issue.Type)
			fmt.Println()
		} else if issue.Confidence == "LOW" {
			color.Blue("[HUSKYCI][!] Confidence: %s", issue.Confidence)
			color.Blue("[HUSKYCI][!] Details: %s", issue.Details)
			color.Blue("[HUSKYCI][!] File: %s", issue.File)
			color.Blue("[HUSKYCI][!] Line: %s", issue.Line)
			color.Blue("[HUSKYCI][!] Code: %s", issue.Code)
			color.Blue("[HUSKYCI][!] Type: %s", issue.Type)
			fmt.Println()
		}
	}

	for _, issue := range outputJSON.JavaScriptResults.RetirejsResult {
		if issue.Severity == "high" {
			color.Red("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Red("[HUSKYCI][!] Code: %s", issue.Code)
			color.Red("[HUSKYCI][!] Version: %s", issue.Version)
			color.Red("[HUSKYCI][!] Details: %s", issue.Details)
			fmt.Println()
		} else if issue.Severity == "medium" {
			color.Yellow("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Yellow("[HUSKYCI][!] Code: %s", issue.Code)
			color.Yellow("[HUSKYCI][!] Version: %s", issue.Version)
			color.Yellow("[HUSKYCI][!] Details: %s", issue.Details)
			fmt.Println()
		} else if issue.Severity == "low" {
			color.Blue("[HUSKYCI][!] Severity: %s", issue.Severity)
			color.Blue("[HUSKYCI][!] Code: %s", issue.Code)
			color.Blue("[HUSKYCI][!] Version: %s", issue.Version)
			color.Blue("[HUSKYCI][!] Details: %s", issue.Details)
			fmt.Println()
		}

	}
}

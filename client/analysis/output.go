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

	"github.com/globocom/huskyCI/client/types"
	"github.com/globocom/huskyCI/client/util"
)

var outputJSON types.JSONOutput
var goResults types.GoResults
var pythonResults types.PythonResults
var javaScriptResults types.JavaScriptResults
var rubyResults types.RubyResults

// prepareSecurityTestResult preares the output of a given securityTest.
func prepareSecurityTestResult(container types.Container) {

	switch container.SecurityTest.Name {
	case "enry":
	case "gosec":
		prepareGosecOutput(container.COutput, container.CInfo)
		outputJSON.GoResults = goResults
	case "bandit":
		prepareBanditOutput(container.COutput, container.CInfo)
		outputJSON.PythonResults.BanditOutput = pythonResults.BanditOutput
	case "retirejs":
		prepareRetirejsOutput(container.COutput, container.CInfo)
		outputJSON.JavaScriptResults.RetirejsResult = javaScriptResults.RetirejsResult
	case "brakeman":
		prepareBrakemanOutput(container.COutput, container.CInfo)
		outputJSON.RubyResults.BrakemanOutput = rubyResults.BrakemanOutput
	case "safety":
		prepareSafetyOutput(container.COutput, container.CInfo)
		outputJSON.PythonResults.SafetyOutput = pythonResults.SafetyOutput
	default:
		fmt.Println("[HUSKYCI][ERROR] securityTest name not recognized:", container.SecurityTest.Name)
		os.Exit(1)
	}
}

// prepareGosecOutput will prepare Gosec output.
func prepareGosecOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	gosecOutput := types.GosecOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &gosecOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal gosecOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, issue := range gosecOutput.GosecIssues {
		gosecVuln := types.HuskyCIVulnerability{}
		gosecVuln.Language = "Go"
		gosecVuln.SecurityTool = "GoSec"
		gosecVuln.Severity = issue.Severity
		gosecVuln.Confidence = issue.Confidence
		gosecVuln.Details = issue.Details
		gosecVuln.File = issue.File
		gosecVuln.Line = issue.Line
		gosecVuln.Code = issue.Code

		goResults.GosecOutput = append(goResults.GosecOutput, gosecVuln)

		if ((issue.Severity == "MEDIUM") || (issue.Severity == "HIGH")) && (issue.Confidence == "HIGH") {
			types.FoundVuln = true
		} else if issue.Severity == "LOW" {
			types.FoundInfo = true
		}
	}
}

// prepareBanditOutput will prepare Bandit output.
func prepareBanditOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	banditOutput := types.BanditOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &banditOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal banditOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, issue := range banditOutput.Results {
		banditVuln := types.HuskyCIVulnerability{}
		banditVuln.Language = "Python"
		banditVuln.SecurityTool = "Bandit"
		banditVuln.Severity = issue.IssueSeverity
		banditVuln.Confidence = issue.IssueConfidence
		banditVuln.Details = issue.IssueText
		banditVuln.File = issue.Filename
		banditVuln.Line = strconv.Itoa(issue.LineNumber)
		banditVuln.Code = issue.Code

		pythonResults.BanditOutput = append(pythonResults.BanditOutput, banditVuln)

		if ((issue.IssueSeverity == "MEDIUM") || (issue.IssueSeverity == "HIGH")) && (issue.IssueConfidence == "HIGH") {
			types.FoundVuln = true
		} else if issue.IssueSeverity == "LOW" {
			types.FoundInfo = true
		}
	}
}

// prepareRetirejsOutput will prepare Retirejs output.
func prepareRetirejsOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	if strings.Contains(mongoDBcontainerInfo, "ERROR_RUNNING_RETIREJS") {
		retirejsVuln := types.HuskyCIVulnerability{}
		retirejsVuln.Language = "JavaScript"
		retirejsVuln.SecurityTool = "RetireJS"
		retirejsVuln.Severity = "info"
		retirejsVuln.Confidence = "high"
		retirejsVuln.Details = "It looks like your project doesn't have package.json or yarn.lock. huskyCI was not able to run RetireJS properly."

		javaScriptResults.RetirejsResult = append(javaScriptResults.RetirejsResult, retirejsVuln)
		types.FoundInfo = true

		return
	}

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
				retirejsVuln.Language = "JavaScript"
				retirejsVuln.SecurityTool = "RetireJS"
				retirejsVuln.Severity = vulnerability.Severity
				retirejsVuln.Code = result.Component
				retirejsVuln.Version = result.Version
				for _, info := range vulnerability.Info {
					retirejsVuln.Details = retirejsVuln.Details + info + "\n"
				}
				retirejsVuln.Details = retirejsVuln.Details + vulnerability.Identifiers.Summary

				javaScriptResults.RetirejsResult = append(javaScriptResults.RetirejsResult, retirejsVuln)

				if retirejsVuln.Severity == "high" || retirejsVuln.Severity == "medium" {
					types.FoundVuln = true
				} else if retirejsVuln.Severity == "low" {
					types.FoundInfo = true
				}
			}
		}
	}
}

// prepareBrakemanOutput will prepare Brakeman output.
func prepareBrakemanOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	brakemanOutput := types.BrakemanOutput{}
	err := json.Unmarshal([]byte(mongoDBcontainerOutput), &brakemanOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not Unmarshal brakemanOutput!", mongoDBcontainerOutput)
		os.Exit(1)
	}

	for _, warning := range brakemanOutput.Warnings {
		brakemanVuln := types.HuskyCIVulnerability{}
		brakemanVuln.Language = "Ruby"
		brakemanVuln.SecurityTool = "Brakeman"
		brakemanVuln.Confidence = warning.Confidence
		brakemanVuln.Details = warning.Details + warning.Message
		brakemanVuln.File = warning.File
		brakemanVuln.Line = strconv.Itoa(warning.Line)
		brakemanVuln.Code = warning.Code
		brakemanVuln.Type = warning.Type

		rubyResults.BrakemanOutput = append(rubyResults.BrakemanOutput, brakemanVuln)

		if brakemanVuln.Confidence == "High" || brakemanVuln.Confidence == "Medium" {
			brakemanVuln.Severity = "High"
			types.FoundVuln = true
		} else {
			brakemanVuln.Severity = "Low"
			types.FoundInfo = true
		}
	}
}

// prepareSafetyOutput will prepare Safety output.
func prepareSafetyOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	if mongoDBcontainerInfo == "Requirements not found." {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
		safetyVuln.Details = "requirements.txt not found"
		types.FoundInfo = true

		pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)

		return
	}

	if mongoDBcontainerInfo == "Internal error running Safety." {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
		safetyVuln.Details = "Internal error running Safety."
		types.FoundInfo = true

		pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)

		return
	}

	// Safety returns warnings and the json output in the same string, which need to be split
	var cOutputSanitized string
	safetyOutput := types.SafetyOutput{}
	warningFound := strings.Contains(mongoDBcontainerOutput, "Warning: unpinned requirement ")
	if !warningFound {
		// only issues found
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
			safetyVuln.Language = "Python"
			safetyVuln.SecurityTool = "Safety"
			safetyVuln.Severity = "warning"
			safetyVuln.Details = util.AdjustWarningMessage(warning)

			pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)
			types.FoundInfo = true
		}
		if onlyWarning {
			return
		}
	}

	for _, issue := range safetyOutput.SafetyIssues {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "high"
		safetyVuln.Details = issue.Comment
		safetyVuln.Code = issue.Dependency + " " + issue.Version
		safetyVuln.VunerableBelow = issue.Below

		pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)
		types.FoundVuln = true
	}
}

// printJSONOutput prints the analysis output in a JSON format
func printJSONOutput() error {
	jsonReady := []byte{}
	var err error
	if jsonReady, err = json.Marshal(outputJSON); err != nil {
		return err
	}
	fmt.Println(string(jsonReady))
	return nil
}

// printSTDOUTOutput prints the analysis output in STDOUT using printfs
func printSTDOUTOutput() {

	for _, issue := range outputJSON.GoResults.GosecOutput {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
	}

	for _, issue := range outputJSON.PythonResults.BanditOutput {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code:\n%s\n", issue.Code)
	}

	for _, issue := range outputJSON.PythonResults.SafetyOutput {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		if issue.Details != "requirements.txt not found" && !strings.Contains(issue.Details, "Unpinned requirement ") {
			fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
			fmt.Printf("[HUSKYCI][!] Vulnerable Below: %s\n", issue.VunerableBelow)
		}
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
	}

	for _, issue := range outputJSON.RubyResults.BrakemanOutput {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		fmt.Printf("[HUSKYCI][!] Type: %s\n", issue.Type)
	}

	for _, issue := range outputJSON.JavaScriptResults.RetirejsResult {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		if !strings.Contains(issue.Details, "doesn't have package.json") {
			fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
			fmt.Printf("[HUSKYCI][!] Version: %s\n", issue.Version)
		}
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
	}

	printAllSummary()
}

// prepareAllSummary prepares how many low, medium and high vulnerabilites were found.
func prepareAllSummary() {
	var totalLow, totalMedium, totalHigh int

	// GoSec summary
	for _, issue := range outputJSON.GoResults.GosecOutput {
		switch issue.Severity {
		case "LOW":
			outputJSON.Summary.GosecSummary.FoundInfo = true
			outputJSON.Summary.GosecSummary.LowVuln++
		case "MEDIUM":
			outputJSON.Summary.GosecSummary.FoundVuln = true
			outputJSON.Summary.GosecSummary.MediumVuln++
		case "HIGH":
			outputJSON.Summary.GosecSummary.FoundVuln = true
			outputJSON.Summary.GosecSummary.HighVuln++
		}
	}

	// Bandit summary
	for _, issue := range outputJSON.PythonResults.BanditOutput {
		switch issue.Severity {
		case "LOW":
			outputJSON.Summary.BanditSummary.FoundInfo = true
			outputJSON.Summary.BanditSummary.LowVuln++
		case "MEDIUM":
			outputJSON.Summary.BanditSummary.FoundVuln = true
			outputJSON.Summary.BanditSummary.MediumVuln++
		case "HIGH":
			outputJSON.Summary.BanditSummary.FoundVuln = true
			outputJSON.Summary.BanditSummary.HighVuln++
		}
	}

	// Safety summary
	for _, issue := range outputJSON.PythonResults.SafetyOutput {
		switch issue.Severity {
		case "info":
			outputJSON.Summary.SafetySummary.FoundInfo = true
			outputJSON.Summary.SafetySummary.LowVuln++
		case "warning":
			outputJSON.Summary.SafetySummary.FoundInfo = true
			outputJSON.Summary.SafetySummary.LowVuln++
		case "high":
			outputJSON.Summary.SafetySummary.FoundVuln = true
			outputJSON.Summary.SafetySummary.HighVuln++
		}
	}

	// Brakeman summary
	for _, issue := range outputJSON.RubyResults.BrakemanOutput {
		switch issue.Severity {
		case "Low":
			outputJSON.Summary.BrakemanSummary.FoundVuln = true
			outputJSON.Summary.BrakemanSummary.LowVuln++
		case "Medium":
			outputJSON.Summary.BrakemanSummary.FoundVuln = true
			outputJSON.Summary.BrakemanSummary.MediumVuln++
		case "High":
			outputJSON.Summary.BrakemanSummary.FoundVuln = true
			outputJSON.Summary.BrakemanSummary.HighVuln++
		}
	}

	// RetireJS summary
	for _, issue := range outputJSON.JavaScriptResults.RetirejsResult {
		switch issue.Severity {
		case "info":
			outputJSON.Summary.RetirejsSummary.FoundInfo = true
			outputJSON.Summary.RetirejsSummary.LowVuln++
		case "low":
			outputJSON.Summary.RetirejsSummary.FoundInfo = true
			outputJSON.Summary.RetirejsSummary.LowVuln++
		case "medium":
			outputJSON.Summary.RetirejsSummary.FoundVuln = true
			outputJSON.Summary.RetirejsSummary.MediumVuln++
		case "high":
			outputJSON.Summary.RetirejsSummary.FoundVuln = true
			outputJSON.Summary.RetirejsSummary.HighVuln++
		}
	}

	// Total summary
	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.BanditSummary.FoundVuln || outputJSON.Summary.SafetySummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.RetirejsSummary.FoundVuln {
		outputJSON.Summary.TotalSummary.FoundVuln = true
	} else if outputJSON.Summary.GosecSummary.FoundInfo || outputJSON.Summary.BanditSummary.FoundInfo || outputJSON.Summary.SafetySummary.FoundInfo || outputJSON.Summary.BrakemanSummary.FoundInfo || outputJSON.Summary.RetirejsSummary.FoundInfo {
		outputJSON.Summary.TotalSummary.FoundInfo = true
	}

	totalLow = outputJSON.Summary.RetirejsSummary.LowVuln + outputJSON.Summary.BrakemanSummary.LowVuln + outputJSON.Summary.SafetySummary.LowVuln + outputJSON.Summary.BanditSummary.LowVuln + outputJSON.Summary.GosecSummary.LowVuln
	totalMedium = outputJSON.Summary.RetirejsSummary.MediumVuln + outputJSON.Summary.BrakemanSummary.MediumVuln + outputJSON.Summary.SafetySummary.MediumVuln + outputJSON.Summary.BanditSummary.MediumVuln + outputJSON.Summary.GosecSummary.MediumVuln
	totalHigh = outputJSON.Summary.RetirejsSummary.HighVuln + outputJSON.Summary.BrakemanSummary.HighVuln + outputJSON.Summary.SafetySummary.HighVuln + outputJSON.Summary.BanditSummary.HighVuln + outputJSON.Summary.GosecSummary.HighVuln

	outputJSON.Summary.TotalSummary.HighVuln = totalHigh
	outputJSON.Summary.TotalSummary.MediumVuln = totalMedium
	outputJSON.Summary.TotalSummary.LowVuln = totalLow

}

func printAllSummary() {

	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.GosecSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Go -> GoSec\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.GosecSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.GosecSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.GosecSummary.LowVuln)
	}

	if outputJSON.Summary.BanditSummary.FoundVuln || outputJSON.Summary.BanditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Python -> Bandit\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.BanditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.BanditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.BanditSummary.LowVuln)
	}

	if outputJSON.Summary.SafetySummary.FoundVuln || outputJSON.Summary.SafetySummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Python -> Safety\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.SafetySummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.SafetySummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.SafetySummary.LowVuln)
	}

	if outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Ruby -> Brakeman\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.BrakemanSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.BrakemanSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.BrakemanSummary.LowVuln)
	}

	if outputJSON.Summary.RetirejsSummary.FoundVuln || outputJSON.Summary.RetirejsSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] JavaScript -> RetireJS\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.RetirejsSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.RetirejsSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.RetirejsSummary.LowVuln)
	}

	if outputJSON.Summary.TotalSummary.FoundVuln || outputJSON.Summary.TotalSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Total\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.TotalSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.TotalSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.TotalSummary.LowVuln)
		fmt.Println()
	}
}

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

// PrepareBanditOutput will prepare Bandit output.
func PrepareBanditOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

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

// PrepareRetirejsOutput will prepare Retirejs output.
func PrepareRetirejsOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

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
				retirejsVuln.Language = "JavaScript"
				retirejsVuln.SecurityTool = "RetireJS"
				retirejsVuln.Severity = vulnerability.Severity
				retirejsVuln.Code = result.Component
				retirejsVuln.Version = result.Version
				for _, info := range vulnerability.Info {
					retirejsVuln.Details = retirejsVuln.Details + info
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

	if foundVuln {
		types.FoundVuln = true
	}
}

// PrepareBrakemanOutput will prepare Brakeman output.
func PrepareBrakemanOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

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
			types.FoundVuln = true
		} else {
			types.FoundInfo = true
		}
	}
}

// PrepareSafetyOutput will prepare Safety output.
func PrepareSafetyOutput(mongoDBcontainerOutput string, mongoDBcontainerInfo string) {

	if mongoDBcontainerInfo == "No issues found." {
		return
	}

	if mongoDBcontainerInfo == "Requirements not found or this project uses latest dependencies." {
		safetyVuln := types.HuskyCIVulnerability{}
		safetyVuln.Language = "Python"
		safetyVuln.SecurityTool = "Safety"
		safetyVuln.Severity = "info"
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
			safetyVuln.Language = "Python"
			safetyVuln.SecurityTool = "Safety"
			safetyVuln.Severity = "warning"
			safetyVuln.Details = warning

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
		safetyVuln.Code = issue.Version
		safetyVuln.VunerableBelow = issue.Below

		pythonResults.SafetyOutput = append(pythonResults.SafetyOutput, safetyVuln)
		types.FoundVuln = true
	}
}

// PrintJSONOutput prints the analysis output in a JSON format
func printJSONOutput() error {
	jsonReady := []byte{}
	var err error
	if jsonReady, err = json.Marshal(outputJSON); err != nil {
		return err
	}
	fmt.Println(string(jsonReady))
	return nil
}

// PrinthuskyCIOutput prints the analysis output in huskyCI's format
func printhuskyCIOutput() {
	for _, issue := range outputJSON.GoResults.GosecOutput {
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		fmt.Println()
	}

	for _, issue := range outputJSON.PythonResults.BanditOutput {
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code:\n%s\n", issue.Code)
		fmt.Println()
	}

	for _, issue := range outputJSON.PythonResults.SafetyOutput {
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		fmt.Printf("[HUSKYCI][!] Vulnerable Below: %s\n", issue.VunerableBelow)
		fmt.Println()
	}

	for _, issue := range outputJSON.RubyResults.BrakemanOutput {
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		fmt.Printf("[HUSKYCI][!] Type: %s\n", issue.Type)
		fmt.Println()
	}

	for _, issue := range outputJSON.JavaScriptResults.RetirejsResult {
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		fmt.Printf("[HUSKYCI][!] Version: %s\n", issue.Version)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Println()
	}
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/client/util"

	"github.com/globocom/huskyCI/client/types"
)

var outputJSON types.JSONOutput

// generateSonarOutput prints the analysis output in a JSON format
func generateSonarOutput() error {

	var allVulns []types.HuskyCIVulnerability

	// gosec
	allVulns = append(allVulns, outputJSON.GoResults.HuskyCIGosecOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.GoResults.HuskyCIGosecOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.GoResults.HuskyCIGosecOutput.HighVulns...)

	// bandit
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCIBanditOutput.NoSecVulns...)
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCIBanditOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCIBanditOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCIBanditOutput.HighVulns...)

	// safety
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCISafetyOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCISafetyOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.PythonResults.HuskyCISafetyOutput.HighVulns...)

	// brakeman
	allVulns = append(allVulns, outputJSON.RubyResults.HuskyCIBrakemanOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.RubyResults.HuskyCIBrakemanOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.RubyResults.HuskyCIBrakemanOutput.HighVulns...)

	// npmaudit
	allVulns = append(allVulns, outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns...)

	// yarnaudit
	allVulns = append(allVulns, outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns...)

	// gitleaks
	allVulns = append(allVulns, outputJSON.GenericResults.HuskyCIGitleaksOutput.LowVulns...)
	allVulns = append(allVulns, outputJSON.GenericResults.HuskyCIGitleaksOutput.MediumVulns...)
	allVulns = append(allVulns, outputJSON.GenericResults.HuskyCIGitleaksOutput.HighVulns...)

	var sonarOutput types.HuskyCISonarOutput

	for _, vuln := range allVulns {
		var issue types.SonarIssue
		issue.EngineID = "huskyCI"
		issue.Type = "VULNERABILITY"
		issue.RuleID = vuln.Language
		switch strings.ToLower(vuln.Severity) {
		case `low`:
			issue.Severity = "MINOR"
		case `medium`:
			issue.Severity = "MAJOR"
		case `high`:
			issue.Severity = "BLOCKER"
		default:
			issue.Severity = "INFO"
		}
		issue.PrimaryLocation.FilePath = vuln.File
		issue.PrimaryLocation.Message = vuln.Details
		issue.PrimaryLocation.TextRange.StartLine = 0
		lineNum, err := strconv.Atoi(vuln.Line)
		if err != nil {
			lineNum = 0
		}
		if lineNum != 0 && lineNum > 0 {
			issue.PrimaryLocation.TextRange.StartLine = lineNum
		}
		sonarOutput.Issues = append(sonarOutput.Issues, issue)
	}

	sonarOutputString, err := json.Marshal(sonarOutput)
	if err != nil {
		return err
	}

	err = util.CreateSonarJSONFile(sonarOutputString)
	if err != nil {
		return err
	}

	return nil
}

// printSTDOUTOutput prints the analysis output in STDOUT using printfs
func printSTDOUTOutput(analysis types.Analysis) {

	// gosec
	printSTDOUTOutputGosec(outputJSON.GoResults.HuskyCIGosecOutput.LowVulns)
	printSTDOUTOutputGosec(outputJSON.GoResults.HuskyCIGosecOutput.MediumVulns)
	printSTDOUTOutputGosec(outputJSON.GoResults.HuskyCIGosecOutput.HighVulns)

	// bandit
	printSTDOUTOutputBandit(outputJSON.PythonResults.HuskyCIBanditOutput.LowVulns)
	printSTDOUTOutputBandit(outputJSON.PythonResults.HuskyCIBanditOutput.MediumVulns)
	printSTDOUTOutputBandit(outputJSON.PythonResults.HuskyCIBanditOutput.HighVulns)

	// safety
	printSTDOUTOutputSafety(outputJSON.PythonResults.HuskyCISafetyOutput.LowVulns)
	printSTDOUTOutputSafety(outputJSON.PythonResults.HuskyCISafetyOutput.MediumVulns)
	printSTDOUTOutputSafety(outputJSON.PythonResults.HuskyCISafetyOutput.HighVulns)

	// brakeman
	printSTDOUTOutputBrakeman(outputJSON.RubyResults.HuskyCIBrakemanOutput.LowVulns)
	printSTDOUTOutputBrakeman(outputJSON.RubyResults.HuskyCIBrakemanOutput.MediumVulns)
	printSTDOUTOutputBrakeman(outputJSON.RubyResults.HuskyCIBrakemanOutput.HighVulns)

	// npmaudit
	printSTDOUTOutputNpmAudit(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns)
	printSTDOUTOutputNpmAudit(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns)
	printSTDOUTOutputNpmAudit(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns)

	// yarnaudit
	printSTDOUTOutputYarnAudit(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns)
	printSTDOUTOutputYarnAudit(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns)
	printSTDOUTOutputYarnAudit(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns)

	// gitleaks
	printSTDOUTOutputGitleaks(outputJSON.GenericResults.HuskyCIGitleaksOutput.LowVulns)
	printSTDOUTOutputGitleaks(outputJSON.GenericResults.HuskyCIGitleaksOutput.MediumVulns)
	printSTDOUTOutputGitleaks(outputJSON.GenericResults.HuskyCIGitleaksOutput.HighVulns)

	printAllSummary(analysis)
}

// prepareAllSummary prepares how many low, medium and high vulnerabilites were found.
func prepareAllSummary(analysis types.Analysis) {
	var totalNoSec, totalLow, totalMedium, totalHigh int

	outputJSON.GoResults = analysis.HuskyCIResults.GoResults
	outputJSON.JavaScriptResults = analysis.HuskyCIResults.JavaScriptResults
	outputJSON.PythonResults = analysis.HuskyCIResults.PythonResults
	outputJSON.RubyResults = analysis.HuskyCIResults.RubyResults
	outputJSON.GenericResults = analysis.HuskyCIResults.GenericResults

	// GoSec summary
	outputJSON.Summary.GosecSummary.LowVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.LowVulns)
	outputJSON.Summary.GosecSummary.MediumVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.MediumVulns)
	outputJSON.Summary.GosecSummary.HighVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.HighVulns)
	if len(outputJSON.GoResults.HuskyCIGosecOutput.LowVulns) > 0 {
		outputJSON.Summary.GosecSummary.FoundInfo = true
	}
	if len(outputJSON.GoResults.HuskyCIGosecOutput.MediumVulns) > 0 || len(outputJSON.GoResults.HuskyCIGosecOutput.HighVulns) > 0 {
		outputJSON.Summary.GosecSummary.FoundVuln = true
	}

	// Bandit summary
	outputJSON.Summary.BanditSummary.NoSecVuln = len(outputJSON.PythonResults.HuskyCIBanditOutput.NoSecVulns)
	outputJSON.Summary.BanditSummary.LowVuln = len(outputJSON.PythonResults.HuskyCIBanditOutput.LowVulns)
	outputJSON.Summary.BanditSummary.MediumVuln = len(outputJSON.PythonResults.HuskyCIBanditOutput.MediumVulns)
	outputJSON.Summary.BanditSummary.HighVuln = len(outputJSON.PythonResults.HuskyCIBanditOutput.HighVulns)
	if len(outputJSON.PythonResults.HuskyCIBanditOutput.LowVulns) > 0 {
		outputJSON.Summary.BanditSummary.FoundInfo = true
	}
	if len(outputJSON.PythonResults.HuskyCIBanditOutput.MediumVulns) > 0 || len(outputJSON.PythonResults.HuskyCIBanditOutput.HighVulns) > 0 {
		outputJSON.Summary.BanditSummary.FoundVuln = true
	}

	// Safety summary
	outputJSON.Summary.SafetySummary.LowVuln = len(outputJSON.PythonResults.HuskyCISafetyOutput.LowVulns)
	outputJSON.Summary.SafetySummary.MediumVuln = len(outputJSON.PythonResults.HuskyCISafetyOutput.MediumVulns)
	outputJSON.Summary.SafetySummary.HighVuln = len(outputJSON.PythonResults.HuskyCISafetyOutput.HighVulns)
	if len(outputJSON.PythonResults.HuskyCISafetyOutput.LowVulns) > 0 {
		outputJSON.Summary.SafetySummary.FoundInfo = true
	}
	if len(outputJSON.PythonResults.HuskyCISafetyOutput.MediumVulns) > 0 || len(outputJSON.PythonResults.HuskyCISafetyOutput.HighVulns) > 0 {
		outputJSON.Summary.SafetySummary.FoundVuln = true
	}

	// Brakeman summary
	outputJSON.Summary.BrakemanSummary.LowVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.LowVulns)
	outputJSON.Summary.BrakemanSummary.MediumVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.MediumVulns)
	outputJSON.Summary.BrakemanSummary.HighVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.HighVulns)
	if len(outputJSON.RubyResults.HuskyCIBrakemanOutput.LowVulns) > 0 {
		outputJSON.Summary.BrakemanSummary.FoundInfo = true
	}
	if len(outputJSON.RubyResults.HuskyCIBrakemanOutput.MediumVulns) > 0 || len(outputJSON.RubyResults.HuskyCIBrakemanOutput.HighVulns) > 0 {
		outputJSON.Summary.BrakemanSummary.FoundVuln = true
	}

	// NpmAudit summary
	outputJSON.Summary.NpmAuditSummary.LowVuln = len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns)
	outputJSON.Summary.NpmAuditSummary.MediumVuln = len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns)
	outputJSON.Summary.NpmAuditSummary.HighVuln = len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns)
	if len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns) > 0 {
		outputJSON.Summary.NpmAuditSummary.FoundInfo = true
	}
	if len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns) > 0 {
		outputJSON.Summary.NpmAuditSummary.FoundVuln = true
	}

	// YarnAudit summary
	outputJSON.Summary.YarnAuditSummary.LowVuln = len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns)
	outputJSON.Summary.YarnAuditSummary.MediumVuln = len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns)
	outputJSON.Summary.YarnAuditSummary.HighVuln = len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns)
	if len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns) > 0 {
		outputJSON.Summary.YarnAuditSummary.FoundInfo = true
	}
	if len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns) > 0 {
		outputJSON.Summary.YarnAuditSummary.FoundVuln = true
	}

	// GitLeaks summary
	outputJSON.Summary.GitleaksSummary.LowVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.LowVulns)
	outputJSON.Summary.GitleaksSummary.MediumVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.MediumVulns)
	outputJSON.Summary.GitleaksSummary.HighVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.HighVulns)
	if len(outputJSON.GenericResults.HuskyCIGitleaksOutput.LowVulns) > 0 {
		outputJSON.Summary.GitleaksSummary.FoundInfo = true
	}
	if len(outputJSON.GenericResults.HuskyCIGitleaksOutput.MediumVulns) > 0 || len(outputJSON.GenericResults.HuskyCIGitleaksOutput.HighVulns) > 0 {
		outputJSON.Summary.GitleaksSummary.FoundVuln = true
	}

	// Total summary
	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.BanditSummary.FoundVuln || outputJSON.Summary.SafetySummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.NpmAuditSummary.FoundVuln || outputJSON.Summary.YarnAuditSummary.FoundVuln || outputJSON.Summary.GitleaksSummary.FoundVuln {
		outputJSON.Summary.TotalSummary.FoundVuln = true
		types.FoundVuln = true
	} else if outputJSON.Summary.GosecSummary.FoundInfo || outputJSON.Summary.BanditSummary.FoundInfo || outputJSON.Summary.SafetySummary.FoundInfo || outputJSON.Summary.BrakemanSummary.FoundInfo || outputJSON.Summary.NpmAuditSummary.FoundInfo || outputJSON.Summary.YarnAuditSummary.FoundInfo || outputJSON.Summary.GitleaksSummary.FoundInfo {
		outputJSON.Summary.TotalSummary.FoundInfo = true
		types.FoundInfo = true
	}

	totalNoSec = outputJSON.Summary.BanditSummary.NoSecVuln
	totalLow = outputJSON.Summary.BrakemanSummary.LowVuln + outputJSON.Summary.SafetySummary.LowVuln + outputJSON.Summary.BanditSummary.LowVuln + outputJSON.Summary.GosecSummary.LowVuln + outputJSON.Summary.NpmAuditSummary.LowVuln + outputJSON.Summary.YarnAuditSummary.LowVuln + outputJSON.Summary.GitleaksSummary.LowVuln
	totalMedium = outputJSON.Summary.BrakemanSummary.MediumVuln + outputJSON.Summary.SafetySummary.MediumVuln + outputJSON.Summary.BanditSummary.MediumVuln + outputJSON.Summary.GosecSummary.MediumVuln + outputJSON.Summary.NpmAuditSummary.MediumVuln + outputJSON.Summary.YarnAuditSummary.MediumVuln + outputJSON.Summary.GitleaksSummary.MediumVuln
	totalHigh = outputJSON.Summary.BrakemanSummary.HighVuln + outputJSON.Summary.SafetySummary.HighVuln + outputJSON.Summary.BanditSummary.HighVuln + outputJSON.Summary.GosecSummary.HighVuln + outputJSON.Summary.NpmAuditSummary.HighVuln + outputJSON.Summary.YarnAuditSummary.HighVuln + outputJSON.Summary.GitleaksSummary.HighVuln

	outputJSON.Summary.TotalSummary.HighVuln = totalHigh
	outputJSON.Summary.TotalSummary.MediumVuln = totalMedium
	outputJSON.Summary.TotalSummary.LowVuln = totalLow
	outputJSON.Summary.TotalSummary.NoSecVuln = totalNoSec

}

func printAllSummary(analysis types.Analysis) {

	var gosecVersion, banditVersion, safetyVersion, brakemanVersion, npmauditVersion, yarnauditVersion, gitleaksVersion string

	for _, container := range analysis.Containers {
		switch container.SecurityTest.Name {
		case "gosec":
			gosecVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "bandit":
			banditVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "safety":
			safetyVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "brakeman":
			brakemanVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "npmaudit":
			npmauditVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "yarnaudit":
			yarnauditVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "gitleaks":
			gitleaksVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		}
	}

	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.GosecSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Go -> %s\n", gosecVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.GosecSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.GosecSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.GosecSummary.LowVuln)
	}

	if outputJSON.Summary.BanditSummary.FoundVuln || outputJSON.Summary.BanditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Python -> %s\n", banditVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.BanditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.BanditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.BanditSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.BanditSummary.NoSecVuln)
	}

	if outputJSON.Summary.SafetySummary.FoundVuln || outputJSON.Summary.SafetySummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Python -> %s\n", safetyVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.SafetySummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.SafetySummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.SafetySummary.LowVuln)
	}

	if outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Ruby -> %s\n", brakemanVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.BrakemanSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.BrakemanSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.BrakemanSummary.LowVuln)
	}

	if outputJSON.Summary.NpmAuditSummary.FoundVuln || outputJSON.Summary.NpmAuditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] JavaScript -> %s\n", npmauditVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.NpmAuditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.NpmAuditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.NpmAuditSummary.LowVuln)
	}

	if outputJSON.Summary.YarnAuditSummary.FoundVuln || outputJSON.Summary.YarnAuditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] JavaScript -> %s\n", yarnauditVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.YarnAuditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.YarnAuditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.YarnAuditSummary.LowVuln)
	}

	if outputJSON.Summary.GitleaksSummary.FoundVuln || outputJSON.Summary.GitleaksSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Generic -> %s\n", gitleaksVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.GitleaksSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.GitleaksSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.GitleaksSummary.LowVuln)
	}

	if outputJSON.Summary.TotalSummary.FoundVuln || outputJSON.Summary.TotalSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Total\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.TotalSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.TotalSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.TotalSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.TotalSummary.NoSecVuln)
	}

	fmt.Println()
}

func printSTDOUTOutputGosec(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
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
}

func printSTDOUTOutputBandit(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
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
}

func printSTDOUTOutputSafety(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
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
}

func printSTDOUTOutputBrakeman(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
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
}

func printSTDOUTOutputNpmAudit(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		if !strings.Contains(issue.Details, "doesn't have package-lock.json.") {
			fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
			fmt.Printf("[HUSKYCI][!] Version: %s\n", issue.Version)
			fmt.Printf("[HUSKYCI][!] Vulnerable Below: %s\n", issue.VunerableBelow)
		}
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
	}
}

func printSTDOUTOutputYarnAudit(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		if !strings.Contains(issue.Details, "doesn't have yarn.lock.") {
			fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
			fmt.Printf("[HUSKYCI][!] Occurrences: %d\n", issue.Occurrences)
			fmt.Printf("[HUSKYCI][!] Version: %s\n", issue.Version)
			fmt.Printf("[HUSKYCI][!] Vulnerable Below: %s\n", issue.VunerableBelow)
		}
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
	}
}

func printSTDOUTOutputGitleaks(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
	}
}

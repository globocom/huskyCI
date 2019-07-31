// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/client/types"
)

var outputJSON types.JSONOutput
var goResults types.GoResults
var pythonResults types.PythonResults
var javaScriptResults types.JavaScriptResults
var rubyResults types.RubyResults

// prepareSecurityTestResult preares the output of a given securityTest.
func prepareSecurityTestResult(analysis types.Analysis) {
	outputJSON.GoResults = analysis.HuskyCIResults.GoResults
	outputJSON.JavaScriptResults = analysis.HuskyCIResults.JavaScriptResults
	outputJSON.JavaScriptResults = analysis.HuskyCIResults.JavaScriptResults
	outputJSON.PythonResults = analysis.HuskyCIResults.PythonResults
	outputJSON.PythonResults = analysis.HuskyCIResults.PythonResults
	outputJSON.RubyResults = analysis.HuskyCIResults.RubyResults
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

	// retirejs
	printSTDOUTOutputRetireJS(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.LowVulns)
	printSTDOUTOutputRetireJS(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.MediumVulns)
	printSTDOUTOutputRetireJS(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.HighVulns)

	// npmaudit
	printSTDOUTOutputNpmAudit(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns)
	printSTDOUTOutputNpmAudit(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns)
	printSTDOUTOutputNpmAudit(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns)

	printAllSummary()
}

// prepareAllSummary prepares how many low, medium and high vulnerabilites were found.
func prepareAllSummary() {
	var totalLow, totalMedium, totalHigh int

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

	// RetireJS summary
	outputJSON.Summary.RetirejsSummary.LowVuln = len(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.LowVulns)
	outputJSON.Summary.RetirejsSummary.MediumVuln = len(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.MediumVulns)
	outputJSON.Summary.RetirejsSummary.HighVuln = len(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.HighVulns)
	if len(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.LowVulns) > 0 {
		outputJSON.Summary.RetirejsSummary.FoundInfo = true
	}
	if len(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.MediumVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCIRetireJSOutput.HighVulns) > 0 {
		outputJSON.Summary.RetirejsSummary.FoundVuln = true
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

	// Total summary
	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.BanditSummary.FoundVuln || outputJSON.Summary.SafetySummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.RetirejsSummary.FoundVuln || outputJSON.Summary.NpmAuditSummary.FoundVuln {
		outputJSON.Summary.TotalSummary.FoundVuln = true
	} else if outputJSON.Summary.GosecSummary.FoundInfo || outputJSON.Summary.BanditSummary.FoundInfo || outputJSON.Summary.SafetySummary.FoundInfo || outputJSON.Summary.BrakemanSummary.FoundInfo || outputJSON.Summary.RetirejsSummary.FoundInfo || outputJSON.Summary.NpmAuditSummary.FoundInfo {
		outputJSON.Summary.TotalSummary.FoundInfo = true
	}

	totalLow = outputJSON.Summary.RetirejsSummary.LowVuln + outputJSON.Summary.BrakemanSummary.LowVuln + outputJSON.Summary.SafetySummary.LowVuln + outputJSON.Summary.BanditSummary.LowVuln + outputJSON.Summary.GosecSummary.LowVuln + outputJSON.Summary.NpmAuditSummary.LowVuln
	totalMedium = outputJSON.Summary.RetirejsSummary.MediumVuln + outputJSON.Summary.BrakemanSummary.MediumVuln + outputJSON.Summary.SafetySummary.MediumVuln + outputJSON.Summary.BanditSummary.MediumVuln + outputJSON.Summary.GosecSummary.MediumVuln + outputJSON.Summary.NpmAuditSummary.MediumVuln
	totalHigh = outputJSON.Summary.RetirejsSummary.HighVuln + outputJSON.Summary.BrakemanSummary.HighVuln + outputJSON.Summary.SafetySummary.HighVuln + outputJSON.Summary.BanditSummary.HighVuln + outputJSON.Summary.GosecSummary.HighVuln + outputJSON.Summary.NpmAuditSummary.HighVuln

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

	if outputJSON.Summary.NpmAuditSummary.FoundVuln || outputJSON.Summary.NpmAuditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] JavaScript -> Npm Audit\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.NpmAuditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.NpmAuditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.NpmAuditSummary.LowVuln)
	}

	if outputJSON.Summary.TotalSummary.FoundVuln || outputJSON.Summary.TotalSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Total\n")
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.TotalSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.TotalSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.TotalSummary.LowVuln)
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

func printSTDOUTOutputRetireJS(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		if !strings.Contains(issue.Details, "doesn't have package.json") {
			fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
			fmt.Printf("[HUSKYCI][!] Version: %s\n", issue.Version)
			fmt.Printf("[HUSKYCI][!] Occurrences: %d\n", issue.Occurrences)
		}
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
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

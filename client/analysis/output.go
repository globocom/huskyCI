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

// printJSONOutput prints the analysis output in a JSON format
func printJSONOutput() error {
	jsonReady, err := json.Marshal(outputJSON)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonReady))
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

	// spotbugs
	printSTDOUTOutputSpotBugs(outputJSON.JavaResults.HuskyCISpotBugsOutput.LowVulns)
	printSTDOUTOutputSpotBugs(outputJSON.JavaResults.HuskyCISpotBugsOutput.MediumVulns)
	printSTDOUTOutputSpotBugs(outputJSON.JavaResults.HuskyCISpotBugsOutput.HighVulns)

	// tfsec
	printSTDOUTOutputTFSec(outputJSON.HclResults.HuskyCITFSecOutput.LowVulns)
	printSTDOUTOutputTFSec(outputJSON.HclResults.HuskyCITFSecOutput.MediumVulns)
	printSTDOUTOutputTFSec(outputJSON.HclResults.HuskyCITFSecOutput.HighVulns)

	// securitycodescan
	printSTDOUTOutputSecurityCodeScan(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.LowVulns)
	printSTDOUTOutputSecurityCodeScan(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.MediumVulns)
	printSTDOUTOutputSecurityCodeScan(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.HighVulns)

	printAllSummary(analysis)
}

// prepareAllSummary prepares how many low, medium and high vulnerabilites were found.
func prepareAllSummary(analysis types.Analysis) {
	var totalNoSec, totalLow, totalMedium, totalHigh int

	outputJSON.GoResults = analysis.HuskyCIResults.GoResults
	outputJSON.JavaScriptResults = analysis.HuskyCIResults.JavaScriptResults
	outputJSON.PythonResults = analysis.HuskyCIResults.PythonResults
	outputJSON.RubyResults = analysis.HuskyCIResults.RubyResults
	outputJSON.JavaResults = analysis.HuskyCIResults.JavaResults
	outputJSON.HclResults = analysis.HuskyCIResults.HclResults
	outputJSON.CSharpResults = analysis.HuskyCIResults.CSharpResults
	outputJSON.GenericResults = analysis.HuskyCIResults.GenericResults

	// GoSec summary
	outputJSON.Summary.GosecSummary.NoSecVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.NoSecVulns)
	outputJSON.Summary.GosecSummary.LowVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.LowVulns)
	outputJSON.Summary.GosecSummary.MediumVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.MediumVulns)
	outputJSON.Summary.GosecSummary.HighVuln = len(outputJSON.GoResults.HuskyCIGosecOutput.HighVulns)
	if len(outputJSON.GoResults.HuskyCIGosecOutput.LowVulns) > 0 || len(outputJSON.GoResults.HuskyCIGosecOutput.NoSecVulns) > 0 {
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
	if len(outputJSON.PythonResults.HuskyCIBanditOutput.LowVulns) > 0 || len(outputJSON.PythonResults.HuskyCIBanditOutput.NoSecVulns) > 0 {
		outputJSON.Summary.BanditSummary.FoundInfo = true
	}
	if len(outputJSON.PythonResults.HuskyCIBanditOutput.MediumVulns) > 0 || len(outputJSON.PythonResults.HuskyCIBanditOutput.HighVulns) > 0 {
		outputJSON.Summary.BanditSummary.FoundVuln = true
	}

	// Safety summary
	outputJSON.Summary.SafetySummary.LowVuln = len(outputJSON.PythonResults.HuskyCISafetyOutput.LowVulns)
	outputJSON.Summary.SafetySummary.MediumVuln = len(outputJSON.PythonResults.HuskyCISafetyOutput.MediumVulns)
	outputJSON.Summary.SafetySummary.HighVuln = len(outputJSON.PythonResults.HuskyCISafetyOutput.HighVulns)
	if len(outputJSON.PythonResults.HuskyCISafetyOutput.LowVulns) > 0 || len(outputJSON.PythonResults.HuskyCISafetyOutput.NoSecVulns) > 0 {
		outputJSON.Summary.SafetySummary.FoundInfo = true
	}
	if len(outputJSON.PythonResults.HuskyCISafetyOutput.MediumVulns) > 0 || len(outputJSON.PythonResults.HuskyCISafetyOutput.HighVulns) > 0 {
		outputJSON.Summary.SafetySummary.FoundVuln = true
	}

	// Brakeman summary
	outputJSON.Summary.BrakemanSummary.NoSecVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.NoSecVulns)
	outputJSON.Summary.BrakemanSummary.LowVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.LowVulns)
	outputJSON.Summary.BrakemanSummary.MediumVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.MediumVulns)
	outputJSON.Summary.BrakemanSummary.HighVuln = len(outputJSON.RubyResults.HuskyCIBrakemanOutput.HighVulns)
	if len(outputJSON.RubyResults.HuskyCIBrakemanOutput.LowVulns) > 0 || len(outputJSON.RubyResults.HuskyCIBrakemanOutput.NoSecVulns) > 0 {
		outputJSON.Summary.BrakemanSummary.FoundInfo = true
	}
	if len(outputJSON.RubyResults.HuskyCIBrakemanOutput.MediumVulns) > 0 || len(outputJSON.RubyResults.HuskyCIBrakemanOutput.HighVulns) > 0 {
		outputJSON.Summary.BrakemanSummary.FoundVuln = true
	}

	// NpmAudit summary
	outputJSON.Summary.NpmAuditSummary.LowVuln = len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns)
	outputJSON.Summary.NpmAuditSummary.MediumVuln = len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns)
	outputJSON.Summary.NpmAuditSummary.HighVuln = len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns)
	if len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.NoSecVulns) > 0 {
		outputJSON.Summary.NpmAuditSummary.FoundInfo = true
	}
	if len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns) > 0 {
		outputJSON.Summary.NpmAuditSummary.FoundVuln = true
	}

	// YarnAudit summary
	outputJSON.Summary.YarnAuditSummary.LowVuln = len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns)
	outputJSON.Summary.YarnAuditSummary.MediumVuln = len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns)
	outputJSON.Summary.YarnAuditSummary.HighVuln = len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns)
	if len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.NoSecVulns) > 0 {
		outputJSON.Summary.YarnAuditSummary.FoundInfo = true
	}
	if len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns) > 0 || len(outputJSON.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns) > 0 {
		outputJSON.Summary.YarnAuditSummary.FoundVuln = true
	}

	// SpotBugs summary
	outputJSON.Summary.SpotBugsSummary.LowVuln = len(outputJSON.JavaResults.HuskyCISpotBugsOutput.LowVulns)
	outputJSON.Summary.SpotBugsSummary.MediumVuln = len(outputJSON.JavaResults.HuskyCISpotBugsOutput.MediumVulns)
	outputJSON.Summary.SpotBugsSummary.HighVuln = len(outputJSON.JavaResults.HuskyCISpotBugsOutput.HighVulns)
	if len(outputJSON.JavaResults.HuskyCISpotBugsOutput.LowVulns) > 0 || len(outputJSON.JavaResults.HuskyCISpotBugsOutput.NoSecVulns) > 0 {
		outputJSON.Summary.SpotBugsSummary.FoundInfo = true
	}
	if len(outputJSON.JavaResults.HuskyCISpotBugsOutput.MediumVulns) > 0 || len(outputJSON.JavaResults.HuskyCISpotBugsOutput.HighVulns) > 0 {
		outputJSON.Summary.SpotBugsSummary.FoundVuln = true
	}

	// GitLeaks summary
	outputJSON.Summary.GitleaksSummary.NoSecVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.NoSecVulns)
	outputJSON.Summary.GitleaksSummary.LowVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.LowVulns)
	outputJSON.Summary.GitleaksSummary.MediumVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.MediumVulns)
	outputJSON.Summary.GitleaksSummary.HighVuln = len(outputJSON.GenericResults.HuskyCIGitleaksOutput.HighVulns)
	if len(outputJSON.GenericResults.HuskyCIGitleaksOutput.LowVulns) > 0 || len(outputJSON.GenericResults.HuskyCIGitleaksOutput.NoSecVulns) > 0 {
		outputJSON.Summary.GitleaksSummary.FoundInfo = true
	}
	if len(outputJSON.GenericResults.HuskyCIGitleaksOutput.MediumVulns) > 0 || len(outputJSON.GenericResults.HuskyCIGitleaksOutput.HighVulns) > 0 {
		outputJSON.Summary.GitleaksSummary.FoundVuln = true
	}

	// TFSec summary
	outputJSON.Summary.TFSecSummary.LowVuln = len(outputJSON.HclResults.HuskyCITFSecOutput.LowVulns)
	outputJSON.Summary.TFSecSummary.MediumVuln = len(outputJSON.HclResults.HuskyCITFSecOutput.MediumVulns)
	outputJSON.Summary.TFSecSummary.HighVuln = len(outputJSON.HclResults.HuskyCITFSecOutput.HighVulns)
	if len(outputJSON.HclResults.HuskyCITFSecOutput.LowVulns) > 0 || len(outputJSON.HclResults.HuskyCITFSecOutput.NoSecVulns) > 0 {
		outputJSON.Summary.TFSecSummary.FoundInfo = true
	}
	if len(outputJSON.HclResults.HuskyCITFSecOutput.MediumVulns) > 0 || len(outputJSON.HclResults.HuskyCITFSecOutput.HighVulns) > 0 {
		outputJSON.Summary.TFSecSummary.FoundVuln = true
	}

	// SecurityCodeScan summary
	outputJSON.Summary.SecurityCodeScanSummary.LowVuln = len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.LowVulns)
	outputJSON.Summary.SecurityCodeScanSummary.MediumVuln = len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.MediumVulns)
	outputJSON.Summary.SecurityCodeScanSummary.HighVuln = len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.HighVulns)
	if len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.LowVulns) > 0 || len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.NoSecVulns) > 0 {
		outputJSON.Summary.SecurityCodeScanSummary.FoundInfo = true
	}
	if len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.MediumVulns) > 0 || len(outputJSON.CSharpResults.HuskyCISecurityCodeScanOutput.HighVulns) > 0 {
		outputJSON.Summary.SecurityCodeScanSummary.FoundVuln = true
	}

	// Total summary
	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.BanditSummary.FoundVuln || outputJSON.Summary.SafetySummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.NpmAuditSummary.FoundVuln || outputJSON.Summary.YarnAuditSummary.FoundVuln || outputJSON.Summary.GitleaksSummary.FoundVuln || outputJSON.Summary.SpotBugsSummary.FoundVuln || outputJSON.Summary.TFSecSummary.FoundVuln || outputJSON.Summary.SecurityCodeScanSummary.FoundVuln {
		outputJSON.Summary.TotalSummary.FoundVuln = true
		types.FoundVuln = true
	} else if outputJSON.Summary.GosecSummary.FoundInfo || outputJSON.Summary.BanditSummary.FoundInfo || outputJSON.Summary.SafetySummary.FoundInfo || outputJSON.Summary.BrakemanSummary.FoundInfo || outputJSON.Summary.NpmAuditSummary.FoundInfo || outputJSON.Summary.YarnAuditSummary.FoundInfo || outputJSON.Summary.GitleaksSummary.FoundInfo || outputJSON.Summary.SpotBugsSummary.FoundInfo || outputJSON.Summary.TFSecSummary.FoundInfo || outputJSON.Summary.SecurityCodeScanSummary.FoundInfo {
		outputJSON.Summary.TotalSummary.FoundInfo = true
		types.FoundInfo = true
	}

	totalNoSec = outputJSON.Summary.BrakemanSummary.NoSecVuln + outputJSON.Summary.BanditSummary.NoSecVuln + outputJSON.Summary.GosecSummary.NoSecVuln + outputJSON.Summary.GitleaksSummary.NoSecVuln

	totalLow = outputJSON.Summary.BrakemanSummary.LowVuln + outputJSON.Summary.SafetySummary.LowVuln + outputJSON.Summary.BanditSummary.LowVuln + outputJSON.Summary.GosecSummary.LowVuln + outputJSON.Summary.NpmAuditSummary.LowVuln + outputJSON.Summary.YarnAuditSummary.LowVuln + outputJSON.Summary.GitleaksSummary.LowVuln + outputJSON.Summary.SpotBugsSummary.LowVuln + outputJSON.Summary.TFSecSummary.LowVuln + outputJSON.Summary.SecurityCodeScanSummary.LowVuln

	totalMedium = outputJSON.Summary.BrakemanSummary.MediumVuln + outputJSON.Summary.SafetySummary.MediumVuln + outputJSON.Summary.BanditSummary.MediumVuln + outputJSON.Summary.GosecSummary.MediumVuln + outputJSON.Summary.NpmAuditSummary.MediumVuln + outputJSON.Summary.YarnAuditSummary.MediumVuln + outputJSON.Summary.GitleaksSummary.MediumVuln + outputJSON.Summary.SpotBugsSummary.MediumVuln + outputJSON.Summary.TFSecSummary.MediumVuln + outputJSON.Summary.SecurityCodeScanSummary.MediumVuln

	totalHigh = outputJSON.Summary.BrakemanSummary.HighVuln + outputJSON.Summary.SafetySummary.HighVuln + outputJSON.Summary.BanditSummary.HighVuln + outputJSON.Summary.GosecSummary.HighVuln + outputJSON.Summary.NpmAuditSummary.HighVuln + outputJSON.Summary.YarnAuditSummary.HighVuln + outputJSON.Summary.GitleaksSummary.HighVuln + outputJSON.Summary.SpotBugsSummary.HighVuln + outputJSON.Summary.TFSecSummary.HighVuln + outputJSON.Summary.SecurityCodeScanSummary.HighVuln

	outputJSON.Summary.TotalSummary.HighVuln = totalHigh
	outputJSON.Summary.TotalSummary.MediumVuln = totalMedium
	outputJSON.Summary.TotalSummary.LowVuln = totalLow
	outputJSON.Summary.TotalSummary.NoSecVuln = totalNoSec

}

func printAllSummary(analysis types.Analysis) {

	var gosecVersion, banditVersion, safetyVersion, brakemanVersion, npmauditVersion, yarnauditVersion, gitleaksVersion, spotbugsVersion, tfsecVersion, securityCodeScanVersion string

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
		case "spotbugs":
			spotbugsVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "gitleaks":
			gitleaksVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "tfsec":
			tfsecVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		case "securitycodescan":
			securityCodeScanVersion = fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		}
	}

	if outputJSON.Summary.GosecSummary.FoundVuln || outputJSON.Summary.GosecSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Go -> %s\n", gosecVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.GosecSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.GosecSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.GosecSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.GosecSummary.NoSecVuln)
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
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.SafetySummary.NoSecVuln)
	}

	if outputJSON.Summary.BrakemanSummary.FoundVuln || outputJSON.Summary.BrakemanSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Ruby -> %s\n", brakemanVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.BrakemanSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.BrakemanSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.BrakemanSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.BrakemanSummary.NoSecVuln)
	}

	if outputJSON.Summary.NpmAuditSummary.FoundVuln || outputJSON.Summary.NpmAuditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] JavaScript -> %s\n", npmauditVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.NpmAuditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.NpmAuditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.NpmAuditSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.NpmAuditSummary.NoSecVuln)
	}

	if outputJSON.Summary.YarnAuditSummary.FoundVuln || outputJSON.Summary.YarnAuditSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] JavaScript -> %s\n", yarnauditVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.YarnAuditSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.YarnAuditSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.YarnAuditSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.YarnAuditSummary.NoSecVuln)
	}

	if outputJSON.Summary.SpotBugsSummary.FoundVuln || outputJSON.Summary.SpotBugsSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Java -> %s\n", spotbugsVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.SpotBugsSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.SpotBugsSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.SpotBugsSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.SpotBugsSummary.NoSecVuln)
	}

	if outputJSON.Summary.TFSecSummary.FoundVuln || outputJSON.Summary.TFSecSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] HCL -> %s\n", tfsecVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.TFSecSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.TFSecSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.TFSecSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.TFSecSummary.NoSecVuln)
	}

	if outputJSON.Summary.SecurityCodeScanSummary.FoundVuln || outputJSON.Summary.SecurityCodeScanSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] C# -> %s\n", securityCodeScanVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.SecurityCodeScanSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.SecurityCodeScanSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.SecurityCodeScanSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.SecurityCodeScanSummary.NoSecVuln)
	}

	if outputJSON.Summary.GitleaksSummary.FoundVuln || outputJSON.Summary.GitleaksSummary.FoundInfo {
		fmt.Println()
		fmt.Printf("[HUSKYCI][SUMMARY] Generic -> %s\n", gitleaksVersion)
		fmt.Printf("[HUSKYCI][SUMMARY] High: %d\n", outputJSON.Summary.GitleaksSummary.HighVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Medium: %d\n", outputJSON.Summary.GitleaksSummary.MediumVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] Low: %d\n", outputJSON.Summary.GitleaksSummary.LowVuln)
		fmt.Printf("[HUSKYCI][SUMMARY] NoSecHusky: %d\n", outputJSON.Summary.GitleaksSummary.NoSecVuln)
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
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
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
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
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
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
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
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
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
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
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
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
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

func printSTDOUTOutputSpotBugs(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Confidence: %s\n", issue.Confidence)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		fmt.Printf("[HUSKYCI][!] Type: %s\n", issue.Type)
	}
}

func printSTDOUTOutputTFSec(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
	}
}

func printSTDOUTOutputGitleaks(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
		fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
	}
}

func printSTDOUTOutputSecurityCodeScan(issues []types.HuskyCIVulnerability) {
	for _, issue := range issues {
		fmt.Println()
		fmt.Printf("[HUSKYCI][!] Title: %s\n", issue.Title)
		fmt.Printf("[HUSKYCI][!] Language: %s\n", issue.Language)
		fmt.Printf("[HUSKYCI][!] Tool: %s\n", issue.SecurityTool)
		fmt.Printf("[HUSKYCI][!] Severity: %s\n", issue.Severity)
		fmt.Printf("[HUSKYCI][!] Details: %s\n", issue.Details)
		if !strings.Contains(issue.Details, "could not run 'security-scan' on your project") {
			fmt.Printf("[HUSKYCI][!] File: %s\n", issue.File)
			fmt.Printf("[HUSKYCI][!] Line: %s\n", issue.Line)
			fmt.Printf("[HUSKYCI][!] Code: %s\n", issue.Code)
		}
	}
}

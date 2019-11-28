// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sonarqube

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/globocom/huskyCI/client/types"
	"github.com/globocom/huskyCI/client/util"
)

const goContainerBasePath = `/go/src/code/`
const placeholderFileName = "huskyCI_Placeholder_File"
const placeholderFileText = `
Placeholder file indicating that no file was associated with this vulnerability.
This usually means that the vulnerability is related to a missing file
or is not associated with any specific file, i.e.: vulnerable dependency versions.
`

// GenerateOutputFile prints the analysis output in a JSON format
func GenerateOutputFile(analysis types.Analysis, outputPath, outputFileName string) error {

	allVulns := make([]types.HuskyCIVulnerability, 0)

	// gosec
	allVulns = append(allVulns, analysis.HuskyCIResults.GoResults.HuskyCIGosecOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.GoResults.HuskyCIGosecOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.GoResults.HuskyCIGosecOutput.HighVulns...)

	// bandit
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCIBanditOutput.NoSecVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCIBanditOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCIBanditOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCIBanditOutput.HighVulns...)

	// safety
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCISafetyOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCISafetyOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.PythonResults.HuskyCISafetyOutput.HighVulns...)

	// brakeman
	allVulns = append(allVulns, analysis.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.HighVulns...)

	// npmaudit
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns...)

	// yarnaudit
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaScriptResults.HuskyCIYarnAuditOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaScriptResults.HuskyCIYarnAuditOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaScriptResults.HuskyCIYarnAuditOutput.HighVulns...)

	// gitleaks
	allVulns = append(allVulns, analysis.HuskyCIResults.GenericResults.HuskyCIGitleaksOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.GenericResults.HuskyCIGitleaksOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.GenericResults.HuskyCIGitleaksOutput.HighVulns...)

	// spotbugs
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaResults.HuskyCISpotBugsOutput.LowVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaResults.HuskyCISpotBugsOutput.MediumVulns...)
	allVulns = append(allVulns, analysis.HuskyCIResults.JavaResults.HuskyCISpotBugsOutput.HighVulns...)

	var sonarOutput HuskyCISonarOutput
	sonarOutput.Issues = make([]SonarIssue, 0)

	for _, vuln := range allVulns {
		var issue SonarIssue
		issue.EngineID = "huskyCI"
		issue.Type = "VULNERABILITY"
		issue.RuleID = vuln.Language + " - " + vuln.SecurityTool
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
		if vuln.File == "" {
			err := util.CreateFile([]byte(placeholderFileText), outputPath, placeholderFileName)
			if err != nil {
				return err
			}
			issue.PrimaryLocation.FilePath = filepath.Join(outputPath, placeholderFileName)
		} else {
			var filePath string
			if vuln.Language == "Go" {
				filePath = strings.Replace(vuln.File, goContainerBasePath, "", 1)
			} else {
				filePath = vuln.File
			}
			issue.PrimaryLocation.FilePath = filePath
		}
		issue.PrimaryLocation.Message = vuln.Details
		issue.PrimaryLocation.TextRange.StartLine = 1
		lineNum, err := strconv.Atoi(vuln.Line)
		if err != nil {
			lineNum = 1
		}
		if lineNum != 1 && lineNum > 0 {
			issue.PrimaryLocation.TextRange.StartLine = lineNum
		}
		sonarOutput.Issues = append(sonarOutput.Issues, issue)
	}
	sonarOutputString, err := json.Marshal(sonarOutput)
	if err != nil {
		return err
	}
	err = util.CreateFile(sonarOutputString, outputPath, outputFileName)
	if err != nil {
		return err
	}

	return nil
}

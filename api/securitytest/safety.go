// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/util"
	"github.com/globocom/huskyCI/api/vulnerability"
)

// SafetyOutput is the struct that holds issues, messages and errors found on a Safety scan.
type SafetyOutput struct {
	SafetyIssues   []SafetyIssue `json:"issues"`
	ReqNotFound    bool
	WarningFound   bool
	OutputWarnings []string
}

// SafetyIssue is a struct that holds the results that were scanned and the file they came from.
type SafetyIssue struct {
	Dependency string `json:"dependency"`
	Below      string `json:"vulnerable_below"`
	Version    string `json:"installed_version"`
	Comment    string `json:"description"`
	ID         string `json:"id"`
}

func (s *SecurityTest) analyzeSafety() error {

	safetyOutput := SafetyOutput{}

	// requirements.txt not found
	if s.WarningFound != "" {
		s.prepareSafetyVulns(safetyOutput)
		return nil
	}

	warningFound := strings.Contains(s.Container.Output, "Warning: unpinned requirement ")

	// check if warnings were found and handle container output
	if warningFound {
		outputJSON := util.GetLastLine(s.Container.Output)
		safetyOutput.OutputWarnings = util.GetAllLinesButLast(s.Container.Output)
		s.Container.Output = outputJSON
	}

	cOutputSanitized := sanitizeSafetyJSON(s.Container.Output)
	s.Container.Output = cOutputSanitized

	// Unmarshall rawOutput into finalOutput, that is a Safety struct.
	if err := json.Unmarshal([]byte(s.Container.Output), &safetyOutput); err != nil {
		log.Error("analyzeSafety", "SAFETY", 1018, s.Container.Output, err)
		s.Result = "error"
		s.Info = log.MsgCode[1018]
		s.ErrorFound = err.Error()
		return err
	}

	s.prepareSafetyVulns(safetyOutput)

	return nil
}

func (s *SecurityTest) prepareSafetyVulns(safetyOutput SafetyOutput) {

	// requirements.txt not found
	if s.WarningFound != "" {

		safetyVuln := vulnerability.New()

		safetyVuln.Language = "Python"
		safetyVuln.SecurityTest = "Safety"
		safetyVuln.Severity = "low"
		safetyVuln.Details = s.WarningFound

		s.Vulnerabilities = append(s.Vulnerabilities, *safetyVuln)

		return
	}

	onlyWarning := false

	// requiments.txt contains values that are not pinned or similiar
	if safetyOutput.WarningFound {
		if len(safetyOutput.SafetyIssues) == 0 {
			onlyWarning = true
		}
		for _, warning := range safetyOutput.OutputWarnings {

			safetyVuln := vulnerability.New()

			safetyVuln.Language = "Python"
			safetyVuln.SecurityTest = "Safety"
			safetyVuln.Severity = "low"
			safetyVuln.Details = adjustWarningMessage(warning)

			s.Vulnerabilities = append(s.Vulnerabilities, *safetyVuln)

		}
		if onlyWarning {
			return
		}
	}

	for _, issue := range safetyOutput.SafetyIssues {

		safetyVuln := vulnerability.New()

		safetyVuln.Language = "Python"
		safetyVuln.SecurityTest = "Safety"
		safetyVuln.Severity = "high"
		safetyVuln.Details = issue.Comment
		safetyVuln.Code = issue.Dependency + " " + issue.Version
		safetyVuln.VunerableBelow = issue.Below

		s.Vulnerabilities = append(s.Vulnerabilities, *safetyVuln)
	}

}

// sanitizeSafetyJSON returns a sanitized string from Safety container logs.
// Safety might return a JSON with the "\" and "\"" characters, which needs to be sanitized to be unmarshalled correctly.
func sanitizeSafetyJSON(s string) string {
	if s == "" {
		return ""
	}
	s1 := strings.Replace(s, "\\", "\\\\", -1)
	s2 := strings.Replace(s1, "\\\"", "\\\\\"", -1)
	return s2
}

// adjustWarningMessage returns the Safety Warning string that will be printed.
func adjustWarningMessage(warningRaw string) string {
	warning := strings.Split(warningRaw, ":")
	if len(warning) > 1 {
		warning[1] = strings.Replace(warning[1], "safety_huskyci_analysis_requirements_raw.txt", "'requirements.txt'", -1)
		warning[1] = strings.Replace(warning[1], " unpinned", "Unpinned", -1)

		return (warning[1] + " huskyCI can check it if you pin it in a format such as this: \"mypacket==3.2.9\" :D")
	}

	return warningRaw
}

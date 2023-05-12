// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// NpmAuditOutput is the struct that stores all npm audit output
type NpmAuditOutput struct {
	Vulnerabilities      map[string]Vulnerability `json:"vulnerabilities"`
	Metadata        	 Metadata                 `json:"metadata"`
	PackageNotFound 	 bool
}

// Vulnerability is the granular output of a security info found
type Vulnerability struct {
	Via                []ViaMessage     `json:"via"`
	ID                 int              `json:"id"`
	Name         	   string           `json:"name"`
	VulnerableVersions string           `json:"range"`
	Severity           string           `json:"severity"`
	FixAvailable       FixAvailableType `json:"fixAvailable"`
	Title              string           `json:"title"`
}

type FixAvailableType struct {
	Text string
}

// FixAvailableType holds the information of the dependency that originated the vulnerability
type FixAvailableTypeNPM struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	IsSemVerMajor bool   `json:"isSemVerMajor"`
}

// Via holds the information of a given security issue found
type Via struct {
	Source     int        `json:"source"`
	Name	   string     `json:"name"`
	Dependency string     `json:"dependency"`
	Title      string     `json:"title"`
	Url        string     `json:"url"`
	Severity   string     `json:"severity"`
	CWE        []string   `json:"cwe"`
	CVSS       CVSSType   `json:"cvss"`
	Range      string     `json:"range"`
}

// CVSSType is the struct that holds CVSS info
type CVSSType struct {
	Score        json.Number `json:"score"`
	VectorString string      `json:"vectorString"`
}

// Metadata is the struct that holds vulnerabilities summary
type Metadata struct {
	Vulnerabilities VulnerabilitiesSummary `json:"vulnerabilities"`
}

// VulnerabilitiesSummary is the struct that has all types of possible vulnerabilities from npm audit
type VulnerabilitiesSummary struct {
	Info     int `json:"info"`
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type ViaMessage struct {
    Text string
}

func (e *ViaMessage) UnmarshalJSON(data []byte) error {
    if len(data) == 0 || string(data) == "null" {
        return nil
    }
    if data[0] == '"' && data[len(data)-1] == '"' {
        return json.Unmarshal(data, &e.Text)
    }
    if data[0] == '{' && data[len(data)-1] == '}' {
		tmp := Via{}
        if err := json.Unmarshal(data, &tmp); err != nil {
			return err
		}
		e.Text = ""
		e.Text += fmt.Sprintf("\tSource: %d\n", tmp.Source)
		e.Text += fmt.Sprintf("\tName: %s\n", tmp.Name)
		e.Text += fmt.Sprintf("\tDependency: %s\n", tmp.Dependency)
		e.Text += fmt.Sprintf("\tTitle: %s\n", tmp.Title)
		e.Text += fmt.Sprintf("\tUrl: %s\n", tmp.Url)
		e.Text += fmt.Sprintf("\tSeverity: %s\n", tmp.Severity)
		e.Text += "\tCWEs: "
		for _, cwe := range tmp.CWE {
			e.Text += fmt.Sprintf("%s, ", cwe)
		}
		e.Text += "\n"
		e.Text += fmt.Sprintf("\tCVSS: %s (%s)\n", tmp.CVSS.Score.String(), tmp.CVSS.VectorString)
		e.Text += fmt.Sprintf("\tVersion Range: %s\n", tmp.Range)
		return nil
    }
    return fmt.Errorf("unsupported Via field")
}


func (e *FixAvailableType) UnmarshalJSON(data []byte) error {
    if len(data) == 0 || string(data) == "null" {
        return nil
    }
	if string(data) == "false" {
		e.Text = "false"
        return nil
    }
	if string(data) == "true" {
		e.Text = "true"
        return nil
    }
    if data[0] == '"' && data[len(data)-1] == '"' {
        return json.Unmarshal(data, &e.Text)
    }
    if data[0] == '{' && data[len(data)-1] == '}' {
		tmp := FixAvailableTypeNPM{}
        if err := json.Unmarshal(data, &tmp); err != nil {
			return err
		}
		e.Text = fmt.Sprintf("Fix available: %s %s", tmp.Name, tmp.Version)
		return nil
    }
    return fmt.Errorf("unsupported fixAvailable field")
}

func analyzeNpmaudit(npmAuditScan *SecTestScanInfo) error {

	npmAuditOutput := NpmAuditOutput{}
	npmAuditScan.FinalOutput = npmAuditOutput

	// if package-lock was not found, a warning will be genrated as a low vuln
	packageNotFound := strings.Contains(npmAuditScan.Container.COutput, "ERROR_PACKAGE_LOCK_NOT_FOUND")
	if packageNotFound {
		npmAuditScan.PackageNotFound = true
		npmAuditScan.prepareNpmAuditVulns()
		npmAuditScan.prepareContainerAfterScan()
		return nil
	}

	// nil cOutput states that no Issues were found.
	if npmAuditScan.Container.COutput == "" {
		npmAuditScan.prepareContainerAfterScan()
		return nil
	}

	// Unmarshall rawOutput into finalOutput, that is a NpmAuditOutput struct.
	if err := json.Unmarshal([]byte(npmAuditScan.Container.COutput), &npmAuditOutput); err != nil {
		log.Error("analyzeNpmaudit", "NPMAUDIT", 1014, npmAuditScan.Container.COutput, err)
		npmAuditScan.ErrorFound = util.HandleScanError(npmAuditScan.Container.COutput, err)
		npmAuditScan.prepareContainerAfterScan()
		return npmAuditScan.ErrorFound
	}
	npmAuditScan.FinalOutput = npmAuditOutput

	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	npmAuditScan.prepareNpmAuditVulns()
	npmAuditScan.prepareContainerAfterScan()
	return nil
}

func (npmAuditScan *SecTestScanInfo) prepareNpmAuditVulns() {

	huskyCInpmauditResults := types.HuskyCISecurityTestOutput{}
	npmAuditOutput := npmAuditScan.FinalOutput.(NpmAuditOutput)

	if npmAuditScan.PackageNotFound {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Severity = "low"
		npmauditVuln.Title = "No package-lock.json found."
		npmauditVuln.Details = "It looks like your project doesn't have a package-lock.json file. If you use NPM to handle your dependencies, it would be a good idea to commit it so huskyCI can check for vulnerabilities."

		npmAuditScan.Vulnerabilities.LowVulns = append(npmAuditScan.Vulnerabilities.LowVulns, npmauditVuln)
		return
	}

	for _, issue := range npmAuditOutput.Vulnerabilities {
		npmauditVuln := types.HuskyCIVulnerability{}
		npmauditVuln.Language = "JavaScript"
		npmauditVuln.SecurityTool = "NpmAudit"
		npmauditVuln.Title = fmt.Sprintf("Vulnerable Dependency: %s %s (%s)", issue.Name, issue.VulnerableVersions, issue.Title)
		if issue.FixAvailable.Text != "true" && issue.FixAvailable.Text != "false" {
			npmauditVuln.Details = issue.FixAvailable.Text
		}
		npmauditVuln.VunerableBelow = issue.VulnerableVersions
		npmauditVuln.Code = issue.Name
		npmauditVuln.Version = ""
		for i, via := range issue.Via {
			npmauditVuln.Version += fmt.Sprintf("Advisories and information (Via %d):\n", i)
			npmauditVuln.Version += fmt.Sprintf("%s\n", via.Text)
		}

		switch issue.Severity {
		case "info", "low":
			npmauditVuln.Severity = "low"
			huskyCInpmauditResults.LowVulns = append(huskyCInpmauditResults.LowVulns, npmauditVuln)
		case "moderate":
			npmauditVuln.Severity = "medium"
			huskyCInpmauditResults.MediumVulns = append(huskyCInpmauditResults.MediumVulns, npmauditVuln)
		case "high", "critical":
			npmauditVuln.Severity = "high"
			huskyCInpmauditResults.HighVulns = append(huskyCInpmauditResults.HighVulns, npmauditVuln)
		}

	}

	npmAuditScan.Vulnerabilities = huskyCInpmauditResults
}

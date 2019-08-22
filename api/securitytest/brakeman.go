// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/db"
	huskydocker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// BrakemanScan holds all information needed for a Brakeman scan.
type BrakemanScan struct {
	RID             string
	CID             string
	URL             string
	Branch          string
	Image           string
	Command         string
	RawOutput       string
	ErrorFound      error
	FinalOutput     BrakemanOutput
	Vulnerabilities types.HuskyCISecurityTestOutput
}

// BrakemanOutput is the struct that holds issues and stats found on a Brakeman scan.
type BrakemanOutput struct {
	Warnings []WarningItem `json:"warnings"`
}

// WarningItem is the struct that holds all detailed information of a vulnerability found.
type WarningItem struct {
	Type       string `json:"warning_type"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Details    string `json:"link"`
	Confidence string `json:"confidence"`
}

func newScanBrakeman(URL, branch, command string) BrakemanScan {
	return BrakemanScan{
		Image:   "huskyci/brakeman",
		URL:     URL,
		Branch:  branch,
		Command: util.HandleCmd(URL, branch, command),
	}
}

func initBrakeman(enryScan EnryScan, allScansResult *AllScansResult) error {
	brakemanScan, brakemanContainer, err := runScanBrakeman(enryScan.URL, enryScan.Branch)
	if err != nil {
		return err
	}

	for _, highVuln := range brakemanScan.Vulnerabilities.HighVulns {
		allScansResult.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.HighVulns = append(allScansResult.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.HighVulns, highVuln)
	}
	for _, mediumVuln := range brakemanScan.Vulnerabilities.MediumVulns {
		allScansResult.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.MediumVulns = append(allScansResult.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.MediumVulns, mediumVuln)
	}
	for _, lowVuln := range brakemanScan.Vulnerabilities.LowVulns {
		allScansResult.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.LowVulns = append(allScansResult.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.LowVulns, lowVuln)
	}

	allScansResult.FinalResult = brakemanContainer.CResult
	allScansResult.Status = brakemanContainer.CStatus
	allScansResult.Containers = append(allScansResult.Containers, brakemanContainer)
	return nil
}

func runScanBrakeman(URL, branch string) (BrakemanScan, types.Container, error) {
	brakemanScan := BrakemanScan{}
	brakemanContainer, err := newContainerBrakeman()
	if err != nil {
		log.Error("runScanBrakeman", "BRAKEMAN", 1029, err)
		return brakemanScan, brakemanContainer, err
	}
	brakemanScan = newScanBrakeman(URL, branch, brakemanContainer.SecurityTest.Cmd)
	if err := brakemanScan.startBrakeman(); err != nil {
		return brakemanScan, brakemanContainer, err
	}

	brakemanScan.prepareContainerAfterScanBrakeman(&brakemanContainer)
	return brakemanScan, brakemanContainer, nil
}

func (brakemanScan *BrakemanScan) startBrakeman() error {
	if err := brakemanScan.dockerRunBrakeman(); err != nil {
		brakemanScan.ErrorFound = err
		return err
	}
	if err := brakemanScan.analyzeBrakeman(); err != nil {
		brakemanScan.ErrorFound = err
		return err
	}
	return nil
}

func (brakemanScan *BrakemanScan) dockerRunBrakeman() error {
	CID, cOutput, err := huskydocker.DockerRun(brakemanScan.Image, brakemanScan.Command)
	if err != nil {
		return err
	}
	brakemanScan.CID = CID
	brakemanScan.RawOutput = cOutput
	return nil
}

func (brakemanScan *BrakemanScan) analyzeBrakeman() error {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(brakemanScan.RawOutput, "ERROR_CLONING")
	if errorClonning {
		errorMsg := errors.New("error clonning")
		log.Error("analyzeBrakeman", "BRAKEMAN", 1031, brakemanScan.URL, brakemanScan.Branch, errorMsg)
		return errorMsg
	}

	// step 2: nil cOutput states that no Issues were found.
	if brakemanScan.RawOutput == "" {
		return nil
	}

	// step 3: Unmarshall rawOutput into finalOutput, that is a Brakeman struct.
	if err := json.Unmarshal([]byte(brakemanScan.RawOutput), &brakemanScan.FinalOutput); err != nil {
		log.Error("analyzeBrakeman", "BRAKEMAN", 1005, brakemanScan.RawOutput, err)
		return err
	}

	// step 4: find Issues that have severity "MEDIUM" or "HIGH" and confidence "HIGH".
	brakemanScan.prepareBrakemanOutput(brakemanScan.FinalOutput)
	return nil
}

func (brakemanScan *BrakemanScan) prepareBrakemanOutput(brakemanOutput BrakemanOutput) {
	huskyCIbrakemanResults := types.HuskyCISecurityTestOutput{}

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

		switch brakemanVuln.Confidence {
		case "High":
			huskyCIbrakemanResults.HighVulns = append(huskyCIbrakemanResults.HighVulns, brakemanVuln)
		case "Medium":
			huskyCIbrakemanResults.MediumVulns = append(huskyCIbrakemanResults.MediumVulns, brakemanVuln)
		case "Low":
			huskyCIbrakemanResults.LowVulns = append(huskyCIbrakemanResults.LowVulns, brakemanVuln)
		}
	}

	brakemanScan.Vulnerabilities = huskyCIbrakemanResults
}

func (brakemanScan *BrakemanScan) prepareContainerAfterScanBrakeman(brakemanContainer *types.Container) {
	if len(brakemanScan.Vulnerabilities.MediumVulns) > 0 || len(brakemanScan.Vulnerabilities.HighVulns) > 0 {
		brakemanContainer.CInfo = "Issues found."
		brakemanContainer.CResult = "failed"
	} else if len(brakemanScan.Vulnerabilities.LowVulns) > 0 && (len(brakemanScan.Vulnerabilities.MediumVulns) == 0 || len(brakemanScan.Vulnerabilities.HighVulns) == 0) {
		brakemanContainer.CInfo = "Warnings found."
		brakemanContainer.CResult = "passed"
	}
	brakemanContainer.CStatus = "finished"
	brakemanContainer.CID = brakemanScan.CID
	brakemanContainer.COutput = brakemanScan.RawOutput
	brakemanContainer.FinishedAt = time.Now()
}

func newContainerBrakeman() (types.Container, error) {
	brakemanContainer := types.Container{}
	brakemanQuery := map[string]interface{}{"name": "brakeman"}
	brakemanSecurityTest, err := db.FindOneDBSecurityTest(brakemanQuery)
	if err != nil {
		log.Error("newContainerBrakeman", "BRAKEMAN", 2012, err)
		return brakemanContainer, err
	}
	return types.Container{
		SecurityTest: brakemanSecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

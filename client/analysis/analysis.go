// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/globocom/huskyCI/client/config"
	"github.com/globocom/huskyCI/client/types"
	"github.com/globocom/huskyCI/client/util"
)

// StartAnalysis starts a container and returns its RID and error.
func StartAnalysis() (string, error) {

	// preparing POST to HuskyCI
	requestPayload := types.JSONPayload{
		RepositoryURL:    config.RepositoryURL,
		RepositoryBranch: config.RepositoryBranch,
		InternalDepURL:   config.InternalDepURL,
	}
	marshalPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}
	huskyStartAnalysisURL := config.HuskyAPI + "/analysis"

	httpClient, err := util.NewClient(config.HuskyUseTLS)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", huskyStartAnalysisURL, bytes.NewBuffer(marshalPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("huskyCIToken", config.HuskyCIToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		errorMsg := fmt.Sprintf("Error sending request to start analysis! StatusCode received: %d", resp.StatusCode)
		return "", errors.New(errorMsg)
	}

	defer resp.Body.Close()
	RID := resp.Header.Get("X-Request-Id")

	if RID == "" {
		errorMsg := fmt.Sprintf("Error sending request to start analysis. RID is empty!")
		return "", errors.New(errorMsg)
	}

	// Setting analysis values on the JSON output
	outputJSON.Summary.URL = requestPayload.RepositoryURL
	outputJSON.Summary.Branch = requestPayload.RepositoryBranch
	outputJSON.Summary.RID = RID

	return RID, nil
}

// GetAnalysis gets the results of an analysis.
func GetAnalysis(RID string) (types.Analysis, error) {

	analysis := types.Analysis{}

	httpClient, err := util.NewClient(config.HuskyUseTLS)
	if err != nil {
		return analysis, err
	}

	resp, err := httpClient.Get(config.HuskyAPI + "/analysis/" + RID)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return analysis, err
	}

	err = json.Unmarshal(body, &analysis)
	if err != nil {
		return analysis, err
	}

	return analysis, nil
}

// MonitorAnalysis will keep monitoring an analysis until it has finished or timed out.
func MonitorAnalysis(RID string) (types.Analysis, error) {

	analysis := types.Analysis{}
	timeout := time.After(15 * time.Minute)
	retryTick := time.Tick(60 * time.Second)

	for {
		select {
		case <-timeout:
			return analysis, errors.New("time out")
		case <-retryTick:
			analysis, err := GetAnalysis(RID)
			if err != nil {
				return analysis, err
			}
			if analysis.Status == "finished" {
				return analysis, nil
			}
			if !types.IsJSONoutput {
				fmt.Println("[HUSKYCI][!] Hold on! huskyCI is still running...")
			}
		}
	}
}

// PrepareResults analyzes the result received from HuskyCI API.
func PrepareResults(analysisResult types.Analysis) {
	for _, container := range analysisResult.Containers {
		prepareSecurityTestResult(container)
	}
}

// PrintResults prints huskyCI output either in JSON or the standard output.
func PrintResults(formatOutput string) error {

	prepareAllSummary()

	if types.IsJSONoutput {
		err := printJSONOutput()
		if err != nil {
			return err
		}
	} else {
		printSTDOUTOutput()
	}

	return nil
}

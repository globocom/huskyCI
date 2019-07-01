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
	resp, err := httpClient.Post(huskyStartAnalysisURL, "application/json", bytes.NewBuffer(marshalPayload))
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
	retryTick := time.Tick(20 * time.Second)

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
			fmt.Println("[HUSKYCI][!] Hold on! HuskyCI is still running...")
		}
	}
}

// AnalyzeResult analyzes the result received from HuskyCI API.
func AnalyzeResult(analysisResult types.Analysis) {
	fmt.Println()
	for _, container := range analysisResult.Containers {
		CheckMongoDBContainerOutput(container)
	}
}

// PrintResults prints huskyCI output either in JSON or the standard output.
func PrintResults(formatOutput string) error {
	if formatOutput == "JSON" {
		err := printJSONOutput()
		if err != nil {
			return err
		}
	} else {
		printhuskyCIOutput()
	}

	return nil
}

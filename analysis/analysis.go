// Copyright 2018 Globo.com authors. All rights reserved.
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

	"github.com/globocom/husky-client/config"
	"github.com/globocom/husky-client/types"
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
	huskyStartAnalysisURL := config.HuskyAPI + "/husky"
	req, err := http.NewRequest("POST", huskyStartAnalysisURL, bytes.NewBuffer(marshalPayload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// sending POST to HuskyCI
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// analyzing response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	responsePayload := types.JSONResponse{}
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return "", err
	}
	return responsePayload.RID, nil
}

// GetAnalysis gets the results of an analysis.
func GetAnalysis(RID string) (types.Analysis, error) {

	analysis := types.Analysis{}

	resp, err := http.Get(config.HuskyAPI + "/husky/" + RID)
	if err != nil {
		return analysis, err
	}
	defer resp.Body.Close()

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
		CheckContainerOutput(container)
	}
}

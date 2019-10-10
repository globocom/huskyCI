// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/globocom/huskyCI/client/types"
)

// Start starts a container at huskyci API and returns its RID and error.
func (hcli *Client) Start(RepositoryURL, RepositoryBranch string) (types.Analysis, error) {

	requestPayload := types.JSONPayload{
		RepositoryURL:    RepositoryURL,
		RepositoryBranch: RepositoryBranch,
	}

	marshalPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return types.Analysis{}, err
	}

	// Make request
	huskyStartAnalysisURL := hcli.target.Endpoint + "/analysis"
	req, err := http.NewRequest("POST", huskyStartAnalysisURL, bytes.NewBuffer(marshalPayload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Husky-Token", hcli.target.Token)
	if err != nil {
		return types.Analysis{}, err
	}

	resp, err := hcli.httpCli.Do(req)
	if err != nil {
		return types.Analysis{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		if resp.StatusCode == 401 {
			return types.Analysis{}, fmt.Errorf("Unauthorized Husky-Token: (%s)", hcli.target.Token)
		}
		return types.Analysis{}, fmt.Errorf("Error sending request to start analysis! StatusCode received: %d", resp.StatusCode)
	}

	RID := resp.Header.Get("X-Request-Id")
	if RID == "" {
		return types.Analysis{}, fmt.Errorf("Error sending request to start analysis. RID is empty")
	}

	// Setting analysis values on the JSON output
	analysis := types.Analysis{}
	analysis.Branch = requestPayload.RepositoryBranch
	analysis.URL = requestPayload.RepositoryURL
	analysis.RID = RID

	return analysis, nil
}

// Monitor will keep monitoring an analysis until it has finished or timed out.
func (hcli *Client) Monitor(RID string, timeoutMonitor, retry time.Duration) (types.Analysis, error) {

	analysisResult := types.Analysis{}
	timeout := time.After(timeoutMonitor)
	retryTick := time.Tick(retry)

	for {
		select {
		case <-timeout:
			return analysisResult, fmt.Errorf("huskyCI monitor timeout (%s)", timeoutMonitor.String())
		case <-retryTick:
			analysisResult, err := hcli.Get(RID)
			if err != nil {
				return analysisResult, err
			}
			if analysisResult.Status == "finished" {
				return analysisResult, nil
			} else if analysisResult.Status == "error running" {
				return analysisResult, fmt.Errorf("huskyCI encountered an error trying to execute this analysis: %v", analysisResult.ErrorFound)
			}
			if !types.IsJSONoutput {
				fmt.Println("[HUSKYCI][!] Hold on! huskyCI is still running...")
			}
		}
	}
}

// Get gets the results of an analysis.
func (hcli *Client) Get(RID string) (types.Analysis, error) {

	getAnalysisURL := hcli.target.Endpoint + "/analysis/" + RID
	req, err := http.NewRequest("GET", getAnalysisURL, nil)
	req.Header.Add("Husky-Token", hcli.target.Token)
	if err != nil {
		return types.Analysis{}, err
	}

	resp, err := hcli.httpCli.Do(req)
	if err != nil {
		return types.Analysis{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Analysis{}, fmt.Errorf("error trying to get this analysis: status code %v (%v)", resp.StatusCode, resp.Body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Analysis{}, err
	}

	analysis := types.Analysis{}
	err = json.Unmarshal(body, &analysis)
	if err != nil {
		return types.Analysis{}, err
	}

	return analysis, nil
}

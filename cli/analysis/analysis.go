// Copyright © 2019 Globo.com
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package analysis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/globocom/huskyCI/client/types"
	"github.com/spf13/viper"
)

// creates a custom httpClient
func createClient() *http.Client {
	// Setting custom HTTP client with timeouts
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Second,
	}
	var netClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}

	return netClient
}

// Start starts a container and returns its RID and error.
func Start(currentTarget types.Target) (types.JSONOutput, error) {

	requestPayload := types.JSONPayload{
		RepositoryURL:    viper.GetString("repo_url"),
		RepositoryBranch: viper.GetString("repo_branch"),
	}

	marshalPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return types.JSONOutput{}, err
	}

	netClient := createClient()

	// Make request
	huskyStartAnalysisURL := currentTarget.Endpoint + "/analysis"
	req, err := http.NewRequest("POST", huskyStartAnalysisURL, bytes.NewBuffer(marshalPayload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Husky-Token", currentTarget.Token)
	if err != nil {
		return types.JSONOutput{}, err
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return types.JSONOutput{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		if resp.StatusCode == 401 {
			return types.JSONOutput{}, fmt.Errorf("Unauthorized Husky-Token: (%s)", currentTarget.Token)
		}
		return types.JSONOutput{}, fmt.Errorf("Error sending request to start analysis! StatusCode received: %d", resp.StatusCode)
	}

	RID := resp.Header.Get("X-Request-Id")
	if RID == "" {
		return types.JSONOutput{}, fmt.Errorf("Error sending request to start analysis. RID is empty")
	}

	// Setting analysis values on the JSON output
	var analisysResult types.JSONOutput
	analisysResult.Summary.URL = requestPayload.RepositoryURL
	analisysResult.Summary.Branch = requestPayload.RepositoryBranch
	analisysResult.Summary.RID = RID

	return analisysResult, nil
}

// Monitor will keep monitoring an analysis until it has finished or timed out.
func Monitor(currentTarget types.Target, RID string) (types.Analysis, error) {

	analysis := types.Analysis{}
	timeout := time.After(15 * time.Minute)
	retryTick := time.Tick(60 * time.Second)

	for {
		select {
		case <-timeout:
			return analysis, errors.New("timeout")
		case <-retryTick:
			analysis, err := Get(currentTarget, RID)
			if err != nil {
				return analysis, err
			}
			if analysis.Status == "finished" {
				return analysis, nil
			} else if analysis.Status == "error running" {
				return analysis, fmt.Errorf("huskyCI encountered an error trying to execute this analysis: %v", analysis.ErrorFound)
			}
			if !types.IsJSONoutput {
				fmt.Println("[HUSKYCI][!] Hold on! huskyCI is still running...")
			}
		}
	}
}

// Get gets the results of an analysis.
func Get(currentTarget types.Target, RID string) (types.Analysis, error) {

	netClient := createClient()
	getAnalysisURL := currentTarget.Endpoint + "/analysis/" + RID
	req, err := http.NewRequest("GET", getAnalysisURL, nil)
	req.Header.Add("Husky-Token", currentTarget.Token)
	if err != nil {
		return types.Analysis{}, err
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return types.Analysis{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Analysis{}, err
	}

	analysis := types.Analysis{}
	err = json.Unmarshal(body, &analysis)
	if err != nil {
		return analysis, err
	}

	return analysis, nil
}

// PrintResults prints huskyCI output either in JSON or the standard output.
func PrintResults(analysis types.Analysis, analysisRunnerResults types.JSONOutput) error {

	prepareAllSummary(analysis, analysisRunnerResults)

	if types.IsJSONoutput {
		err := printJSONOutput(analysisRunnerResults)
		if err != nil {
			return err
		}
	} else {
		printSTDOUTOutput(analysis, analysisRunnerResults)
	}

	return nil
}

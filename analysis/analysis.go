package analysis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/globocom/husky-client/config"
	"github.com/globocom/husky-client/types"
)

// StartAnalysis starts a container and returns its error
func StartAnalysis() (string, error) {

	requestPayload := types.JSONPayload{
		RepositoryURL: config.RepositoryURL,
	}

	marshalPayload, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Could not Marshal requestPayload:", err)
		return "", err
	}

	huskyStartAnalysisURL := config.HuskyAPI + "husky"
	req, err := http.NewRequest("POST", huskyStartAnalysisURL, bytes.NewBuffer(marshalPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error during POST to Husky API:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body response of POST to Husky API:", err)
		return "", err
	}

	responsePayload := types.JSONResponse{}
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		fmt.Println("Could not Unmarshal responsePayload:", err)
		return "", err
	}

	return responsePayload.RID, nil
}

// GetAnalyisis gets
func GetAnalyisis(RID string) (types.Analysis, error) {

	analysis := types.Analysis{}
	huskyMonitorAnalysisURL := config.HuskyAPI + "husky/" + RID

	resp, err := http.Get(huskyMonitorAnalysisURL)
	if err != nil {
		fmt.Println("Error during GET to Husky API:", err)
		return analysis, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading GET response to Husky API:", err)
		return analysis, err
	}

	err = json.Unmarshal(body, &analysis)
	if err != nil {
		fmt.Println("Could not Unmarshal analysis:", err)
		return analysis, err
	}

	return analysis, nil
}

// MonitorAnalysis gets
func MonitorAnalysis(RID string) (types.Analysis, error) {

	analysis := types.Analysis{}
	timeout := time.After(10 * time.Minute)
	retryTick := time.Tick(10 * time.Second)

	for {
		select {
		case <-timeout:
			return analysis, errors.New("time out")
		case <-retryTick:
			analysis, err := GetAnalyisis(RID)
			if err != nil {
				io.WriteString(os.Stderr, "Error!")
				return analysis, err
			}
			if analysis.Status == "finished" {
				return analysis, nil
			}
		}
	}
}

// AnalyzeResult analyzes.
func AnalyzeResult(analysisResult types.Analysis) {
	// result = passed? sucess! Close client. result = failed? Throw error. Output cOutput where cResult = failed.
	if analysisResult.Result != "passed" {
		// print cOutput of each container that has cResult == failed and throw an error
		for _, container := range analysisResult.Containers {
			if container.CResult == "failed" {
				// cOutput is a string that needs to become a JSON
				fmt.Println(container.COuput)
			}
		}
	} else {
		// print Sucess! Warnings!
		fmt.Println(`{"Husky":"Success"}`)
	}
}

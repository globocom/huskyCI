// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/globocom/huskyCI/client/integration/sonarqube"

	"github.com/globocom/huskyCI/client/analysis"
	"github.com/globocom/huskyCI/client/config"
	"github.com/globocom/huskyCI/client/types"
)

func main() {

	types.FoundVuln = false
	types.IsJSONoutput = false

	if len(os.Args) > 1 && os.Args[1] == "JSON" {
		types.IsJSONoutput = true
	}

	// step 0: check and set huskyci-client configuration
	if err := config.CheckEnvVars(); err != nil {
		if !types.IsJSONoutput {
			fmt.Println("[HUSKYCI][ERROR] Check environment variables:", err)
		}
		os.Exit(1)
	}
	config.SetConfigs()

	// step 1: start analysis and get its RID.
	if !types.IsJSONoutput {
		s := fmt.Sprintf("[HUSKYCI][*] %s -> %s", config.RepositoryBranch, config.RepositoryURL)
		fmt.Println(s)
	}
	RID, err := analysis.StartAnalysis()
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Sending request to huskyCI:", err)
		os.Exit(1)
	}
	if !types.IsJSONoutput {
		fmt.Println("[HUSKYCI][*] huskyCI analysis started!", RID)
	}

	// step 2.1: keep querying huskyCI API to check if a given analysis has already finished.
	huskyAnalysis, err := analysis.MonitorAnalysis(RID)
	if err != nil {
		s := fmt.Sprintf("[HUSKYCI][ERROR] Monitoring analysis %s: %s", RID, err)
		fmt.Println(s)
		os.Exit(1)
	}

	// step 2.2: prepare the list of securityTests that ran in the analysis.
	var passedList []string
	var failedList []string
	var errorList []string
	for _, container := range huskyAnalysis.Containers {
		securityTestFullName := fmt.Sprintf("%s:%s", container.SecurityTest.Image, container.SecurityTest.ImageTag)
		if container.CResult == "passed" && container.SecurityTest.Name != "gitauthors" {
			passedList = append(passedList, securityTestFullName)
		} else if container.CResult == "failed" {
			failedList = append(failedList, securityTestFullName)
		} else if container.CResult == "error" {
			failedList = append(errorList, securityTestFullName)
		}
	}

	// step 3: print output based on os.Args(1) parameter received
	types.IsJSONoutput = false
	if len(os.Args) > 1 {
		types.IsJSONoutput = true
	}

	err = analysis.PrintResults(huskyAnalysis)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Printing output:", err)
		os.Exit(1)
	}

	// step 3.5: integration with SonarQube
	outputPath := "./huskyCI/"
	outputFileName := "sonarqube.json"
	err = sonarqube.GenerateOutputFile(huskyAnalysis, outputPath, outputFileName)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Could not create SonarQube integration file: ", err)
	}

	// step 4: block developer CI if vulnerabilities were found
	if !types.FoundVuln && !types.FoundInfo {
		if !types.IsJSONoutput {
			if len(errorList) > 0 {
				fmt.Println("[HUSKYCI][*] The following securityTests failed to run:")
				fmt.Println("[HUSKYCI][*]", errorList)
			}
			fmt.Println("[HUSKYCI][*] The following securityTests were executed and no blocking vulnerabilities were found:")
			fmt.Println("[HUSKYCI][*]", passedList)
			fmt.Println("[HUSKYCI][*] No issues were found.")
		}
		os.Exit(0)
	}

	if !types.FoundVuln && types.FoundInfo {
		if !types.IsJSONoutput {
			if len(errorList) > 0 {
				fmt.Println("[HUSKYCI][*] The following securityTests failed to run:")
				fmt.Println("[HUSKYCI][*]", errorList)
			}
			fmt.Println("[HUSKYCI][*] The following securityTests were executed and no blocking vulnerabilities were found:")
			fmt.Println("[HUSKYCI][*]", passedList)
			fmt.Println("[HUSKYCI][*] However, some LOW/INFO issues were found...")
		}
		os.Exit(0)
	}

	if types.FoundVuln && !types.IsJSONoutput {
		if len(errorList) > 0 {
			fmt.Println("[HUSKYCI][*] The following securityTests failed to run:")
			fmt.Println("[HUSKYCI][*]", errorList)
		}
		if len(passedList) > 0 {
			fmt.Println("[HUSKYCI][*] The following securityTests were executed and no blocking vulnerabilities were found:")
			fmt.Println("[HUSKYCI][*]", passedList)
		}
		fmt.Println("[HUSKYCI][*] Some HIGH/MEDIUM issues were found in these securityTests:")
		fmt.Println("[HUSKYCI][*]", failedList)
	}

	os.Exit(190)
}

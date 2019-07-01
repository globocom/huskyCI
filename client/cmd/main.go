// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/globocom/huskyCI/client/analysis"
	"github.com/globocom/huskyCI/client/config"
	"github.com/globocom/huskyCI/client/types"
)

func main() {

	types.FoundVuln = false

	// step 0: check and set huskyci-client configuration
	if err := config.CheckEnvVars(); err != nil {
		fmt.Println("[HUSKYCI][ERROR] Check environment variables:", err)
		os.Exit(1)
	}
	config.SetConfigs()
	fmt.Println(fmt.Sprintf("[HUSKYCI][*] %s -> %s", config.RepositoryBranch, config.RepositoryURL))

	// step 1: start analysis and get a RID.
	RID, err := analysis.StartAnalysis()
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Sending request to HuskyCI:", err)
		os.Exit(1)
	}

	fmt.Println("[HUSKYCI][*] huskyCI analysis started!", RID)

	// step 2: keep querying husky API to check if a given analysis has already finished.
	huskyAnalysis, err := analysis.MonitorAnalysis(RID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[HUSKYCI][ERROR] Monitoring analysis %s: %s", RID, err))
		os.Exit(1)
	}

	// step 3: analyze result
	analysis.AnalyzeResult(huskyAnalysis)

	// step 4: print output
	formatOutput := ""
	if len(os.Args) > 1 {
		formatOutput = "JSON"
	}
	err = analysis.PrintResults(formatOutput)
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Printing output:", err)
		os.Exit(1)
	}

	// step 5: block developer CI if vulnerabilities were found
	if types.FoundVuln == true {
		os.Exit(1)
	}

	if types.FoundInfo == true {
		fmt.Printf("[HUSKYCI][*] Some low/info issues were found :|\n")
		fmt.Println()
		os.Exit(0)
	}

	fmt.Printf("[HUSKYCI][*] Nice! No issues Found :)\n")
	fmt.Println()
	os.Exit(0)
}

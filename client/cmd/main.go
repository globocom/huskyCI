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

	// step 1: start analysis and get its RID.
	RID, err := analysis.StartAnalysis()
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Sending request to huskyCI:", err)
		os.Exit(1)
	}

	fmt.Println("[HUSKYCI][*] huskyCI analysis started!", RID)

	// step 2: keep querying huskyCI API to check if a given analysis has already finished.
	huskyAnalysis, err := analysis.MonitorAnalysis(RID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[HUSKYCI][ERROR] Monitoring analysis %s: %s", RID, err))
		os.Exit(1)
	}

	// step 3: prepares huskyCI results into structs
	analysis.PrepareResults(huskyAnalysis)

	// step 4: print output based on os.Args(1) parameter received
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
		if len(os.Args) < 1 {
			fmt.Printf("[HUSKYCI][*] Some high issues were found :(\n")
			os.Exit(1)
		}
		os.Exit(1)
	}

	if types.FoundInfo == true && len(os.Args) < 1 {
		fmt.Printf("[HUSKYCI][*] Some low/info issues were found :|\n")
		os.Exit(0)
	}

	if len(os.Args) < 1 {
		fmt.Printf("[HUSKYCI][*] Nice! No issues were found :)\n")
	}
	os.Exit(0)
}

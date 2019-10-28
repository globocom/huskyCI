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
		fmt.Println(fmt.Sprintf("[HUSKYCI][*] %s -> %s", config.RepositoryBranch, config.RepositoryURL))
	}
	RID, err := analysis.StartAnalysis()
	if err != nil {
		fmt.Println("[HUSKYCI][ERROR] Sending request to huskyCI:", err)
		os.Exit(1)
	}
	if !types.IsJSONoutput {
		fmt.Println("[HUSKYCI][*] huskyCI analysis started!", RID)
	}

	// step 2: keep querying huskyCI API to check if a given analysis has already finished.
	huskyAnalysis, err := analysis.MonitorAnalysis(RID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[HUSKYCI][ERROR] Monitoring analysis %s: %s", RID, err))
		os.Exit(1)
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

	// step 4: block developer CI if vulnerabilities were found
	if !types.FoundVuln && !types.FoundInfo {
		if !types.IsJSONoutput {
			fmt.Printf("[HUSKYCI][*] Nice! No issues were found :)\n")
		}
		os.Exit(0)
	}

	if !types.FoundVuln && types.FoundInfo {
		if !types.IsJSONoutput {
			fmt.Printf("[HUSKYCI][*] Some LOW/INFO issues were found :|\n")
		}
		os.Exit(0)
	}

	if types.FoundVuln && !types.IsJSONoutput {
		fmt.Printf("[HUSKYCI][*] Some HIGH/MEDIUM issues were found :(\n")
	}

	os.Exit(1)
}

package main

import (
	"fmt"
	"os"

	"github.com/globocom/husky-client/analysis"
	"github.com/globocom/husky-client/config"
)

func main() {

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

	fmt.Println("[HUSKYCI][*] HuskyCI analysis started!", RID)

	// step 2: keep querying husky API to check if a given analysis has already finished.
	huskyAnalysis, err := analysis.MonitorAnalysis(RID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[HUSKYCI][ERROR] Monitoring analysis %s: %s", RID, err))
		os.Exit(1)
	}

	// step 3: analyze result and return to CI the final result.
	analysis.AnalyzeResult(huskyAnalysis)

}

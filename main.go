package main

import (
	"fmt"
	"io"
	"os"

	"github.com/globocom/husky-client/analysis"
	"github.com/globocom/husky-client/config"
)

func main() {

	// step 0: set all configs needed.
	config.SetConfigs()

	// step 1: start analysis using received parameters and get a RID.
	RID, err := analysis.StartAnalysis()
	if err != nil {
		fmt.Println("Error StartAnalysis():", err)
	}

	// step 2: keep querying husky API to check if a given analysis has already finished.
	huskyAnalysis, err := analysis.MonitorAnalysis(RID)
	if err != nil {
		fmt.Println("Error MonitorAnalysis():", err)
	}

	// step 3: analyze result and return to CI the final result.
	finalResultError := analysis.AnalyzeResult(huskyAnalysis)
	if finalResultError != nil {
		io.WriteString(os.Stderr, "[HUSKY][FAILED]")
	} else {
		io.WriteString(os.Stdout, "[HUSKY][PASSED]")
	}
}

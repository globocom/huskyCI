// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"time"

	"github.com/globocom/huskyCI/client/integration/sonarqube"

	"github.com/globocom/huskyCI/cli/client"
	"github.com/globocom/huskyCI/cli/config"
	"github.com/globocom/huskyCI/client/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an huskyCI analysis",
	Long:  `Run a security analysis using huskyCI backend.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get current target
		currentTarget, err := config.GetCurrentTarget()
		if err != nil {
			fmt.Printf("[HUSKYCI][ERROR] Current target not detected: (%v)", err)
		}
		fmt.Printf("[HUSKYCI][*] Using target %s -> %s\n", currentTarget.Label, currentTarget.Endpoint)

		fmt.Printf("[HUSKYCI][*] %s -> %s\n", viper.Get("repo_branch"), viper.Get("repo_url"))

		hcli := client.NewClient(*currentTarget)
		analysisRunnerResults, err := hcli.Start(viper.GetString("repo_url"), viper.GetString("repo_branch"))
		if err != nil {
			return fmt.Errorf("[HUSKYCI][ERROR] Sending request to huskyCI: %s", err.Error())
		}

		fmt.Printf("[HUSKYCI][*] huskyCI analysis started: %s\n", analysisRunnerResults.RID)

		timeoutMonitor, err := time.ParseDuration("10m")
		if err != nil {
			return fmt.Errorf("[HUSKYCI][ERROR] Internal error %v", err)
		}
		retryMonitor, err := time.ParseDuration("30s")
		if err != nil {
			return fmt.Errorf("[HUSKYCI][ERROR] Internal error %v", err)
		}

		analysisResult, err := hcli.Monitor(analysisRunnerResults.RID, timeoutMonitor, retryMonitor)
		if err != nil {
			return fmt.Errorf("[HUSKYCI][ERROR] Monitoring analysis %s: %v", analysisRunnerResults.RID, err)
		}

		var outputJSON types.JSONOutput
		hcli.PrintResults(analysisResult, outputJSON)

		if !viper.GetBool("found_vuln") {
			if !viper.GetBool("found_info") {
				fmt.Printf("[HUSKYCI][*] Nice! No issues were found :)\n")
			} else {
				fmt.Printf("[HUSKYCI][*] Some LOW/INFO issues were found :|\n")
			}
		} else {
			fmt.Printf("[HUSKYCI][*] Some HIGH/MEDIUM issues were found :(\n")
		}

		outputPath := "./huskyCI/"
		outputFileName := "sonarqube.json"
		err = sonarqube.GenerateOutputFile(analysisResult, outputPath, outputFileName)
		if err != nil {
			fmt.Println("[HUSKYCI][ERROR] Could not create SonarQube integration file: ", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

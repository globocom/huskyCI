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

package cmd

import (
	"fmt"
	"time"

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

		timeoutMonitor, _ := time.ParseDuration("10m")
		retryMonitor, _ := time.ParseDuration("30s")

		analysisResult, err := hcli.Monitor(analysisRunnerResults.RID, timeoutMonitor, retryMonitor)
		if err != nil {
			return fmt.Errorf("[HUSKYCI][ERROR] Monitoring analysis %s: %v", analysisRunnerResults.RID, err)
		}

		var outputJSON types.JSONOutput
		err = hcli.PrintResults(analysisResult, outputJSON)
		if err != nil {
			return fmt.Errorf("[HUSKYCI][ERROR] Printing output: (%v)", err)
		}

		if viper.GetBool("found_vuln") == false {
			if viper.GetBool("found_info") == false {
				fmt.Printf("[HUSKYCI][*] Nice! No issues were found :)\n")
			}
			if viper.GetBool("found_info") == true {
				fmt.Printf("[HUSKYCI][*] Some LOW/INFO issues were found :|\n")
			}
		} else {
			fmt.Printf("[HUSKYCI][*] Some HIGH/MEDIUM issues were found :(\n")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

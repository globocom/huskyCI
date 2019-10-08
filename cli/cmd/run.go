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
	"os"

	"github.com/globocom/huskyCI/cli/analysis"
	"github.com/globocom/huskyCI/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an huskyCI analysis",
	Long:  `Run a security analysis using huskyCI backend.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get current target
		currentTarget, err := config.GetCurrentTarget()
		if err != nil {
			fmt.Printf("[HUSKYCI][ERROR] Current target not detected: (%v)", err)
		}
		fmt.Printf("[HUSKYCI][*] Using target %s -> %s\n", currentTarget.Label, currentTarget.Endpoint)

		fmt.Printf("[HUSKYCI][*] %s -> %s\n", viper.Get("repo_branch"), viper.Get("repo_url"))

		analysisRunnerResults, err := analysis.Start(*currentTarget)
		if err != nil {
			fmt.Println("[HUSKYCI][ERROR] Sending request to huskyCI:", err)
			os.Exit(1)
		}

		fmt.Printf("[HUSKYCI][*] huskyCI analysis started: %s\n", analysisRunnerResults.Summary.RID)

		analysisResult, err := analysis.Monitor(*currentTarget, analysisRunnerResults.Summary.RID)
		if err != nil {
			fmt.Printf("[HUSKYCI][ERROR] Monitoring analysis %s: %v\n", analysisRunnerResults.Summary.RID, err)
			os.Exit(1)
		}

		err = analysis.PrintResults(analysisResult, analysisRunnerResults)
		if err != nil {
			fmt.Printf("[HUSKYCI][ERROR] Printing output: (%v)\n", err)
			os.Exit(1)
		}

		if viper.GetBool("found_vuln") == false {
			if viper.GetBool("found_info") == false {
				fmt.Printf("[HUSKYCI][*] Nice! No issues were found :)\n")
				os.Exit(0)
			}
			if viper.GetBool("found_info") == true {
				fmt.Printf("[HUSKYCI][*] Some LOW/INFO issues were found :|\n")
				os.Exit(0)
			}
		} else {
			fmt.Printf("[HUSKYCI][*] Some HIGH/MEDIUM issues were found :(\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// targetSetCmd represents the targetSet command
var targetSetCmd = &cobra.Command{
	Use:   "target-set",
	Short: "Change current target (huskyci api)",
	Long: `
Change current target (huskyci api).
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// check if target name is used
		notUsed := true
		endpoint := ""
		targets := viper.GetStringMap("targets")
		for k, v := range targets {
			if k == args[0] {
				notUsed = false
				// set target as current
				target := v.(map[string]interface{})
				target["current"] = true
				endpoint = target["endpoint"].(string)
			} else {
				// unset all others targets as not current
				target := v.(map[string]interface{})
				target["current"] = false
			}
		}
		if notUsed {
			fmt.Printf("Client error, target does not exist: %s\n", args[0])
			os.Exit(1)
		}

		// save config (only if target is found)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("Client error saving config with current target: (%s)\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("New target is %s -> %s\n", args[0], endpoint)
	},
}

func init() {
	rootCmd.AddCommand(targetSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetSetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

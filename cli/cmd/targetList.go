// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// targetListCmd represents the targetList command
var targetListCmd = &cobra.Command{
	Use:   "target-list",
	Short: "Displays the list of targets, marking the current",
	Long: `
Displays the list of targets, marking the current.
`,
	Run: func(cmd *cobra.Command, args []string) {

		targets := viper.GetStringMap("targets")
		for k, v := range targets {
			target := v.(map[string]interface{})

			// format output for activated target
			if target["current"].(bool) {
				target["currented"] = "*"
			} else {
				target["currented"] = " "
			}

			fmt.Printf("%s %s (%s)\n", target["currented"], k, target["endpoint"])
		}
	},
}

func init() {
	rootCmd.AddCommand(targetListCmd)
}

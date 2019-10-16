// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {

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
			return fmt.Errorf("Client error, target does not exist: %s", args[0])
		}

		// save config (only if target is found)
		err := viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("Client error saving config with current target: (%s)", err.Error())
		}
		fmt.Printf("New target is %s -> %s\n", args[0], endpoint)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(targetSetCmd)
}

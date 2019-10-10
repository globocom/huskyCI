// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// targetRemoveCmd represents the targetRemove command
var targetRemoveCmd = &cobra.Command{
	Use:   "target-remove",
	Short: "Remove a target from target-list (huskyci api)",
	Long: `
	Remove a target from target-list (huskyci api).
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// check if target name is used
		notUsed := true
		targets := viper.GetStringMap("targets")
		for k := range targets {
			if k == args[0] {
				notUsed = false
			}
		}
		if notUsed {
			return fmt.Errorf("Error, target does not exist: %s", args[0])
		}

		// remove entry from data struct but, before, storing data to show to user
		target := targets[args[0]].(map[string]interface{})
		endpoint := target["endpoint"].(string)
		targets[args[0]] = nil

		// save config
		err := viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("Client error saving config without target: (%s)", err.Error())
		}

		fmt.Printf("Target %s -> %s removed from target list\n", args[0], endpoint)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(targetRemoveCmd)
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// targetAddCmd represents the targetAdd command
var targetAddCmd = &cobra.Command{
	Use:   "target-add [name] [https://huskyci-api-endpoint.example.com]",
	Short: "Adds a new entry to the list of available targets",
	Long: `
Adds a new entry to the list of available targets.
	`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		match, err := regexp.MatchString(`^\w+$`, args[0])
		if err != nil {
			return fmt.Errorf("Client error validanting label name: %s (%s)", args[0], err.Error())
		}
		if !match {
			return fmt.Errorf("Client error parsing target name: %s (must have number, letters and/or underscores)", args[0])
		}

		// check huskyci-api-endpoint
		_, err = url.Parse(args[1])
		if err != nil {
			return fmt.Errorf("Client error parsing target endpoint: %s (%s)", args[1], err.Error())
		}

		// check if target name is used before
		targets := viper.GetStringMap("targets")
		for k, v := range targets {
			if k == args[0] {
				target := v.(map[string]interface{})
				return fmt.Errorf("Client error, target name already exists: %s (with endpoint: %s)", k, target["endpoint"])
			}
		}

		// if new target must be current, we unset all others
		setCurrent, err := cmd.Flags().GetBool("set-current")
		if err != nil {
			return fmt.Errorf("Client error parsing set-current option: (%s)", err.Error())
		}

		if setCurrent {
			for _, v := range targets {
				target := v.(map[string]interface{})
				target["current"] = false
			}
		}

		// add new entry to data struct
		targets[args[0]] = map[string]interface{}{"current": setCurrent, "endpoint": args[1]}

		// save config
		viper.Set("targets", targets)
		err = viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("Client error saving config with new target: (%s)", err.Error())
		}
		fmt.Printf("New target %s -> %s added to target list\n", args[0], args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(targetAddCmd)
	targetAddCmd.Flags().BoolP("set-current", "s", false, "Add and define the target as the current target")
}

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

		match, _ := regexp.MatchString(`^\w+$`, args[0])
		if !match {
			return fmt.Errorf("Client error parsing target name: %s (must have number, letters and/or underscores)", args[0])
		}

		// check huskyci-api-endpoint
		_, err := url.Parse(args[1])
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
			fmt.Errorf("Client error parsing set-current option: (%s)", err.Error())
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	targetAddCmd.Flags().BoolP("set-current", "s", false, "Add and define the target as the current target")
}

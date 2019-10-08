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

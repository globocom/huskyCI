// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"

	"github.com/globocom/huskyCI/cli/analysis"
	"github.com/globocom/huskyCI/cli/errorcli"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a huskyCI analysis",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			errorcli.Handle(errors.New("path is missing"))
		}
		return nil
	},
	// Long:  `Run a security analysis using huskyCI backend.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		pathReceived := args[0]
		currentAnalysis := analysis.New()

		fmt.Println()
		if err := currentAnalysis.CheckPath(pathReceived); err != nil {
			errorcli.Handle(err)
		}

		fmt.Println()
		if err := currentAnalysis.CompressFiles(pathReceived); err != nil {
			errorcli.Handle(err)
		}

		fmt.Println()
		if err := currentAnalysis.SendZip(); err != nil {
			errorcli.Handle(err)
		}

		fmt.Println()
		if err := currentAnalysis.CheckStatus(); err != nil {
			errorcli.Handle(err)
		}

		fmt.Println()
		currentAnalysis.PrintVulns()

		if err := currentAnalysis.HouseCleaning(); err != nil {
			errorcli.Handle(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

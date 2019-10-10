/// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/globocom/huskyCI/cli/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "huskyci",
	Short: "huskyci is a CLI to use huskyCI",
	Long: `huskyci is a CLI to use huskyCI
	
huskyCI is an Open Source tool that performs security tests
inside CI pipelines of multiple projects and centralizes all results 
into a database for further analysis and metrics.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.huskyci/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Printf("Client error reading home folder: %s (%s)\n", home, err.Error())
			os.Exit(1)
		}

		// check if .huskyci folder exists and creates if it not exists
		path, err := config.CheckAndCreateConfigFolder(home, true)
		if err != nil {
			os.Exit(1)
		}

		// add path to viper config
		viper.AddConfigPath(path)

		// set config file name and type
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// test if config file exists and creates if it not exists
		err = viper.ReadInConfig()
		if err != nil {
			_, err = config.CreateConfigFile(path, true)
			if err != nil {
				os.Exit(1)
			}
		}
	}

	viper.SetEnvPrefix("HUSKYCI_CLIENT")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Client error reading config file (%s)\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
}

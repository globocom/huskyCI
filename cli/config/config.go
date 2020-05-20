// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"

	"github.com/globocom/huskyCI/cli/types"
	"github.com/spf13/viper"
)

// GetCurrentTarget return a types.Target with current target
func GetCurrentTarget() (*types.Target, error) {
	// Get current targets
	currentTarget := new(types.Target)

	if len(os.Getenv("HUSKYCI_CLIENT_API_ADDR")) > 0 {
		currentTarget.Endpoint = os.Getenv("HUSKYCI_CLIENT_API_ADDR")
		currentTarget.Label = "env-var"
		currentTarget.TokenStorage = "env-var"
		currentTarget.Token = os.Getenv("HUSKYCI_CLIENT_TOKEN")
	} else {
		targets := viper.GetStringMap("targets")
		for k, v := range targets {
			target := v.(map[string]interface{})

			// check if target is properly configured
			if target["current"] == nil {
				return nil, fmt.Errorf("You need to configure a target using target-add command")
			}

			// format output for activated target
			if target["current"] != nil && target["current"].(bool) {
				currentTarget.Label = k
				currentTarget.Endpoint = target["endpoint"].(string)

				// check token storage
				if target["token-storage"] == nil {
					fmt.Printf("Token storage is not set. You can set it using the -s flag at 'huskyci login' command\n")
					currentTarget.TokenStorage = ""
				} else {
					currentTarget.TokenStorage = target["token-storage"].(string)
				}

			}
		}

	}

	return currentTarget, nil
}

// CheckAndCreateConfigFolder check if config folder exists and create if it doesn't exists
func CheckAndCreateConfigFolder(home string, debug bool) (string, error) {
	// check if .huskyci folder exists and creates if it not exists
	path := home + "/.huskyci"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0750)
		if err != nil {
			if debug {
				fmt.Printf("Client error creating config folder: %s (%s)\n", path, err.Error())
			}
			return "", err
		}
	}
	if debug {
		fmt.Printf("Client created config folder: %s\n", path)
	}
	return path, nil
}

// CreateConfigFile creates a config file for huskyci CLI
func CreateConfigFile(path string, debug bool) (string, error) {
	configFile := path + "/config.yaml"
	file, err := os.Create(configFile)
	if err != nil {
		if debug {
			fmt.Printf("Client error creating config file: %s (%s)\n", configFile, err.Error())
		}
		return "", err
	}
	err = file.Close()
	if err != nil {
		if debug {
			fmt.Printf("Client error closing config file: %s (%s)\n", configFile, err.Error())
		}
		return "", err
	}
	if debug {
		fmt.Printf("Client created new config file: %s\n", configFile)
	}
	return configFile, nil
}

// GetHuskyZipFilePath returns "$HOME/.huskyci/compressed-code.zip" and an error.
// If .huskyci folder is not present, the CLI will create it.
func GetHuskyZipFilePath() (string, error) {

	var fullFilePath string

	home, err := os.UserHomeDir()
	if err != nil {
		return fullFilePath, err
	}

	huskyHome, err := CheckAndCreateConfigFolder(home, false)
	if err != nil {
		return fullFilePath, err
	}

	fullFilePath = fmt.Sprintf("%s/%s", huskyHome, "compressed-code.zip")

	return fullFilePath, nil
}

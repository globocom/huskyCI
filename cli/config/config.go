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

package config

import (
	"fmt"
	"os"

	"github.com/globocom/huskyCI/client/types"
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

			if target["current"] != nil {
				// format output for activated target
				if target["current"].(bool) {
					currentTarget.Label = k
					currentTarget.Endpoint = target["endpoint"].(string)

					// check token storage
					if target["token-storage"] == nil {
						fmt.Printf("Token storage is not set. You can set it using the -s flag at 'husky login' command\n")
						currentTarget.TokenStorage = ""
					} else {
						currentTarget.TokenStorage = target["token-storage"].(string)
					}
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

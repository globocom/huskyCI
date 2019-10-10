// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGetCurrentTarget(t *testing.T) {
	t.Run(
		"Test GetCurrentTarget() from env-var",
		func(t *testing.T) {
			// Set control env var
			os.Setenv("HUSKYCI_CLIENT_API_ADDR", "https://env.example.com:443")

			currentTarget, err := GetCurrentTarget()
			if err != nil {
				t.Fatalf("CONFIG: fail to read API ADDR from env var (%v)", err)
			}
			if currentTarget.Endpoint != "https://env.example.com:443" {
				t.Fatalf("CONFIG: fail to read API ADDR from env var (%v)", currentTarget)
			}
		})
	t.Run(
		"Test GetCurrentTarget() from viper (without targets)",
		func(t *testing.T) {
			// Unset control env var
			os.Unsetenv("HUSKYCI_CLIENT_API_ADDR")

			// Generate targets
			targets := viper.GetStringMap("targets")
			targets["test"] = map[string]interface{}{"current": nil, "endpoint": "https://without.example.com:443", "token-storage": "keychain"}
			viper.Set("targets", targets)

			_, err := GetCurrentTarget()
			if err == nil {
				t.Fatalf("CONFIG: fail to identify targets from viper without targets (%v)", err)
			}
		})
	t.Run(
		"Test GetCurrentTarget() from viper (without token-storage)",
		func(t *testing.T) {
			// Unset control env var
			os.Unsetenv("HUSKYCI_CLIENT_API_ADDR")

			// Generate targets
			targets := viper.GetStringMap("targets")
			targets["test"] = map[string]interface{}{"current": true, "endpoint": "https://viper.example.com:443"}
			viper.Set("targets", targets)

			currentTarget, err := GetCurrentTarget()
			if err != nil {
				t.Fatalf("CONFIG: fail to read API ADDR from viper (%v)", err)
			}
			if currentTarget.Endpoint != "https://viper.example.com:443" {
				t.Fatalf("CONFIG: fail to read API ADDR from viper (%v)", currentTarget)
			}
		})
	t.Run(
		"Test GetCurrentTarget() from viper (with token-storage)",
		func(t *testing.T) {
			// Unset control env var
			os.Unsetenv("HUSKYCI_CLIENT_API_ADDR")

			// Generate targets
			targets := viper.GetStringMap("targets")
			targets["test"] = map[string]interface{}{"current": true, "endpoint": "https://viper.example.com:443", "token-storage": "keychain"}
			viper.Set("targets", targets)

			currentTarget, err := GetCurrentTarget()
			if err != nil {
				t.Fatalf("CONFIG: fail to read API ADDR from viper (%v)", err)
			}
			if currentTarget.Endpoint != "https://viper.example.com:443" {
				t.Fatalf("CONFIG: fail to read API ADDR from viper (%v)", currentTarget)
			}
		})
}

func TestCheckAndCreateConfigFolder(t *testing.T) {
	t.Run(
		"Test CheckAndCreateConfigFolder()",
		func(t *testing.T) {
			// Create temp dir
			dir, err := ioutil.TempDir("", "TestCheckAndCreateConfigFolder")
			if err != nil {
				t.Fatalf("CONFIG: (pre-test) fail to create config folder (%v)", err)
			}

			_, err = CheckAndCreateConfigFolder(dir, false)
			if err != nil {
				t.Fatalf("CONFIG: fail to create config folder (%v)", err)
			}

			// Clean environment
			defer os.RemoveAll(dir)
		})
	t.Run(
		"Test CheckAndCreateConfigFolder() without permissions",
		func(t *testing.T) {
			// Create temp dir
			dir, err := ioutil.TempDir("", "TestCheckAndCreateConfigFolder")
			if err != nil {
				t.Fatalf("CONFIG: (pre-test) fail to create config folder (%v)", err)
			}

			// Change permissions
			if err := os.Chmod(dir, 0000); err != nil {
				t.Fatalf("Internal Error: (%v)", err)
			}

			_, err = CheckAndCreateConfigFolder(dir, true)
			if err != nil {
				t.Fatalf("CONFIG: fail to create config folder (%v)", err)
			}

			// Change permissions
			if err := os.Chmod(dir, 0777); err != nil {
				t.Fatalf("Internal Error: (%v)", err)
			}

			// Clean environment
			defer os.RemoveAll(dir)
		})
}

func TestCreateConfigFile(t *testing.T) {
	t.Run(
		"Test CreateConfigFile()",
		func(t *testing.T) {
			// Create temp dir
			dir, err := ioutil.TempDir("", "TestCreateConfigFile")
			if err != nil {
				t.Fatalf("CONFIG: (pre-test) fail to create config folder (%v)", err)
			}
			configFolder, err := CheckAndCreateConfigFolder(dir, true)
			if err != nil {
				t.Fatalf("Internal Error: (%v)", err)
			}

			_, err = CreateConfigFile(configFolder, true)
			if err != nil {
				t.Fatalf("CONFIG: fail to create config file (%v)", err)
			}

			// Clean environment
			defer os.RemoveAll(dir)
		})

}

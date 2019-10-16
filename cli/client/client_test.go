// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"os"
	"testing"
	"time"

	"github.com/globocom/huskyCI/client/types"
)

func TestCreateHTTPClient(t *testing.T) {
	t.Run(
		"Test CreateHTTPClient() without env-var",
		func(t *testing.T) {
			os.Unsetenv("HUSKYCI_CLIENT_TIMEOUT")
			client := createHTTPClient()

			defaultDuration, err := time.ParseDuration("10s")
			if err != nil {
				t.Fatalf("Internal Error")
			}

			if client.Timeout != defaultDuration {
				t.Fatalf("CLIENT: fail to read default timeout (%v)", client.Timeout)
			}
		})
	t.Run(
		"Test CreateHTTPClient() with env-var",
		func(t *testing.T) {
			os.Setenv("HUSKYCI_CLIENT_TIMEOUT", "15s")
			client := createHTTPClient()

			changedDuration, err := time.ParseDuration("15s")
			if err != nil {
				t.Fatalf("Internal Error")
			}

			if client.Timeout != changedDuration {
				t.Fatalf("CLIENT: fail to read non default timeout (%v)", client.Timeout)
			}
		})
}

func TestNewClient(t *testing.T) {
	t.Run(
		"Test NewClient() without env-var",
		func(t *testing.T) {
			os.Unsetenv("HUSKYCI_CLIENT_TIMEOUT")

			hclient := NewClient(types.Target{Endpoint: "https://example.com"})

			defaultDuration, err := time.ParseDuration("10s")
			if err != nil {
				t.Fatalf("Internal Error")
			}

			if hclient.httpCli.Timeout != defaultDuration {
				t.Fatalf("CLIENT: fail to read default timeout (%v)", hclient.httpCli.Timeout)
			}
		})
	t.Run(
		"Test NewClient() with env-var",
		func(t *testing.T) {
			os.Setenv("HUSKYCI_CLIENT_TIMEOUT", "15s")
			hclient := NewClient(types.Target{Endpoint: "https://example.com"})

			changedDuration, err := time.ParseDuration("15s")
			if err != nil {
				t.Fatalf("Internal Error")
			}

			if hclient.httpCli.Timeout != changedDuration {
				t.Fatalf("CLIENT: fail to read non default timeout (%v)", hclient.httpCli.Timeout)
			}
		})
}

// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"net"
	"net/http"
	"time"

	"github.com/globocom/huskyCI/client/types"
	"github.com/spf13/viper"
)

// Client has data to make API requests
type Client struct {
	target  types.Target
	httpCli *http.Client
}

// NewClient creates a custom Client
func NewClient(target types.Target) *Client {

	// Init config
	viper.SetEnvPrefix("HUSKYCI_CLIENT")
	viper.AutomaticEnv()

	timeout := viper.GetDuration("timeout")
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Setting custom HTTP client with timeouts
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}
	var netClient = &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}

	cli := Client{
		target:  target,
		httpCli: netClient,
	}

	return &cli
}

// creates a custom httpClient
func createHTTPClient() *http.Client {

	// Init config
	viper.SetEnvPrefix("HUSKYCI_CLIENT")
	viper.AutomaticEnv()

	timeout := viper.GetDuration("timeout")
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Setting custom HTTP client with timeouts
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}
	var netClient = &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}

	return netClient
}

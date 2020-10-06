// Copyright 2020 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"net/url"
)

const (
	// ClientID is a GitHub App client id.
	ClientID = "Iv1.d0c686433b9fc26f"

	// GrantTypeDeviceCode is a grant type for the device authorization flow.
	GrantTypeDeviceCode = "urn:ietf:params:oauth:grant-type:device_code"
)

// GitHub API paths.
const (
	LoginDeviceCodePath       = "login/device/code"
	LoginOAuthAccessTokenPath = "login/oauth/access_token" // #nosec - just a GitHub API path
)

// DefaultBaseURI is a default GitHub base URI.
var DefaultBaseURI, _ = url.Parse("https://github.com")

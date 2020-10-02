package github

import (
	"net/url"
)

const (
	// ClientID is a GitHub App client id.
	ClientID = "9c1379d55962c56e11a8"

	// GrantTypeDeviceCode is a grant type for the device authorization flow.
	GrantTypeDeviceCode = "urn:ietf:params:oauth:grant-type:device_code"
)

// DefaultBaseURI is a default GitHub base URI.
var DefaultBaseURI, _ = url.Parse("https://github.com")

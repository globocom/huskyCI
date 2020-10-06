// Copyright 2020 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// DeviceFlow allows to authorize users to GitHub.
// See https://docs.github.com/en/developers/apps/authorizing-oauth-apps#device-flow.
type DeviceFlow struct {
	baseURI *url.URL
	client  *http.Client
}

// NewDeviceFlow creates a new GitHub device authorization flow.
func NewDeviceFlow(baseURL *url.URL, client *http.Client) DeviceFlow {
	return DeviceFlow{baseURI: baseURL, client: client}
}

// GetCodes requests the device and user verification codes from GitHub.
func (df DeviceFlow) GetCodes(req *GetCodesRequest) (*GetCodesResponse, error) {
	uri := df.uri(LoginDeviceCodePath)
	resp := new(GetCodesResponse)
	if err := df.do(http.MethodPost, uri, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCodesRequest contains input parameters to request the device and user
// verification codes from GitHub.
type GetCodesRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope,omitempty"`
}

// GetCodesResponse contains response parameters when requesting the device
// and user verification codes from GitHub.
type GetCodesResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// GetAccessToken requests the access token for the authorized user.
func (df DeviceFlow) GetAccessToken(
	req *GetAccessTokenRequest,
) (*GetAccessTokenResponse, error) {
	uri := df.uri(LoginOAuthAccessTokenPath)
	resp := new(GetAccessTokenResponse)
	if err := df.do(http.MethodPost, uri, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAccessTokenRequest contains input parameters to request access token for
// the authorized user.
type GetAccessTokenRequest struct {
	ClientID   string `json:"client_id"`
	DeviceCode string `json:"device_code,omitempty"`
	GrantType  string `json:"grant_type,omitempty"`
}

// GetAccessTokenResponse contains response parameters when requesting access
// token for the authorized user.
type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (df DeviceFlow) uri(p string) string {
	result := *df.baseURI
	result.Path = path.Join(result.Path, p)
	return result.String()
}

func (df DeviceFlow) do(method, url string, in, out interface{}) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(in); err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := df.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := df.checkResponse(data); err != nil {
		return err
	}

	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}

func (DeviceFlow) checkResponse(data []byte) error {
	err := new(ErrResponse)
	decodeErr := json.Unmarshal(data, err)
	if decodeErr != nil {
		return decodeErr
	}
	if err.Error() != "" {
		return err
	}

	return nil
}

// ErrResponse is a GitHub error response.
type ErrResponse struct {
	Err string `json:"error"`
}

func (er ErrResponse) Error() string {
	return er.Err
}

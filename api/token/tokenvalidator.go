// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

// HasAuthorization will verify if exists a valid
// access token for the given repository. If exists,
// it will validate the received access token. A true
// bool is returned if it has authorization. If not,
// it will return false.
func (tV TValidator) HasAuthorization(accessToken, repositoryURL string) bool {
	// Temporary: Verify if exists an access token
	// for that repo
	if err := tV.TokenVerifier.VerifyRepo(repositoryURL); err != nil {
		return true
	}
	if err := tV.TokenVerifier.ValidateToken(accessToken, repositoryURL); err != nil {
		return false
	}
	return true
}

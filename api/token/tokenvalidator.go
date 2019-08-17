// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

func (tV TokenValidator) HasAuthorization(accessToken, repositoryURL string) bool {
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

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

package token

import (
	"errors"
	"github.com/globocom/huskyCI/api/types"
)

func (tH TokenHandler) GenerateAccessToken(repo types.TokenRequest) (types.AccessToken, error) {
	accessToken := types.AccessToken{}
	validatedURL, err := tH.External.ValidateURL(repo.RepositoryURL)
	if err != nil {
		return accessToken, err
	}
	token, err := tH.External.GenerateToken()
	if err != nil {
		return accessToken, err
	}
	accessToken.HuskyToken = token
	accessToken.URL = validatedURL
	accessToken.IsValid = true
	accessToken.CreatedAt = tH.External.GetTimeNow()
	if err := tH.External.StoreAccessToken(accessToken); err != nil {
		return types.AccessToken{}, err
	}
	return accessToken, nil
}

func (tH TokenHandler) ValidateToken(token, repositoryURL string) error {
	accessToken, err := tH.External.FindAccessToken(token, repositoryURL)
	if err != nil {
		return err
	}
	if !accessToken.IsValid {
		return errors.New("Access token is invalid")
	}
	return nil
}

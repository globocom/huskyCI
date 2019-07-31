package token

import (
	"github.com/globocom/huskyCI/api/types"
	"time"
)

type ExternalCalls interface {
	ValidateURL(url string) (string, error)
	GenerateToken() (string, error)
	GetTimeNow() time.Time
	StoreAccessToken(accessToken types.AccessToken) error
	FindAccessToken(token, repositoryURL string) (types.AccessToken, error)
}

type TokenHandler struct {
	External ExternalCalls
}

type TokenCaller struct{}

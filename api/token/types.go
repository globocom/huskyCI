package token

import (
	"github.com/globocom/huskyCI/api/auth"
	"github.com/globocom/huskyCI/api/types"
	"time"
)

type ExternalCalls interface {
	ValidateURL(url string) (string, error)
	GenerateToken() (string, error)
	GetTimeNow() time.Time
	StoreAccessToken(accessToken types.DBToken) error
	FindAccessToken(id string) (types.DBToken, error)
	UpdateAccessToken(id string, accesstoken types.DBToken) error
	FindRepoURL(repositoryURL string) error
	GenerateUuid() string
	EncodeBase64(m string) string
	DecodeToStringBase64(encodedVal string) (string, error)
}

type TokenHandler struct {
	External ExternalCalls
	HashGen  auth.Pbkdf2Generator
}

type TokenCaller struct{}

type TokenInterface interface {
	GenerateAccessToken(repo types.TokenRequest) (string, error)
	ValidateToken(token, repositoryURL string) error
	VerifyRepo(repositoryURL string) error
}

type TokenValidator struct {
	TokenVerifier TokenInterface
}

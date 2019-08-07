package token

import (
	"github.com/globocom/huskyCI/api/auth"
	"github.com/globocom/huskyCI/api/types"
	"time"
)

// ExternalCalls defines a group of functions
// used for external calls and validate some
// necessary information about TokenHandler.
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

// TokenHandler is a struct used to handle with
// token generation, validation and deactivation.
// It implements TokenInterface interface.
type TokenHandler struct {
	External ExternalCalls
	HashGen  auth.Pbkdf2Generator
}

// TokenCaller implements ExternalCalls interface.
type TokenCaller struct{}

// TokenInterface is used to define functions that
// handle with access token management.
type TokenInterface interface {
	GenerateAccessToken(repo types.TokenRequest) (string, error)
	ValidateToken(token, repositoryURL string) error
	VerifyRepo(repositoryURL string) error
}

// TokenValidator is used to validate an access token
// using the defined functions TokenInterface
type TokenValidator struct {
	TokenVerifier TokenInterface
}

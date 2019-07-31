package token

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"time"
)

func (tC *TokenCaller) ValidateURL(url string) (string, error) {
	return util.CheckMaliciousRepoURL(url)
}

func generateRandomBytes() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return b, err
}

func (tC *TokenCaller) GenerateToken() (string, error) {
	b, err := generateRandomBytes()
	return base64.URLEncoding.EncodeToString(b), err
}

func (tC *TokenCaller) GetTimeNow() time.Time {
	return time.Now()
}

func (tC *TokenCaller) StoreAccessToken(accessToken types.AccessToken) error {
	return db.InsertAccessToken(accessToken)
}

func (tC *TokenCaller) FindAccessToken(token, repositoryURL string) (types.AccessToken, error) {
	aTokenQuery := map[string]interface{}{"huskytoken": token, "repositoryURL": repositoryURL}
	return db.FindOneAccessToken(aTokenQuery)
}

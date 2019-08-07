package token

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/google/uuid"
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

func (tC *TokenCaller) StoreAccessToken(accessToken types.DBToken) error {
	return db.InsertAccessToken(accessToken)
}

func (tC *TokenCaller) FindAccessToken(id string) (types.DBToken, error) {
	aTokenQuery := map[string]interface{}{"uuid": id}
	return db.FindOneAccessToken(aTokenQuery)
}

func (tC *TokenCaller) FindRepoURL(repositoryURL string) error {
	repoQuery := map[string]interface{}{"repositoryURL": repositoryURL, "isValid": true}
	_, err := db.FindOneAccessToken(repoQuery)
	return err
}

func (tC *TokenCaller) GenerateUuid() string {
	return uuid.New().String()
}

func (tC *TokenCaller) EncodeBase64(m string) string {
	return base64.URLEncoding.EncodeToString([]byte(m))
}

func (tC *TokenCaller) DecodeToStringBase64(encodedVal string) (string, error) {
	decodedVal, err := base64.URLEncoding.DecodeString(encodedVal)
	return string(decodedVal), err
}

func (tC *TokenCaller) UpdateAccessToken(id string, accesstoken types.DBToken) error {
	aTokenQuery := map[string]interface{}{"uuid": id}
	return db.UpdateOneDBAccessToken(aTokenQuery, accesstoken)
}

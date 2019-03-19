package util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/globocom/huskyci/api/log"
)

const (
	// CertFile contains the address for the API's TLS certificate.
	CertFile = "api/api-tls-cert.pem"
	// KeyFile contains the address for the API's TLS certificate key file.
	KeyFile = "api/api-tls-key.pem"
)

// NewClientTLS returns an http client with certificate authentication.
func NewClientTLS() (*http.Client, error) {

	// Tries to find system's certificate pool
	caCertPool, _ := x509.SystemCertPool() // #nosec - SystemCertPool tries to get local cert pool, if it fails, a new cer pool is created
	if caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}

	cert, err := ioutil.ReadFile(CertFile)
	if err != nil {
		log.Error("NewClientTLS", "UTIL", 4001, err)
	}
	if ok := caCertPool.AppendCertsFromPEM(cert); !ok {
		log.Error("NewClientTLS", "UTIL", 4002, err)
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:               tls.VersionTLS11,
				MaxVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
				InsecureSkipVerify:       false,
				RootCAs:                  caCertPool,
			},
		},
	}
	return client, nil
}

// HandleCmd will extract %GIT_REPO% and %GIT_BRANCH% from cmd and replace it with the proper repository URL.
func HandleCmd(repositoryURL, repositoryBranch, cmd string) string {
	replace1 := strings.Replace(cmd, "%GIT_REPO%", repositoryURL, -1)
	replace2 := strings.Replace(replace1, "%GIT_BRANCH%", repositoryBranch, -1)
	return replace2
}

// HandlePrivateSSHKey will extract %GIT_PRIVATE_SSH_KEY% from cmd and replace it with the proper private SSH key.
func HandlePrivateSSHKey(rawString string) string {
	cmdReplaced := strings.Replace(rawString, "GIT_PRIVATE_SSH_KEY", os.Getenv("GIT_PRIVATE_SSH_KEY"), -1)
	return cmdReplaced
}

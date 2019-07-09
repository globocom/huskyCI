package util

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/globocom/huskyCI/api/log"
)

const (
	// CertFile contains the address for the API's TLS certificate.
	CertFile = "api/api-tls-cert.pem"
	// KeyFile contains the address for the API's TLS certificate key file.
	KeyFile = "api/api-tls-key.pem"
)

// NewClient returns an http client.
func NewClient(httpsEnable bool) (*http.Client, error) {

	if httpsEnable {
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

	client := &http.Client{}
	return client, nil
}

// GetLastLine receives a string with multiple lines and returns it's last
func GetLastLine(s string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines[len(lines)-1]
}

// GetAllLinesButLast receives a string with multiple lines and returns all but the last line.
func GetAllLinesButLast(s string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	lines = lines[:len(lines)-1]
	return lines
}

// SanitizeSafetyJSON returns a sanitized string from Safety container logs.
// Safety might return a JSON with the "\" and "\"" characters, which needs to be sanitized to be unmarshalled correctly.
func SanitizeSafetyJSON(s string) string {
	s1 := strings.Replace(s, "\\", "\\\\", -1)
	s2 := strings.Replace(s1, "\\\"", "\\\\\"", -1)
	return s2
}

// AdjustWarningMessage returns the Safety Warning string that will be printed.
func AdjustWarningMessage(warningRaw string) string {
	warning := strings.Split(warningRaw, ":")
	if len(warning) > 1 {
		warning[1] = strings.Replace(warning[1], "safety_huskyci_analysis_requirements_raw.txt", "'requirements.txt'", -1)
		warning[1] = strings.Replace(warning[1], " unpinned", "Unpinned", -1)

		return (warning[1] + " huskyCI can check it if you pin it in a format such as this: \"mypacket==3.2.9\" :D")
	}

	return warningRaw
}

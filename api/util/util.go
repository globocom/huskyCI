package util

import (
	"bufio"
	"os"
	"strings"
)

const (
	// CertFile contains the address for the API's TLS certificate.
	CertFile = "api/api-tls-cert.pem"
	// KeyFile contains the address for the API's TLS certificate key file.
	KeyFile = "api/api-tls-key.pem"
)

// RequestResponse returns a map string interface to be used in echo JSON requests response.
func RequestResponse(success bool, errMsg string) map[string]interface{} {
	if errMsg == "" {
		return map[string]interface{}{
			"success": success,
			"error":   "none",
		}
	}
	return map[string]interface{}{
		"success": success,
		"error":   errMsg,
	}
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
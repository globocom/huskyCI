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

// HandleCmd will extract %GIT_REPO% and %GIT_BRANCH% from cmd and replace it with the proper repository URL.
func HandleCmd(repositoryURL, repositoryBranch, cmd string) string {
	if repositoryURL != "" && repositoryBranch != "" && cmd != "" {
		replace1 := strings.Replace(cmd, "%GIT_REPO%", repositoryURL, -1)
		replace2 := strings.Replace(replace1, "%GIT_BRANCH%", repositoryBranch, -1)
		return replace2
	}
	return ""
}

// HandlePrivateSSHKey will extract %GIT_PRIVATE_SSH_KEY% from cmd and replace it with the proper private SSH key.
func HandlePrivateSSHKey(rawString string) string {
	privKey := os.Getenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY")
	cmdReplaced := strings.Replace(rawString, "GIT_PRIVATE_SSH_KEY", privKey, -1)
	return cmdReplaced
}

// GetLastLine receives a string with multiple lines and returns it's last
func GetLastLine(s string) string {
	if s == "" {
		return ""
	}
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines[len(lines)-1]
}

// GetAllLinesButLast receives a string with multiple lines and returns all but the last line.
func GetAllLinesButLast(s string) []string {
	if s == "" {
		return []string{}
	}
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

// RemoveDuplicates remove duplicated itens from a slice.
func RemoveDuplicates(s []string) []string {
	mapS := make(map[string]string, len(s))
	i := 0
	for _, v := range s {
		if _, ok := mapS[v]; !ok {
			mapS[v] = v
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// CreateContainerName generates a name for a container based on project's URL and branch
// Example:
//  - Input: {repositoryURL: https://github.com/globocom/huskyCI.git, repositoryBranch: "master", imageName: "huskyci/enry"}
//  - Output: globocom_huskyCI_master_enry
func CreateContainerName(repositoryURL, repositoryBranch, image string) string {
	if repositoryURL != "" && repositoryBranch != "" && image != "" {
		imageSlices := strings.Split(image, "/")
		imageName := imageSlices[len(imageSlices)-1]
		url := strings.TrimSuffix(repositoryURL, ".git")
		urlSlices := strings.Split(url, "/")[3:]
		myStrings := []string{strings.Join(urlSlices, "_"), repositoryBranch, imageName}
		return strings.Join(myStrings, "_")
	}
	return ""
}

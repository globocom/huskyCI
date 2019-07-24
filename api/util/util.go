package util

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

const (
	// CertFile contains the address for the API's TLS certificate.
	CertFile = "api/api-tls-cert.pem"
	// KeyFile contains the address for the API's TLS certificate key file.
	KeyFile = "api/api-tls-key.pem"
)

// HandleCmd will extract %GIT_REPO%, %GIT_BRANCH% and %INTERNAL_DEP_URL% from cmd and replace it with the proper repository URL.
func HandleCmd(repositoryURL, repositoryBranch, internalDepURL, cmd string) string {
	if repositoryURL != "" && repositoryBranch != "" && cmd != "" {
		replace1 := strings.Replace(cmd, "%GIT_REPO%", repositoryURL, -1)
		replace2 := strings.Replace(replace1, "%GIT_BRANCH%", repositoryBranch, -1)
		replace3 := strings.Replace(replace2, "%INTERNAL_DEP_URL%", internalDepURL, -1)
		return replace3
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
	if s == "" {
		return ""
	}
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

// CheckMaliciousInput checks if an user's input is "malicious" or not
func CheckMaliciousInput(repository types.Repository) (string, error) {

	sanitiziedURL, err := CheckMaliciousRepoURL(repository.URL)
	if err != nil {
		return "", err
	}

	if err := CheckMaliciousRepoBranch(repository.Branch); err != nil {
		return "", err
	}

	if repository.InternalDepURL != "" {
		if err := CheckMaliciousRepoInternalDepURL(repository.InternalDepURL); err != nil {
			return "", err
		}
	}

	return sanitiziedURL, nil
}

// CheckMaliciousRepoURL verifies if a given URL is a git repository and returns the sanitizied string and its error
func CheckMaliciousRepoURL(repositoryURL string) (string, error) {

	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`
	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repositoryURL)
	if err != nil {
		log.Error("CheckMaliciousRepoURL", "UTIL", 1008, "Repository URL regexp ", err)
		return "", types.ErrorInternal
	}
	if !valid {
		log.Error("CheckMaliciousRepoURL", "UTIL", 1016, repositoryURL)
		return "", types.ErrorInvalidRepository
	}

	return r.FindString(repositoryURL), nil
}

// CheckMaliciousRepoBranch verifies if a given branch is "malicious" or not
func CheckMaliciousRepoBranch(repositoryBranch string) error {

	regexpBranch := `^[a-zA-Z0-9_\/.-]*$`
	valid, err := regexp.MatchString(regexpBranch, repositoryBranch)
	if err != nil {
		log.Error("CheckMaliciousRepoBranch", "UTIL", 1008, "Repository Branch regexp ", err)
		return types.ErrorInternal
	}
	if !valid {
		log.Error("CheckMaliciousRepoBranch", "UTIL", 1017, repositoryBranch)
		return types.ErrorInvalidBranch
	}

	return nil
}

// CheckMaliciousRepoInternalDepURL verifies if a given internal dependecy URL is "malicious" or not
func CheckMaliciousRepoInternalDepURL(repositoryInternalDepURL string) error {

	regexpInternalDepURL := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
	valid, err := regexp.MatchString(regexpInternalDepURL, repositoryInternalDepURL)
	if err != nil {
		log.Error("CheckMaliciousRepoInternalDepURL", "UTIL", 1008, "Repository Branch regexp ", err)
		return types.ErrorInternal
	}
	if !valid {
		log.Error("CheckMaliciousRepoInternalDepURL", "UTIL", 1021, repositoryInternalDepURL)
		return types.ErrorInvalidDependencyURL
	}

	return nil
}

// CheckMaliciousRID verifies if a given RID is "malicious" or not
func CheckMaliciousRID(RID string) error {

	regexpRID := `^[a-zA-Z0-9]*$`
	valid, err := regexp.MatchString(regexpRID, RID)
	if err != nil {
		log.Error("CheckMaliciousRID", "UTIL", 1008, "RID regexp ", err)
		return types.ErrorInternal
	}
	if !valid {
		log.Warning("CheckMaliciousRID", "UTIL", 107, RID)
		return types.ErrorInvalidRID
	}

	return nil
}

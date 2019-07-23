package util

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
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
func CheckMaliciousInput(repository types.Repository, c echo.Context) (string, error) {

	sanitiziedURL, err := CheckMaliciousRepoURL(repository.URL, c)
	if err != nil {
		return "", err
	}

	if err := CheckMaliciousRepoBranch(repository.Branch, c); err != nil {
		return "", err
	}

	if repository.InternalDepURL != "" {
		if err := CheckMaliciousRepoInternalDepURL(repository.InternalDepURL, c); err != nil {
			return "", err
		}
	}

	return sanitiziedURL, nil
}

// CheckMaliciousRepoURL verifies if a given URL is a git repository and returns the sanitizied string and its error
func CheckMaliciousRepoURL(repositoryURL string, c echo.Context) (string, error) {

	errorInternal := errors.New("internal error")
	errorInvalidRepository := errors.New("invalid repository url")
	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`

	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repositoryURL)
	if err != nil {
		log.Error("CheckMaliciousRepoURL", "UTIL", 1008, "Repository URL regexp ", err)
		return "", errorInternal
	}
	if !valid {
		log.Error("CheckMaliciousRepoURL", "UTIL", 1016, repositoryURL)
		return "", errorInvalidRepository
	}
	return r.FindString(repositoryURL), nil
}

// CheckMaliciousRepoBranch verifies if a given branch is "malicious" or not
func CheckMaliciousRepoBranch(repositoryBranch string, c echo.Context) error {

	errorInternal := errors.New("internal error")
	errorInvalidBranch := errors.New("invalid branch")
	regexpBranch := `^[a-zA-Z0-9_\/.-]*$`

	valid, err := regexp.MatchString(regexpBranch, repositoryBranch)
	if err != nil {
		log.Error("CheckMaliciousRepoBranch", "UTIL", 1008, "Repository Branch regexp ", err)
		return errorInternal
	}
	if !valid {
		log.Error("CheckMaliciousRepoBranch", "UTIL", 1017, repositoryBranch)
		return errorInvalidBranch
	}
	return nil
}

// CheckMaliciousRepoInternalDepURL verifies if a given internal dependecy URL is "malicious" or not
func CheckMaliciousRepoInternalDepURL(repositoryInternalDepURL string, c echo.Context) error {

	errorInternal := errors.New("internal error")
	errorInvalidDepURL := errors.New("invalid internal dependency url")
	regexpInternalDepURL := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`

	valid, err := regexp.MatchString(regexpInternalDepURL, repositoryInternalDepURL)
	if err != nil {
		log.Error("CheckMaliciousRepoInternalDepURL", "UTIL", 1008, "Repository Branch regexp ", err)
		return errorInternal
	}
	if !valid {
		log.Error("CheckMaliciousRepoInternalDepURL", "UTIL", 1021, repositoryInternalDepURL)
		return errorInvalidDepURL
	}
	return nil
}

// CheckMaliciousRID verifies if a given RID is "malicious" or not
func CheckMaliciousRID(RID string, c echo.Context) error {

	errorInternal := errors.New("internal error")
	errorInvalidRID := errors.New("invalid rid")
	regexpRID := `^[a-zA-Z0-9]*$`

	valid, err := regexp.MatchString(regexpRID, RID)
	if err != nil {
		log.Error("CheckMaliciousRID", "UTIL", 1008, "RID regexp ", err)
		return errorInternal
	}
	if !valid {
		log.Warning("CheckMaliciousRID", "UTIL", 107, RID)
		return c.JSON(http.StatusBadRequest, errorInvalidRID)
	}
	return nil
}

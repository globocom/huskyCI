// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"errors"
	"fmt"

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

const logInfoAnalysis = "ANALYSIS"
const logActionReceiveRequest = "ReceiveRequest"

// HandleCmd will extract %GIT_REPO%, %GIT_BRANCH% from cmd and replace it with the proper repository URL.
func HandleCmd(repositoryURL, repositoryBranch, cmd string) string {
	if repositoryURL != "" && repositoryBranch != "" && cmd != "" {
		replace1 := strings.Replace(cmd, "%GIT_REPO%", repositoryURL, -1)
		replace2 := strings.Replace(replace1, "%GIT_BRANCH%", repositoryBranch, -1)
		return replace2
	}
	return ""
}

// HandleGitURLSubstitution will extract GIT_SSH_URL and GIT_URL_TO_SUBSTITUTE from cmd and replace it with the SSH equivalent.
func HandleGitURLSubstitution(rawString string) string {
	gitSSHURL := os.Getenv("HUSKYCI_API_GIT_SSH_URL")
	gitURLToSubstitute := os.Getenv("HUSKYCI_API_GIT_URL_TO_SUBSTITUTE")

	if gitSSHURL == "" || gitURLToSubstitute == "" {
		gitSSHURL = "nil"
		gitURLToSubstitute = "nil"
	}
	cmdReplaced := strings.Replace(rawString, "%GIT_SSH_URL%", gitSSHURL, -1)
	cmdReplaced = strings.Replace(cmdReplaced, "%GIT_URL_TO_SUBSTITUTE%", gitURLToSubstitute, -1)

	return cmdReplaced
}

// HandlePrivateSSHKey will extract %GIT_PRIVATE_SSH_KEY% from cmd and replace it with the proper private SSH key.
func HandlePrivateSSHKey(rawString string) string {
	privKey := os.Getenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY")
	cmdReplaced := strings.Replace(rawString, "%GIT_PRIVATE_SSH_KEY%", privKey, -1)
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

// CheckValidInput checks if an user's input is "malicious" or not
func CheckValidInput(repository types.Repository, c echo.Context) (string, error) {

	sanitiziedURL, err := CheckMaliciousRepoURL(repository.URL)
	if err != nil {
		if sanitiziedURL == "" {
			log.Error(logActionReceiveRequest, logInfoAnalysis, 1016, repository.URL)
			reply := map[string]interface{}{"success": false, "error": "invalid repository URL"}
			return "", c.JSON(http.StatusBadRequest, reply)
		}
		log.Error(logActionReceiveRequest, logInfoAnalysis, 1008, "Repository URL regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return "", c.JSON(http.StatusInternalServerError, reply)
	}

	if err := CheckMaliciousRepoBranch(repository.Branch, c); err != nil {
		return "", err
	}

	return sanitiziedURL, nil
}

// CheckMaliciousRepoURL verifies if a given URL is a git repository and returns the sanitizied string and its error
func CheckMaliciousRepoURL(repositoryURL string) (string, error) {
	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`
	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repositoryURL)
	if err != nil {
		return "matchStringError", err
	}
	if !valid {
		errorMsg := fmt.Sprintf("Invalid URL format: %s", repositoryURL)
		return "", errors.New(errorMsg)
	}
	return r.FindString(repositoryURL), nil
}

// CheckMaliciousRepoBranch verifies if a given branch is "malicious" or not
func CheckMaliciousRepoBranch(repositoryBranch string, c echo.Context) error {
	regexpBranch := `^[a-zA-Z0-9_\/.-]*$`
	valid, err := regexp.MatchString(regexpBranch, repositoryBranch)
	if err != nil {
		log.Error(logActionReceiveRequest, logInfoAnalysis, 1008, "Repository Branch regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	if !valid {
		log.Error(logActionReceiveRequest, logInfoAnalysis, 1017, repositoryBranch)
		reply := map[string]interface{}{"success": false, "error": "invalid repository branch"}
		return c.JSON(http.StatusBadRequest, reply)
	}
	return nil
}

// CheckMaliciousRID verifies if a given RID is "malicious" or not
func CheckMaliciousRID(RID string, c echo.Context) error {
	regexpRID := `^[-a-zA-Z0-9]*$`
	valid, err := regexp.MatchString(regexpRID, RID)
	if err != nil {
		log.Error("GetAnalysis", logInfoAnalysis, 1008, "RID regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	if !valid {
		log.Warning("GetAnalysis", logInfoAnalysis, 107, RID)
		reply := map[string]interface{}{"success": false, "error": "invalid RID"}
		return c.JSON(http.StatusBadRequest, reply)
	}
	return nil
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

// EndOfTheDay returns the the time at the end of the day t.
func EndOfTheDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

// BeginningOfTheDay returns the the time at the beginning of the day t.
func BeginningOfTheDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// CountDigits returns the number of digits in an integer.
func CountDigits(i int) int {
	count := 0
	for i != 0 {
		i /= 10
		count = count + 1
	}

	return count
}

func banditCase(code string, lineNumber int) bool {
	lineNumberLength := CountDigits(lineNumber)
	splitCode := strings.Split(code, "\n")
	for _, codeLine := range splitCode {
		if len(codeLine) > 0 {
			codeLineNumber := codeLine[:lineNumberLength]
			if strings.Contains(codeLine, "#nohusky") && (codeLineNumber == strconv.Itoa(lineNumber)) {
				return true
			}
		}
	}
	return false
}

// VerifyNoHusky verifies if the code string is marked with the #nohusky tag.
func VerifyNoHusky(code string, lineNumber int, securityTool string) bool {
	m := map[string]types.NohuskyFunction{
		"Bandit": banditCase,
	}

	return m[securityTool](code, lineNumber)

}

// SliceContains returns true if a given value is present on the given slice
func SliceContains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

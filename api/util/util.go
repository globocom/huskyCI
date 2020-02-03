// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/log"
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

// SliceContains returns true if a given value is present on the given slice
func SliceContains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

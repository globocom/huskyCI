// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/labstack/echo"
)

// AnalysisCount is the endpoint that returns data about HuskyCI scans.
func AnalysisCount(c echo.Context) error {
	timeRange := strings.ToLower(c.Param("time_range"))
	result, err := db.AnalysisCountMetric(timeRange)
	if err != nil {
		if err == errors.New("invalid time_range type") {
			log.Warning("AnalysisCount", "STATS", 111, err, timeRange)
			reply := map[string]interface{}{"success": false, "error": "invalid time_range type"}
			return c.JSON(http.StatusBadRequest, reply)
		}
		log.Error("AnalysisCount", "STATS", 2017, "AnalysisCount", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	return c.JSON(http.StatusOK, result)
}

// LanguageCount is the endpoint that returns the counter for each language scanned.
func LanguageCount(c echo.Context) error {
	result, err := db.LanguageCountMetric()
	if err != nil {
		log.Error("LanguageCount", "STATS", 2017, "LanguageCount", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	return c.JSON(http.StatusOK, result)
}

// RepositoryCount is the endpoint that returns the counter for each repository scanned.
func RepositoryCount(c echo.Context) error {
	result, err := db.RepositoryCountMetric()
	if err != nil {
		log.Error("RepositoryCount", "STATS", 2017, "RepositoryCount", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	return c.JSON(http.StatusOK, result)
}

// AuthorCount is the endpoint that returns the counter for each author from repositories scanned.
func AuthorCount(c echo.Context) error {
	result, err := db.AuthorCountMetric()
	if err != nil {
		log.Error("AuthorCount", "STATS", 2017, "AuthorCount", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	return c.JSON(http.StatusOK, result)
}

// ContainerCount is the endpoint that returns the counter for each container deployed.
func ContainerCount(c echo.Context) error {
	result, err := db.ContainerCountMetric()
	if err != nil {
		log.Error("ContainerCount", "STATS", 2017, "ContainerCount", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	return c.JSON(http.StatusOK, result)
}

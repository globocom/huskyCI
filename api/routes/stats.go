// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/patrickmn/go-cache"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
)

const logActionGetMetric = "GetMetric"
const logInfoStats = "STATS"

// GetMetric returns data about the metric received
func GetMetric(c echo.Context) error {
	url := c.Request().URL.String()
	metricType := strings.ToLower(c.Param("metric_type"))
	queryParams := c.QueryParams()

	if result, ok := apiContext.APIConfiguration.Cache.Get(url); ok {
		return c.JSON(http.StatusOK, result)
	}

	result, err := apiContext.APIConfiguration.DBInstance.GetMetricByType(metricType, queryParams)
	if err != nil {
		httpStatus, reply := checkError(err, metricType)
		return c.JSON(httpStatus, reply)
	}

	apiContext.APIConfiguration.Cache.Set(url, result, cache.DefaultExpiration)

	return c.JSON(http.StatusOK, result)
}

func checkError(err error, metricType string) (int, map[string]interface{}) {
	switch err.Error() {
	case "invalid time_range query string param":
		log.Warning(logActionGetMetric, logInfoStats, 111, err)
		reply := map[string]interface{}{"success": false, "error": "invalid time_range type"}
		return http.StatusBadRequest, reply
	case "invalid metric type":
		log.Warning(logActionGetMetric, logInfoStats, 112, metricType, err)
		reply := map[string]interface{}{"success": false, "error": "invalid metric type"}
		return http.StatusBadRequest, reply
	default:
		log.Error(logActionGetMetric, logInfoStats, 2017, metricType, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return http.StatusInternalServerError, reply
	}
}

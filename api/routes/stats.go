// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"
	"strings"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	"github.com/labstack/echo"
)

// GetMetric returns data about the metric received
func GetMetric(c echo.Context) error {
	metricType := strings.ToLower(c.Param("metric_type"))
	queryParams := c.QueryParams()
	result, err := apiContext.APIConfiguration.DbInstance.GetMetricByType(metricType, queryParams)
	if err != nil {
		httpStatus, reply := checkError(err, metricType)
		return c.JSON(httpStatus, reply)
	}
	return c.JSON(http.StatusOK, result)
}

func checkError(err error, metricType string) (int, map[string]interface{}) {
	switch err.Error() {
	case "invalid time_range query string param":
		log.Warning("GetMetric", "STATS", 111, err)
		reply := map[string]interface{}{"success": false, "error": "invalid time_range type"}
		return http.StatusBadRequest, reply
	case "invalid metric type":
		log.Warning("GetMetric", "STATS", 112, metricType, err)
		reply := map[string]interface{}{"success": false, "error": "invalid metric type"}
		return http.StatusBadRequest, reply
	default:
		log.Error("GetMetric", "STATS", 2017, metricType, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return http.StatusInternalServerError, reply
	}
}

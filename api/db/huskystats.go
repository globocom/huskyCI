// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

var statsQueryStringParams = map[string][]string{
	"language":    []string{"time_range"},
	"container":   []string{"time_range"},
	"analysis":    []string{"time_range"},
	"repository":  []string{"time_range"},
	"author":      []string{"time_range"},
	"severity":    []string{"time_range"},
	"time-to-fix": []string{"time_range"},
}

var aggrTimeFilterStage = map[string][]bson.M{
	"today":      generateTimeFilterStage(0, 0),
	"yesterday":  generateTimeFilterStage(-1, -1),
	"last7days":  generateTimeFilterStage(-6, 0),
	"last30days": generateTimeFilterStage(-29, 0),
}

// GetMetricByType returns data about the metric received
func GetMetricByType(metricType string, queryStringParams map[string][]string) (interface{}, error) {
	if !validMetric(metricType) {
		return nil, errors.New("invalid metric type")
	}
	validParams := validQueryStringParams(metricType, queryStringParams)
	err := validateParams(validParams)
	if err != nil {
		return nil, err
	}

	query := statsQueryBase[metricType]

	for param, values := range validParams {
		switch param {
		case "time_range":
			value := values[len(values)-1]
			query = append(aggrTimeFilterStage[value], query...)
		}
	}

	switch metricType {
	case "time-to-fix":
		return TimeToFixData(query)
	default:
		var obj interface{}
		err = mongoHuskyCI.Conn.Aggregation(query, mongoHuskyCI.AnalysisCollection, &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
}

// validTimeRange returns if a user inputted type is valid
func validTimeRange(timeRange string) bool {
	for tRange := range aggrTimeFilterStage {
		if timeRange == tRange {
			return true
		}
	}
	return false
}

// validMetric returns if a user inputted metric type is valid
func validMetric(metricType string) bool {
	for metric := range statsQueryBase {
		if metricType == metric {
			return true
		}
	}
	return false
}

// validateParams returns error if theres an invalid parameter
func validateParams(params map[string][]string) error {
	for param, values := range params {
		switch param {
		case "time_range":
			value := values[len(values)-1]
			if !validTimeRange(value) {
				return errors.New("invalid time_range query string param")
			}
		}
	}
	return nil
}

// validQueryStringParams returns a list of valid query string params
func validQueryStringParams(metric string, params map[string][]string) map[string][]string {
	validParams := make(map[string][]string)
	for key, value := range params {
		if util.SliceContains(statsQueryStringParams[metric], key) {
			validParams[key] = value
		}
	}
	return validParams
}

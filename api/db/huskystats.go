// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"time"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

var statsQueryStringParams = map[string][]string{
	"language":   []string{"time_range"},
	"container":  []string{"time_range"},
	"analysis":   []string{"time_range"},
	"repository": []string{"time_range"},
	"author":     []string{"time_range"},
}

var statsQueryBase = map[string][]bson.M{
	"language":  generateSimpleAggr("codes", "language", "codes.language"),
	"container": generateSimpleAggr("containers", "container", "containers.securityTest.name"),
	"analysis": []bson.M{
		bson.M{
			"$group": bson.M{
				"_id": "$result",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"count":  1,
				"result": "$_id",
			},
		},
	},
	"repository": []bson.M{
		bson.M{
			"$match": bson.M{
				"repositoryURL": bson.M{
					"$exists": true,
				},
			},
		},
		bson.M{
			"$match": bson.M{
				"repositoryBranch": bson.M{
					"$exists": true,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"repositoryBranch": "$repositoryBranch",
					"repositoryURL":    "$repositoryURL",
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"repositoryURL": "$_id.repositoryURL",
				},
				"branches": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "repositories",
				"totalBranches": bson.M{
					"$sum": "$branches",
				},
				"totalRepositories": bson.M{
					"$sum": 1,
				},
			},
		},
	},
	"author": []bson.M{
		bson.M{
			"$project": bson.M{
				"commitAuthors": 1,
			},
		},
		bson.M{
			"$unwind": "$commitAuthors",
		},
		bson.M{
			"$group": bson.M{
				"_id": "$commitAuthors",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "commitAuthors",
				"totalAuthors": bson.M{
					"$sum": 1,
				},
			},
		},
	},
}

var aggrTimeFilterStage = map[string][]bson.M{
	"today":      generateTimeFilterStage(-1, 0),
	"yesterday":  generateTimeFilterStage(-2, -1),
	"last7days":  generateTimeFilterStage(-7, 0),
	"last30days": generateTimeFilterStage(-30, 0),
}

// GetMetricByType returns data about the metric received
func GetMetricByType(metricType string, queryStringParams map[string][]string) (interface{}, error) {
	if !validMetric(metricType) {
		return nil, fmt.Errorf("invalid metric type")
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

	return mongoHuskyCI.Conn.Aggregation(query, mongoHuskyCI.AnalysisCollection)
}

// generateSimpleAggr generates an aggregation that counts each field group.
func generateSimpleAggr(field, finalName, groupID string) []bson.M {
	return []bson.M{
		bson.M{
			"$project": bson.M{
				field: 1,
			},
		},
		bson.M{
			"$unwind": fmt.Sprintf("$%s", field),
		},
		bson.M{
			"$group": bson.M{
				"_id": fmt.Sprintf("$%s", groupID),
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$project": bson.M{
				finalName: "$_id",
				"count":   1,
			},
		},
	}
}

// generateTimeFilterStage generates a stage that filter records by time range
func generateTimeFilterStage(rangeInitDays, rangeEndDays int) []bson.M {
	return []bson.M{
		bson.M{
			"$match": bson.M{
				"finishedAt": bson.M{
					"$gte": util.BeginningOfTheDay(time.Now().AddDate(0, 0, rangeInitDays)),
					"$lte": util.EndOfTheDay(time.Now().AddDate(0, 0, rangeEndDays)),
				},
			},
		},
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
				return fmt.Errorf("invalid time_range query string param")
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

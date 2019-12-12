// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"fmt"
	"time"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

const timeRangeQS = "time_range"

var statsQueryStringParams = map[string][]string{
	"language":        []string{timeRangeQS},
	"container":       []string{timeRangeQS},
	"analysis":        []string{timeRangeQS},
	"repository":      []string{timeRangeQS},
	"author":          []string{timeRangeQS},
	"severity":        []string{timeRangeQS},
	"historyanalysis": []string{timeRangeQS},
}

const aggHour = 1000 * 60 * 60

var statsQueryBase = map[string][]bson.M{
	"language":  generateSimpleAggr("codes", "language", "codes.language"),
	"container": generateSimpleAggr("containers", "container", "containers.securityTest.name"),
	"analysis": []bson.M{
		bson.M{
			"$project": bson.M{
				"finishedAt": 1,
				"result":     1,
			},
		},
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
	"severity": []bson.M{
		bson.M{
			"$project": bson.M{
				"huskyresults": bson.M{
					"$objectToArray": "$huskyciresults",
				},
			},
		},
		bson.M{
			"$unwind": "$huskyresults",
		},
		bson.M{
			"$project": bson.M{
				"languageresults": bson.M{
					"$objectToArray": "$huskyresults.v",
				},
			},
		},
		bson.M{
			"$unwind": "$languageresults",
		},
		bson.M{
			"$project": bson.M{
				"results": bson.M{
					"$objectToArray": "$languageresults.v",
				},
			},
		},
		bson.M{
			"$unwind": "$results",
		},
		bson.M{
			"$group": bson.M{
				"_id": "$results.k",
				"count": bson.M{
					"$sum": bson.M{
						"$size": "$results.v",
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"severity": "$_id",
				"count":    1,
			},
		},
	},
	"historyanalysis": []bson.M{
		bson.M{
			"$project": bson.M{
				"result": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []string{
								"$result",
								"warning",
							},
						},
						"then": "passed",
						"else": "$result",
					},
				},
				"finishedAt": 1,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dateNumber": bson.M{
					"$toLong": "$finishedAt",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dateMod": bson.M{
					"$mod": []interface{}{
						"$dateNumber",
						aggHour,
					},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"aggDate": bson.M{
					"$toDate": bson.M{
						"$subtract": []string{
							"$dateNumber",
							"$dateMod",
						},
					},
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"date":   "$aggDate",
					"result": "$result",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$_id.date",
				"results": bson.M{
					"$push": bson.M{
						"result": "$_id.result",
						"count":  "$count",
					},
				},
				"total": bson.M{
					"$sum": "$count",
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"_id": -1,
			},
		},
		bson.M{
			"$project": bson.M{
				"date":    "$_id",
				"_id":     0,
				"total":   1,
				"results": "$results",
			},
		},
	},
}

var validAggrTimeFilterStages = []string{"today", "yesterday", "last7days", "last30days"}

// GetMetricByType returns data about the metric received
func (mR *MongoRequests) GetMetricByType(metricType string, queryStringParams map[string][]string) (interface{}, error) {
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
		case timeRangeQS:
			value := values[len(values)-1]
			aggrTimeFilterStage := getTimeFilterStage(value)
			if aggrTimeFilterStage != nil {
				query = append(aggrTimeFilterStage, query...)
			}
		}
	}
	return mongoHuskyCI.Conn.Aggregation(query, mongoHuskyCI.AnalysisCollection)
}

func getTimeFilterStage(timeRange string) []bson.M {
	switch timeRange {
	case "today":
		return generateTimeFilterStage(0, 0)
	case "yesterday":
		return generateTimeFilterStage(-1, -1)
	case "last7days":
		return generateTimeFilterStage(-6, 0)
	case "last30days":
		return generateTimeFilterStage(-29, 0)
	default:
		return nil
	}
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
	for _, tRange := range validAggrTimeFilterStages {
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
		case timeRangeQS:
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

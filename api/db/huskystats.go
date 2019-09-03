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

var runCountAggregations = map[string]interface{}{
	"today":      generateTBAggr(-1, 0),
	"yesterday":  generateTBAggr(-2, -1),
	"last7days":  generateTBAggr(-7, 0),
	"last30days": generateTBAggr(-30, 0),
}

// AnalysisCountMetric returns predefined stats about huskyCI runs.
func AnalysisCountMetric(timeRange string) (interface{}, error) {
	if !validTimeRange(timeRange) {
		return nil, errors.New("invalid time_range type")
	}
	aggr := runCountAggregations[timeRange].([]bson.M)
	return mongoHuskyCI.Conn.Aggregation(aggr, mongoHuskyCI.AnalysisCollection)
}

// LanguageCountMetric returns the counter for each language scanned.
func LanguageCountMetric() (interface{}, error) {
	return mongoHuskyCI.Conn.Aggregation(generateSimpleAggr("codes", "language", "codes.language"), mongoHuskyCI.AnalysisCollection)
}

// RepositoryCountMetric returns the counter for each repository scanned.
func RepositoryCountMetric() (interface{}, error) {
	var repositoryCountAggr = []bson.M{
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
	}
	return mongoHuskyCI.Conn.Aggregation(repositoryCountAggr, mongoHuskyCI.AnalysisCollection)
}

// AuthorCountMetric returns the counter for each author from repositories scanned.
func AuthorCountMetric() (interface{}, error) {
	var AuthorCountAggr = []bson.M{
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
	}
	return mongoHuskyCI.Conn.Aggregation(AuthorCountAggr, mongoHuskyCI.AnalysisCollection)
}

// ContainerCountMetric returns the counter for each container deployed.
func ContainerCountMetric() (interface{}, error) {
	return mongoHuskyCI.Conn.Aggregation(generateSimpleAggr("containers", "container", "containers.securityTest.name"), mongoHuskyCI.AnalysisCollection)
}

// validTimeRange returns if a user inputted type is valid.
func validTimeRange(timeRange string) bool {
	for tRange := range runCountAggregations {
		if timeRange == tRange {
			return true
		}
	}
	return false
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

// generateTBAggr generates a time based aggregation in the given time range in days.
func generateCountAggr(field, groupID string) []bson.M {
	return []bson.M{
		bson.M{
			"$project": bson.M{
				field: 1,
			},
		},
		bson.M{
			"$match": bson.M{
				field: bson.M{
					"$exists": true,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": groupID,
				"count": bson.M{
					"$sum": bson.M{
						"$size": fmt.Sprintf("$%s", field),
					},
				},
			},
		},
	}
}

// generateTBAggr generates a time based aggregation in the given time range in days.
func generateTBAggr(rangeInitDays, rangeEndDays int) []bson.M {
	return []bson.M{
		bson.M{
			"$match": bson.M{
				"finishedAt": bson.M{
					"$gte": util.BeginningOfTheDay(time.Now().AddDate(0, 0, rangeInitDays)),
					"$lte": util.EndOfTheDay(time.Now().AddDate(0, 0, rangeEndDays)),
				},
			},
		}, bson.M{
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
	}
}

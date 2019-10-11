// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"fmt"
	"sort"
	"time"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

type timeToFixData struct {
	TotalAnalyses    int        `bson:"totalAnalyses"`
	RepositoryURL    string     `bson:"repositoryURL"`
	RepositoryBranch string     `bson:"repositoryBranch"`
	Analyses         []analysis `bson:"analyses"`
}

type analysis struct {
	StartedAt      time.Time `bson:"startedAt"`
	FinishedAt     time.Time `bson:"finishedAt"`
	Result         string    `bson:"result"`
	HuskyCIResults bson.M    `bson:"huskyciresults"`
}

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
				"_id":    0,
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
				"_id":      0,
			},
		},
	},
	"time-to-fix": []bson.M{
		bson.M{
			"$project": bson.M{
				"repositoryBranch": 1,
				"startedAt":        1,
				"finishedAt":       1,
				"result":           1,
				"repositoryURL":    1,
				"huskyciresults": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []string{"{}", "$huskyciresults"},
						},
						"then": "$REMOVE",
						"else": "$huskyciresults",
					},
				},
			},
		}, bson.M{
			"$match": bson.M{
				"result": bson.M{
					"$in": []string{"passed", "failed"},
				},
			},
		}, bson.M{
			"$unwind": bson.M{
				"path":                       "$startedAt",
				"preserveNullAndEmptyArrays": false,
			},
		}, bson.M{
			"$unwind": bson.M{
				"path":                       "$finishedAt",
				"preserveNullAndEmptyArrays": false,
			},
		}, bson.M{
			"$unwind": bson.M{
				"path":                       "$huskyciresults",
				"preserveNullAndEmptyArrays": false,
			},
		}, bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"repositoryURL":    "$repositoryURL",
					"repositoryBranch": "$repositoryBranch",
				},
				"analyses": bson.M{
					"$push": bson.M{
						"startedAt":      "$startedAt",
						"finishedAt":     "$finishedAt",
						"huskyciresults": "$huskyciresults",
						"result":         "$result",
					},
				},
				"results": bson.M{
					"$addToSet": "$result",
				},
				"totalAnalyses": bson.M{
					"$sum": 1,
				},
			},
		}, bson.M{
			"$match": bson.M{
				"results": bson.M{
					"$size": 2,
				},
			},
		}, bson.M{
			"$project": bson.M{
				"_id":              0,
				"repositoryURL":    "$_id.repositoryURL",
				"repositoryBranch": "$_id.repositoryBranch",
				"analyses":         1,
				"totalAnalyses":    1,
			},
		},
	},
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
				"_id":     0,
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

// TimeToFixData queries MongoDB and parse the result to consolidate time to fix vulnerability data
func TimeToFixData(query []bson.M) (interface{}, error) {
	var aggrResult []timeToFixData
	err := mongoHuskyCI.Conn.Aggregation(query, mongoHuskyCI.AnalysisCollection, &aggrResult)
	if err != nil {
		return nil, err
	}
	if len(aggrResult) > 0 {
		return computeTimeToFixResults(aggrResult)
	}

	return bson.M{
			"averageTimeToFix": nil,
			"totalAnalyses":    0,
			"totalSolutions":   0,
		},
		nil
}

func computeTimeToFixResults(data []timeToFixData) (bson.M, error) {
	totalAnalysed := 0
	prev := &analysis{}
	currentlyFailed := false
	var timeToFix []time.Duration
	countSolutions := 0
	for _, repo := range data {

		totalAnalysed += repo.TotalAnalyses
		analyses := repo.Analyses
		sort.Slice(analyses, func(i, j int) bool {
			return repo.Analyses[i].FinishedAt.Before(repo.Analyses[j].FinishedAt)
		})
		for _, analysis := range analyses {
			switch analysis.Result {
			case "failed":
				if !currentlyFailed {
					currentlyFailed = true
					*prev = analysis
				}
			case "passed":
				if currentlyFailed {
					currentlyFailed = false
					if analysis.FinishedAt.After(prev.FinishedAt) {
						timeToFix = append(timeToFix, analysis.FinishedAt.Sub(prev.FinishedAt))
						countSolutions++
					}
				}
			default:
				return nil, errors.New("unrecognized analysis result")
			}
		}
	}

	var totalTimeNanoseconds int64

	for _, fixTime := range timeToFix {
		totalTimeNanoseconds += fixTime.Nanoseconds()
	}

	if len(timeToFix) > 0 && totalTimeNanoseconds > 0 {
		averageTimeToFix := time.Duration(totalTimeNanoseconds / int64(len(timeToFix)))
		return bson.M{
				"averageTimeToFix": averageTimeToFix.String(),
				"totalAnalyses":    totalAnalysed,
				"totalSolutions":   countSolutions,
			},
			nil
	}
	return bson.M{
			"averageTimeToFix": nil,
			"totalAnalyses":    totalAnalysed,
			"totalSolutions":   countSolutions,
		},
		nil

}

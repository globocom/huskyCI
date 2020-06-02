package database

import (
	"github.com/globocom/huskyCI/api/analysis"
	"gopkg.in/mgo.v2/bson"
)

// InsertAnalysis insters a new analysis into AnalysisCollection
func (m *MongoDB) InsertAnalysis(analysis *analysis.Analysis) error {
	newAnalysis := bson.M{
		"ID":         analysis.ID,
		"repository": analysis.Repository,
		"branch":     analysis.Branch,
		"startedAt":  analysis.StartedAt,
	}
	return m.Insert(newAnalysis, "Analysis")
}

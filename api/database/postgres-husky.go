package database

import (
	"github.com/globocom/huskyCI/api/analysis"
)

// InsertAnalysis insters a new analysis into Postgres DB
func (p *Postgres) InsertAnalysis(analysis *analysis.Analysis) error {
	return nil
}

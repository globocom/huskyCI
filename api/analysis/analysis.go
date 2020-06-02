package analysis

import (
	"time"

	"github.com/google/uuid"
)

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID         string    `bson:"ID" json:"ID"`
	Repository string    `bson:"repository" json:"repository"`
	Branch     string    `bson:"branch" json:"branch"`
	StartedAt  time.Time `bson:"startedAt" json:"startedAt"`
	FinishedAt time.Time `bson:"finishedAt" json:"finishedAt"`
	// Vulnerabilities []vulnerability.Vulnerability `bson:"vulnerabilities" json:"vulnerabilities"`
	// SecurityTests   []*securitytest.SecurityTest  `bson:"securityTests" json:"securityTests"`
}

// New returns a new analysis struct based on a repository
func New(repository, branch string) *Analysis {
	return &Analysis{
		ID:         uuid.New().String(),
		Repository: repository,
		Branch:     branch,
		StartedAt:  time.Now(),
	}
}

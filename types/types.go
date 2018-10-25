package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// JSONPayload is a struct that represents the JSON payload needed to make a Husky API request.
type JSONPayload struct {
	RepositoryURL    string `json:"repositoryURL"`
	RepositoryBranch string `json:"repositoryBranch"`
}

// JSONResponse is a struct that represents the JSON reponse of a Husky API request.
type JSONResponse struct {
	RID     string `json:"RID"`
	Details string `json:"details"`
	Result  string `json:"result"`
}

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID            bson.ObjectId  `bson:"_id,omitempty"`
	RID           string         `bson:"RID" json:"RID"`
	URL           string         `bson:"repositoryURL" json:"repositoryURL"`
	Branch        string         `bson:"repositoryBranch" json:"repositoryBranch"`
	SecurityTests []SecurityTest `bson:"securityTests" json:"securityTests"`
	Status        string         `bson:"status" json:"status"`
	Result        string         `bson:"result" json:"result"`
	Containers    []Container    `bson:"containers" json:"containers"`
}

// Container is the struct that stores all data from a container run.
type Container struct {
	CID          string       `bson:"CID" json:"CID"`
	VM           string       `bson:"VM" json:"VM"`
	SecurityTest SecurityTest `bson:"securityTest" json:"securityTest"`
	CStatus      string       `bson:"cStatus" json:"cStatus"`
	COutput      string       `bson:"cOutput" json:"cOutput"`
	CResult      string       `bson:"cResult" json:"cResult"`
	StartedAt    time.Time    `bson:"startedAt" json:"startedAt"`
	FinishedAt   time.Time    `bson:"finishedAt" json:"finishedAt"`
}

// SecurityTest is the struct that stores all data from the security tests to be executed.
type SecurityTest struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Name             string        `bson:"name" json:"name"`
	Image            string        `bson:"image" json:"image"`
	Cmd              string        `bson:"cmd" json:"cmd"`
	Language         string        `bson:"language" json:"language"`
	Default          bool          `bson:"default" json:"default"`
	TimeOutInSeconds int           `bson:"timeOutSeconds" json:"timeOutSeconds"`
}

// GasOutput is the struct that holds all data from the output of the Gas tool.
type GasOutput struct {
	Issues []GasIssue `json:"Issues"`
	Stats  GasStats   `json:"Stats"`
}

// GasIssue is the struct that holds all issues from the output of the Gas tool.
type GasIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// GasStats is the struct that holds all stats from the output of the Gas tool.
type GasStats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	Nosec int `json:"nosec"`
	Found int `json:"found"`
}

package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// JSONPayload is a struct that represents the JSON payload needed to make a Husky API request.
type JSONPayload struct {
	RepositoryURL string `json:"repositoryURL"`
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
	URL           string         `bson:"URL" json:"URL"`
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
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name" json:"name"`
	Image    string        `bson:"image" json:"image"`
	Cmd      string        `bson:"cmd" json:"cmd"`
	Language string        `bson:"language" json:"language"`
	Default  bool          `bson:"default" json:"default"`
}

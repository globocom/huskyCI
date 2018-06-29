package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Repository is the struct that stores all data from repository to be analyzed.
type Repository struct {
	ID               bson.ObjectId  `bson:"_id,omitempty"`
	URL              string         `bson:"URL" json:"repositoryURL"`
	SecurityTest     []SecurityTest `bson:"securityTest" json:"securityTest"`
	SecurityTestName string         `json:"securityTestName"`
	VM               string         `bson:"VM" json:"vm"`
	CreatedAt        time.Time      `bson:"createdAt" json:"createdAt"`
	DeletedAt        time.Time      `bson:"deletedAt" json:"deletedAt"`
	Language         string         `bson:"language" json:"language"`
}

// SecurityTest is the struct that stores all data from the security tests to be executed.
type SecurityTest struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name" json:"name"`
	Image    string        `bson:"image" json:"image"`
	Cmd      []string      `bson:"cmd" json:"cmd"`
	Language string        `bson:"language" json:"language"`
	Default  bool          `bson:"default" json:"default"`
}

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	RID            string          `bson:"RID" json:"RID"`
	URL            string          `bson:"URL" json:"URL"`
	SecurityTestID []bson.ObjectId `bson:"securityTest" json:"securityTestID"`
	Status         string          `bson:"status" json:"status"`
	Result         string          `bson:"result" json:"result"`
	Container      []Container     `bson:"container" json:"container"`
}

// Container is the struct that stores all data from a container run.
type Container struct {
	CID            string        `bson:"CID" json:"CID"`
	VM             string        `bson:"VM" json:"VM"`
	SecurityTestID bson.ObjectId `bson:"securityTest" json:"securityTestID"`
	CStatus        string        `bson:"cStatus" json:"cStatus"`
	COuput         []string      `bson:"cOutput" json:"cOutput"`
	StartedAt      time.Time     `bson:"startedAt" json:"startedAt"`
	FinishedAt     time.Time     `bson:"finishedAt" json:"finishedAt"`
}

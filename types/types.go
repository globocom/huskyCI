package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Repository is the struct of all data from repository to be analyzed.
type Repository struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	URL            string          `json:"repositoryURL" bson:"URL"`
	VM             string          `bson:"VM"`
	SecurityTestID []bson.ObjectId `bson:"securityTest"`
	CreatedAt      time.Time       `bson:"createdAt"`
	DeletedAt      time.Time       `bson:"deletedAt"`
}

// SecurityTest is the struct of all data from the security tests to be executed.
type SecurityTest struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name" json:"name"`
	Image string        `bson:"image" json:"image"`
	Cmd   []string      `bson:"cmd" json:"cmd"`
}

// Analysis is the struct of all data from analysis performed.
type Analysis struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	RID            string          `bson:"RID"`
	URL            string          `bson:"URL"`
	SecurityTestID []bson.ObjectId `bson:"securityTest"`
	Status         string          `bson:"status"`
	Result         string          `bson:"result"`
	Output         []string        `bson:"output"`
	CID            []string        `bson:"container"`
}

// Container is the struct of all data from a container run.
type Container struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	CID            string        `bson:"CID"`
	RID            string        `bson:"RID"`
	VM             string        `bson:"VM"`
	SecurityTestID bson.ObjectId `bson:"securityTest"`
	CStatus        string        `bson:"cStatus"`
	COuput         []string      `bson:"cOutput"`
}

package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Repository is the struct of all data from repository to be analyzed.
type Repository struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	URL            string          `bson:"URL" json:"repositoryURL"`
	VM             string          `bson:"VM" json:"vm"`
	SecurityTestID []bson.ObjectId `bson:"securityTest" json:"securityTestID"`
	Language       string          `bson:"language" json:"language"`
	CreatedAt      time.Time       `bson:"createdAt" json:"createdAt"`
	DeletedAt      time.Time       `bson:"deletedAt" json:"deletedAt"`
}

// SecurityTest is the struct of all data from the security tests to be executed.
type SecurityTest struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name" json:"name"`
	Image    string        `bson:"image" json:"image"`
	Cmd      []string      `bson:"cmd" json:"cmd"`
	Language string        `bson:"language" json:"language"`
	Default  bool          `bson:"default" json:"default"`
}

// Analysis is the struct of all data from analysis performed.
type Analysis struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	RID            string          `bson:"RID" json:"RID"`
	URL            string          `bson:"URL" json:"URL"`
	SecurityTestID []bson.ObjectId `bson:"securityTest" json:"securityTestID"`
	Status         string          `bson:"status" json:"status"`
	Result         string          `bson:"result" json:"result"`
	CID            []string        `bson:"container" json:"container"`
}

// Container is the struct of all data from a container run.
type Container struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	CID            string        `bson:"CID" json:"CID"`
	RID            string        `bson:"RID" json:"RID"`
	VM             string        `bson:"VM" json:"VM"`
	SecurityTestID bson.ObjectId `bson:"securityTest" json:"securityTestID"`
	CStatus        string        `bson:"cStatus" json:"cStatus"`
	COuput         []string      `bson:"cOutput" json:"cOutput"`
}

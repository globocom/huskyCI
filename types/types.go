package types

import "gopkg.in/mgo.v2/bson"

// Repository is the struct of all data from repository to be analyzed.
type Repository struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	URL            string          `json:"repositoryURL" bson:"URL"`
	VM             string          `bson:"VM"`
	SecurityTestID []bson.ObjectId `bson:"securityTest"`
	CreatedAt      string          `bson:"createdAt"`
	DeletedAt      string          `bson:"deletedAt"`
}

// SecurityTest is the struct of all data from the security tests to be executed.
type SecurityTest struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name" json:"securityTestName"`
	Image string        `bson:"image"`
	Cmd   []string      `bson:"cmd" json:"cmd"`
}

// Analysis is the struct of all data from analysis performed.
type Analysis struct {
	RID            bson.ObjectId   `bson:"_id,omitempy"`
	URL            string          `bson:"URL"`
	SecurityTestID []bson.ObjectId `bson:"securityTest"`
	Status         string          `bson:"status"`
	Result         string          `bson:"result"`
	Output         []string        `bson:"output"`
	Container      []bson.ObjectId `bson:"containers"`
}

// Container is the struct of all data from a container run.
type Container struct {
	CID            bson.ObjectId `bson:"_id,omitempy"`
	RID            bson.ObjectId `bson:"RID"`
	VM             string        `bson:"VM"`
	SecurityTestID bson.ObjectId `bson:"securityTest"`
	CStatus        string        `bson:"cStatus"`
	COuput         string        `bson:"cOutput"`
}

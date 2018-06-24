package types

import "gopkg.in/mgo.v2/bson"

// Repository is the struct of all data from repository to be analyzed.
type Repository struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	URL          string        `json:"repositoryURL" bson:"URL"`
	VM           string        `bson:"VM"`
	SecurityTest []string      `bson:"securityTest"`
	CreatedAt    string        `bson:"createdAt"`
	DeletedAt    string        `bson:"deletedAt"`
}

// SecurityTest is the struct of all data from the security tests to be executed.
type SecurityTest struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Image string        `bson:"image"`
	Cmd   []string      `bson:"cmd"`
}

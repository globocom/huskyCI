package types

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Repository is the struct that stores all data from repository to be analyzed.
type Repository struct {
	ID               bson.ObjectId  `bson:"_id,omitempty"`
	URL              string         `bson:"repositoryURL" json:"repositoryURL"`
	Branch           string         `bson:"repositoryBranch" json:"repositoryBranch"`
	SecurityTests    []SecurityTest `bson:"securityTests" json:"securityTests"`
	SecurityTestName []string       `bson:"securityTestName,omitempty" json:"securityTestName"`
	VM               string         `bson:"VM" json:"vm"`
	CreatedAt        time.Time      `bson:"createdAt" json:"createdAt"`
	DeletedAt        time.Time      `bson:"deletedAt" json:"deletedAt"`
	Languages        []Language     `bson:"languages" json:"languages"`
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
	COuput       string       `bson:"cOutput" json:"cOutput"`
	CResult      string       `bson:"cResult" json:"cResult"`
	StartedAt    time.Time    `bson:"startedAt" json:"startedAt"`
	FinishedAt   time.Time    `bson:"finishedAt" json:"finishedAt"`
}

// Language is the struct that stores all data from a language's repository.
type Language struct {
	Name  string   `bson:"name" json:"language_name"`
	Files []string `bson:"files" json:"language_files"`
}

// VersionAPI is the struct that stores all data about version api.
type VersionAPI struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (v VersionAPI) Print() {
	version := v.getValue(v.Version, "N/A")
	commit := v.getValue(v.Commit, "N/A")
	date := v.getValue(v.Date, "N/A")

	printVersion := fmt.Sprintf(`
************************************************
project: %s
version: %s
commit: %s
data build: %s
************************************************
	`, v.Project, version, commit, date)

	fmt.Println(printVersion)
}

func (v VersionAPI) getValue(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

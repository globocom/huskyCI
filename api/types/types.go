// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Repository is the struct that stores all data from repository to be analyzed.
type Repository struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	URL            string        `bson:"repositoryURL" json:"repositoryURL"`
	Branch         string        `json:"repositoryBranch"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	InternalDepURL string        `json:"internaldepURL"`
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
	ID             bson.ObjectId  `bson:"_id,omitempty"`
	RID            string         `bson:"RID" json:"RID"`
	URL            string         `bson:"repositoryURL" json:"repositoryURL"`
	Branch         string         `bson:"repositoryBranch" json:"repositoryBranch"`
	SecurityTests  []SecurityTest `bson:"securityTests" json:"securityTests"`
	Status         string         `bson:"status" json:"status"`
	Result         string         `bson:"result" json:"result"`
	Containers     []Container    `bson:"containers" json:"containers"`
	StartedAt      time.Time      `bson:"startedAt" json:"startedAt"`
	FinishedAt     time.Time      `bson:"finishedAt" json:"finishedAt"`
	InternalDepURL string         `bson:"internaldepURL,omitempty" json:"internaldepURL"`
	Codes          []Code         `bson:"codes" json:"codes"`
	HuskyCIResults HuskyCIResults `bson:"huskyciresults,omitempty" json:"huskyciresults"`
}

// Container is the struct that stores all data from a container run.
type Container struct {
	CID          string       `bson:"CID" json:"CID"`
	SecurityTest SecurityTest `bson:"securityTest" json:"securityTest"`
	CStatus      string       `bson:"cStatus" json:"cStatus"`
	COutput      string       `bson:"cOutput" json:"cOutput"`
	CResult      string       `bson:"cResult" json:"cResult"`
	CInfo        string       `bson:"cInfo" json:"cInfo"`
	StartedAt    time.Time    `bson:"startedAt" json:"startedAt"`
	FinishedAt   time.Time    `bson:"finishedAt" json:"finishedAt"`
}

// Code is the struct that stores all data from code found in a repository.
type Code struct {
	Language string   `bson:"language" json:"language"`
	Files    []string `bson:"files" json:"files"`
}

// HuskyCIResults is a struct that represents huskyCI scan results.
type HuskyCIResults struct {
	GoResults         GoResults         `bson:"goresults,omitempty" json:"goresults,omitempty"`
	PythonResults     PythonResults     `bson:"pythonresults,omitempty" json:"pythonresults,omitempty"`
	JavaScriptResults JavaScriptResults `bson:"javascriptresults,omitempty" json:"javascriptresults,omitempty"`
	RubyResults       RubyResults       `bson:"rubyresults,omitempty" json:"rubyresults,omitempty"`
}

// GoResults represents all Golang security tests results.
type GoResults struct {
	HuskyCIGosecOutput HuskyCIGosecOutput `bson:"gosecoutput,omitempty" json:"gosecoutput,omitempty"`
}

// PythonResults represents all Python security tests results.
type PythonResults struct {
	HuskyCIBanditOutput HuskyCIBanditOutput    `bson:"banditoutput,omitempty" json:"banditoutput,omitempty"`
	HuskyCISafetyOutput []HuskyCIVulnerability `bson:"safetyoutput,omitempty" json:"safetyoutput,omitempty"`
}

// JavaScriptResults represents all JavaScript security tests results.
type JavaScriptResults struct {
	HuskyCIRetireJSOutput HuskyCIRetireJSOutput `bson:"retirejsoutput,omitempty" json:"retirejsoutput,omitempty"`
	HuskyCINpmAuditOutput HuskyCINpmAuditOutput `bson:"npmauditoutput,omitempty" json:"npmauditoutput,omitempty"`
}

// RubyResults represents all Ruby security tests results.
type RubyResults struct {
	HuskyCIBrakemanOutput HuskyCIBrakemanOutput `bson:"brakemanoutput,omitempty" json:"brakemanoutput,omitempty"`
}

// HuskyCIGosecOutput stores all Low, Medium and High Gosec vulnerabilities
type HuskyCIGosecOutput struct {
	LowVulnsGosec    []HuskyCIVulnerability `bson:"lowvulnsgosec,omitempty" json:"lowvulnsbandit,omitempty"`
	MediumVulnsGosec []HuskyCIVulnerability `bson:"mediumvulnsgosec,omitempty" json:"mediumvulnsbandit,omitempty"`
	HighVulnsGosec   []HuskyCIVulnerability `bson:"highvulnsgosec,omitempty" json:"highvulnsbandit,omitempty"`
}

// HuskyCIBanditOutput stores all Low, Medium and High Bandit vulnerabilities
type HuskyCIBanditOutput struct {
	LowVulnsBandit    []HuskyCIVulnerability `bson:"lowvulnsbandit,omitempty" json:"lowvulnsbandit,omitempty"`
	MediumVulnsBandit []HuskyCIVulnerability `bson:"mediumvulnsbandit,omitempty" json:"mediumvulnsbandit,omitempty"`
	HighVulnsBandit   []HuskyCIVulnerability `bson:"highvulnsbandit,omitempty" json:"highvulnsbandit,omitempty"`
}

// HuskyCIBrakemanOutput stores all Low, Medium and High Brakeman vulnerabilities
type HuskyCIBrakemanOutput struct {
	LowVulnsBrakeman    []HuskyCIVulnerability `bson:"lowvulnsbrakeman,omitempty" json:"lowvulnsbrakeman,omitempty"`
	MediumVulnsBrakeman []HuskyCIVulnerability `bson:"mediumvulnsbrakeman,omitempty" json:"mediumvulnsbrakeman,omitempty"`
	HighVulnsBrakeman   []HuskyCIVulnerability `bson:"highvulnsbrakeman,omitempty" json:"highvulnsbrakeman,omitempty"`
}

// HuskyCINpmAuditOutput stores all Low, Medium and High Npm Audit vulnerabilities
type HuskyCINpmAuditOutput struct {
	LowVulnsNpmAudit    []HuskyCIVulnerability `bson:"lowvulnsnpmaudit,omitempty" json:"lowvulnsnpmaudit,omitempty"`
	MediumVulnsNpmAudit []HuskyCIVulnerability `bson:"mediumvulnsnpmaudit,omitempty" json:"mediumvulnsnpmaudit,omitempty"`
	HighVulnsNpmAudit   []HuskyCIVulnerability `bson:"highvulnsnpmaudit,omitempty" json:"highvulnsnpmaudit,omitempty"`
}

// HuskyCIRetireJSOutput stores all Low, Medium and High RetireJS vulnerabilities
type HuskyCIRetireJSOutput struct {
	LowVulnsNpmRetireJS []HuskyCIVulnerability `bson:"lowvulnsretireJS,omitempty" json:"lowvulnsretireJS,omitempty"`
	MediumVulnsRetireJS []HuskyCIVulnerability `bson:"mediumvulnsretireJS,omitempty" json:"mediumvulnsretireJS,omitempty"`
	HighVulnsRetireJS   []HuskyCIVulnerability `bson:"highvulnsretireJS,omitempty" json:"highvulnsretireJS,omitempty"`
}

// HuskyCIVulnerability is the struct that stores vulnerability information.
type HuskyCIVulnerability struct {
	Language       string `bson:"language" json:"language,omitempty"`
	SecurityTool   string `bson:"securitytool" json:"securitytool,omitempty"`
	Severity       string `bson:"severity,omitempty" json:"severity,omitempty"`
	Confidence     string `bson:"confidence,omitempty" json:"confidence,omitempty"`
	File           string `bson:"file,omitempty" json:"file,omitempty"`
	Line           string `bson:"line,omitempty" json:"line,omitempty"`
	Code           string `bson:"code,omitempty" json:"code,omitempty"`
	Details        string `bson:"details" json:"details,omitempty"`
	Type           string `bson:"type,omitempty" json:"type,omitempty"`
	VunerableBelow string `bson:"vulnerablebelow,omitempty" json:"vulnerablebelow,omitempty"`
	Version        string `bson:"version,omitempty" json:"version,omitempty"`
	Occurrences    int    `bson:"occurrences,omitempty" json:"occurrences,omitempty"`
}

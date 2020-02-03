// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"sync"
	"time"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/repository"
	"github.com/globocom/huskyCI/api/securitytest"
	"github.com/globocom/huskyCI/api/vulnerability"
	"github.com/google/uuid"
)

// Analysis is the struct that stores all data from analysis performed.
// new
type Analysis struct {
	ID              string                        `bson:"ID" json:"ID"`
	StartedAt       time.Time                     `bson:"startedAt" json:"startedAt"`
	FinishedAt      time.Time                     `bson:"finishedAt" json:"finishedAt"`
	ErrorsFound     []string                      `bson:"errorsFound,omitempty" json:"errorsFound"`
	Result          Result                        `bson:"result,omitempty" json:"result"`
	Repository      *repository.Repository        `bson:"repository" json:"repository"`
	Vulnerabilities []vulnerability.Vulnerability `bson:"vulnerabilities" json:"vulnerabilities"`
	SecurityTests   []*securitytest.SecurityTest  `bson:"securityTests" json:"securityTests"`
}

// Result holds the status and the info of an analysis.
type Result struct {
	Status string `bson:"status" json:"status"`
	Info   string `bson:"info,omitempty" json:"info"`
}

// New returns a new analysis struct based on a repository
func New(repository *repository.Repository) *Analysis {
	return &Analysis{
		ID:         uuid.New().String(),
		Repository: repository,
		StartedAt:  time.Now(),
	}
}

// Start runs a new analysis
func (a *Analysis) Start() error {

	log.Info("Start", "ANALYSIS", 101, a.ID)

	if err := a.Repository.Scan(); err != nil {
		return err
	}

	// if err := a.checkCacheHit(); err != nil {
	// 	return err
	// }

	a.setSecurityTests()

	if err := a.startSecurityTests(); err != nil {
		return err
	}

	if err := a.RegisterInDB(); err != nil {
		return err
	}

	log.Info("Start", "ANALYSIS", 102, a.ID)

	return nil
}

// RegisterInDB register in database the current analysis
func (a *Analysis) RegisterInDB() error {

	// if err := apiContext.APIConfiguration.DBInstance.InsertDBAnalysis(a); err != nil {
	// 	log.Error("registerInDatabase", "ANALYSIS", 2011, err)
	// 	return err
	// }

	return nil
}

// CheckResult queries the database to check the result of an analysis
func (a *Analysis) CheckResult() error {

	// if err := apiContext.APIConfiguration.DBInstance.InsertDBAnalysis(a); err != nil {
	// 	log.Error("registerInDatabase", "ANALYSIS", 2011, err)
	// 	return err
	// }

	return nil
}

func (a *Analysis) checkCacheHit() error {

	// var cacheHit bool

	// if cacheHit {
	// 	if err := a.registerInDatabase(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (a *Analysis) setSecurityTests() {

	var allSecurityTests []*securitytest.SecurityTest

	allSecurityTests = append(allSecurityTests, securitytest.GetAllGeneric()...)

	for _, language := range a.Repository.Languages {
		languageSecTest := securitytest.GetAllByLanguage(language)
		allSecurityTests = append(allSecurityTests, languageSecTest...)
	}

	a.SecurityTests = allSecurityTests
}

func (a *Analysis) startSecurityTests() error {

	// run and analyze all securityTests in parallel
	var wg sync.WaitGroup

	errChan := make(chan error)
	waitChan := make(chan struct{})
	syncChan := make(chan struct{})

	defer close(errChan)

	for _, securityTest := range a.SecurityTests {

		wg.Add(1)

		go func(secTest *securitytest.SecurityTest) {
			defer wg.Done()
			if err := secTest.Container.Run(a.Repository.URL, a.Repository.Branch); err != nil {
				select {
				case <-syncChan:
					return
				case errChan <- err:
					return
				}
			}
			if err := secTest.Analyze(); err != nil {
				select {
				case <-syncChan:
					return
				case errChan <- err:
					return
				}
			}
		}(securityTest)

	}

	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		return nil
	case err := <-errChan:
		close(syncChan)
		return err
	}

}

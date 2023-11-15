// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"testing"
	"time"

	mongoHuskyCI "github.com/globocom/huskyCI/api/db/mongo"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
)

var huskydbMongoRequestsTest MongoRequests

func TestMongoRequestsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	huskydbMongoRequestsTest = MongoRequests{}
	log.InitLog(true, "", "", "log_test", "log_test")
	RegisterFailHandler(Fail)
	RunSpecs(t, "MongoDB Suite")
}

var _ = BeforeSuite(func() {
	mongoAddress := "localhost"
	dbName := "integration-test"
	username := ""
	password := ""
	dbPort := 27017
	connectionTimeout := time.Duration(1 * time.Second)
	connectionPool := 10
	maxOpenConns := 10
	maxIdleConns := 10
	connMaxLifetime := time.Duration(1 * time.Second)

	errConnect := huskydbMongoRequestsTest.ConnectDB(mongoAddress, dbName, username, password, connectionTimeout, connectionPool, dbPort, maxOpenConns, maxIdleConns, connMaxLifetime)
	Expect(errConnect).To(BeNil())
})

var _ = AfterSuite(func() {
	err := mongoHuskyCI.Conn.Session.DB("").DropDatabase()
	Expect(err).To(BeNil())
})

var _ = Describe("FindOneDBRepository", func() {
	Context("When find one DB Repository", func() {
		It("Should return mgo.ErrNotFound and empty types.Repository{} when not found", func() {
			mapParams := map[string]interface{}{"repositoryURL": "not found URL"}

			repository, err := huskydbMongoRequestsTest.FindOneDBRepository(mapParams)
			Expect(err).To(Equal(mgo.ErrNotFound))

			expectedResult := types.Repository{}
			Expect(repository).To(Equal(expectedResult))
		})
		It("Should return no error and types.Repository{} correctly", func() {
			repositoryToInsert := types.Repository{
				URL: "http://github.com/findonerepository",
			}

			errInsert := huskydbMongoRequestsTest.InsertDBRepository(repositoryToInsert)
			Expect(errInsert).To(BeNil())

			mapParams := map[string]interface{}{"repositoryURL": repositoryToInsert.URL}

			repository, errGet := huskydbMongoRequestsTest.FindOneDBRepository(mapParams)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(repositoryToInsert))
		})
	})
})

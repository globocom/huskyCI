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
	Context("When try to find one", func() {
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

var _ = Describe("FindOneDBSecurityTest", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.SecurityTest{} when not found", func() {
			mapParams := map[string]interface{}{"name": "security_test_name"}

			repository, err := huskydbMongoRequestsTest.FindOneDBSecurityTest(mapParams)
			Expect(err).To(Equal(mgo.ErrNotFound))

			expectedResult := types.SecurityTest{}
			Expect(repository).To(Equal(expectedResult))
		})
		It("Should return no error and  types.SecurityTest{} correctly", func() {
			securityTestToInsert := types.SecurityTest{
				Name:             "security_test_name",
				Image:            "some image",
				Cmd:              "some cmd",
				Language:         "some language",
				Type:             "some tipe",
				Default:          false,
				TimeOutInSeconds: 10,
			}

			errInsert := huskydbMongoRequestsTest.InsertDBSecurityTest(securityTestToInsert)
			Expect(errInsert).To(BeNil())

			mapParams := map[string]interface{}{"name": "security_test_name"}

			repository, errGet := huskydbMongoRequestsTest.FindOneDBSecurityTest(mapParams)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(securityTestToInsert))
		})
	})
})

var _ = Describe("FindOneDBAnalysis", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.Analysis{} when not found", func() {
			mapParams := map[string]interface{}{"RID": "test-id"}

			analysis, err := huskydbMongoRequestsTest.FindOneDBAnalysis(mapParams)
			Expect(err).To(Equal(mgo.ErrNotFound))

			expectedResult := types.Analysis{}
			Expect(analysis).To(Equal(expectedResult))
		})
		It("Should return no error and types.Analysis{} correctly", func() {
			analysisToInsert := types.Analysis{
				RID:        "test-id",
				URL:        "some url",
				Branch:     "some branch",
				Status:     "some status",
				Result:     "some result",
				Containers: []types.Container{},
			}

			errInsert := huskydbMongoRequestsTest.InsertDBAnalysis(analysisToInsert)
			Expect(errInsert).To(BeNil())

			mapParams := map[string]interface{}{"RID": "test-id"}

			analysis, errGet := huskydbMongoRequestsTest.FindOneDBAnalysis(mapParams)
			Expect(errGet).To(BeNil())
			Expect(analysis).To(Equal(analysisToInsert))
		})
	})
})

var _ = Describe("FindOneDBUser", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.User{} when not found", func() {
			mapParams := map[string]interface{}{"username": "some user name"}

			user, err := huskydbMongoRequestsTest.FindOneDBUser(mapParams)
			Expect(err).To(Equal(mgo.ErrNotFound))

			expectedResult := types.User{}
			Expect(user).To(Equal(expectedResult))
		})
		It("Should return no error and types.User{} correctly", func() {
			userToInsert := types.User{
				Username:     "some user name",
				Password:     "some password",
				Salt:         "some salt",
				Iterations:   1,
				KeyLen:       10,
				HashFunction: "some hash function",
			}

			errInsert := huskydbMongoRequestsTest.InsertDBUser(userToInsert)
			Expect(errInsert).To(BeNil())

			mapParams := map[string]interface{}{"username": "some user name"}

			user, errGet := huskydbMongoRequestsTest.FindOneDBUser(mapParams)
			Expect(errGet).To(BeNil())
			Expect(user).To(Equal(userToInsert))
		})
	})
})

var _ = Describe("FindOneDBAccessToken", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.DBToken{} when not found", func() {
			mapParams := map[string]interface{}{"uuid": "some uuid"}

			user, err := huskydbMongoRequestsTest.FindOneDBAccessToken(mapParams)
			Expect(err).To(Equal(mgo.ErrNotFound))

			expectedResult := types.DBToken{}
			Expect(user).To(Equal(expectedResult))
		})
		It("Should return no error and types.DBToken{} correctly", func() {
			dbToken := types.DBToken{
				HuskyToken: "some token",
				URL:        "some url",
				IsValid:    true,
				Salt:       "some salt",
				UUID:       "some uuid",
			}

			errInsert := huskydbMongoRequestsTest.InsertDBAccessToken(dbToken)
			Expect(errInsert).To(BeNil())

			mapParams := map[string]interface{}{"uuid": "some uuid"}

			user, errGet := huskydbMongoRequestsTest.FindOneDBAccessToken(mapParams)
			Expect(errGet).To(BeNil())
			Expect(user).To(Equal(dbToken))
		})
	})
})

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
	"gopkg.in/mgo.v2/bson"
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

var _ = Describe("DBRepository", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty types.Repository{} when not found", func() {
			query := map[string]interface{}{"repositoryURL": "not found URL"}

			repository, err := huskydbMongoRequestsTest.FindOneDBRepository(query)
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

			query := map[string]interface{}{"repositoryURL": repositoryToInsert.URL}

			repository, errGet := huskydbMongoRequestsTest.FindOneDBRepository(query)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(repositoryToInsert))
		})
	})
	Context("When try to find all", func() {
		It("Should return no error and empty repository array when not found", func() {
			query := map[string]interface{}{"repositoryURL": "not found URL"}

			repository, err := huskydbMongoRequestsTest.FindAllDBRepository(query)
			Expect(err).To(BeNil())

			expectedResult := []types.Repository{}
			Expect(repository).To(Equal(expectedResult))
		})
		It("Should return no error and types.Repository{} correctly", func() {
			repositoryToInsert1 := types.Repository{
				URL: "http://github.com/findallrepository",
			}

			errInsert := huskydbMongoRequestsTest.InsertDBRepository(repositoryToInsert1)
			Expect(errInsert).To(BeNil())

			repositoryToInsert2 := types.Repository{
				URL: "http://github.com/findallrepository",
			}

			errInsert = huskydbMongoRequestsTest.InsertDBRepository(repositoryToInsert2)
			Expect(errInsert).To(BeNil())

			query := map[string]interface{}{
				"repositoryURL": "http://github.com/findallrepository",
			}

			expectedResult := []types.Repository{
				repositoryToInsert1,
				repositoryToInsert2,
			}

			result, errGet := huskydbMongoRequestsTest.FindAllDBRepository(query)
			Expect(errGet).To(BeNil())
			Expect(result).To(Equal(expectedResult))
		})
	})
	Context("When update one", func() {
		It("Should return no error and update correctly", func() {
			repositoryToInsert := types.Repository{
				URL:       "http://github.com/updateOneRepository",
				CreatedAt: time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC).Local(),
			}

			errInsert := huskydbMongoRequestsTest.InsertDBRepository(repositoryToInsert)
			Expect(errInsert).To(BeNil())

			repositoryToUpdate := types.Repository{
				URL:       "http://github.com/updateOneRepository",
				CreatedAt: time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC).Local(),
			}

			query := map[string]interface{}{"repositoryURL": repositoryToInsert.URL}
			updateQuery := map[string]interface{}{
				"repositoryURL": repositoryToUpdate.URL,
				"createdAt":     repositoryToUpdate.CreatedAt,
			}
			errUpdate := huskydbMongoRequestsTest.UpdateOneDBRepository(query, updateQuery)
			Expect(errUpdate).To(BeNil())

			repository, errGet := huskydbMongoRequestsTest.FindOneDBRepository(query)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(repositoryToUpdate))
		})
	})
})

var _ = Describe("DBSecurityTest", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.SecurityTest{} when not found", func() {
			query := map[string]interface{}{"name": "security_test_name"}

			repository, err := huskydbMongoRequestsTest.FindOneDBSecurityTest(query)
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
				Type:             "some type",
				Default:          false,
				TimeOutInSeconds: 10,
			}

			errInsert := huskydbMongoRequestsTest.InsertDBSecurityTest(securityTestToInsert)
			Expect(errInsert).To(BeNil())

			query := map[string]interface{}{"name": "security_test_name"}

			repository, errGet := huskydbMongoRequestsTest.FindOneDBSecurityTest(query)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(securityTestToInsert))
		})
	})
	Context("When try to find all", func() {
		It("Should return no error and empty types.SecurityTest{} array when not found", func() {
			query := map[string]interface{}{"type": "type_find_all"}

			repository, err := huskydbMongoRequestsTest.FindAllDBSecurityTest(query)
			Expect(err).To(BeNil())

			expectedResult := []types.SecurityTest{}
			Expect(repository).To(Equal(expectedResult))
		})
		It("Should return no error and  types.SecurityTest{} correctly", func() {
			securityTestToInsert1 := types.SecurityTest{
				Name:             "security_test_name_1",
				Image:            "some image",
				Cmd:              "some cmd",
				Language:         "some language",
				Type:             "type_find_all",
				Default:          false,
				TimeOutInSeconds: 10,
			}

			errInsert := huskydbMongoRequestsTest.InsertDBSecurityTest(securityTestToInsert1)
			Expect(errInsert).To(BeNil())

			securityTestToInsert2 := types.SecurityTest{
				Name:             "security_test_name_2",
				Image:            "some image",
				Cmd:              "some cmd",
				Language:         "some language",
				Type:             "type_find_all",
				Default:          false,
				TimeOutInSeconds: 10,
			}

			errInsert = huskydbMongoRequestsTest.InsertDBSecurityTest(securityTestToInsert2)
			Expect(errInsert).To(BeNil())

			expectResult := []types.SecurityTest{
				securityTestToInsert1,
				securityTestToInsert2,
			}

			query := map[string]interface{}{"type": "type_find_all"}
			repository, errGet := huskydbMongoRequestsTest.FindAllDBSecurityTest(query)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(expectResult))
		})
	})
	Context("When updated one", func() {
		It("Should return no error and update correctly", func() {
			securityTestToInsert := types.SecurityTest{
				Name:             "security_test_update_one",
				Image:            "some image",
				Cmd:              "some cmd",
				Language:         "some language",
				Type:             "some type",
				Default:          false,
				TimeOutInSeconds: 10,
			}

			errInsert := huskydbMongoRequestsTest.InsertDBSecurityTest(securityTestToInsert)
			Expect(errInsert).To(BeNil())

			securityTestToUpdate := types.SecurityTest{
				Name:             "security_test_update_one",
				Image:            "changed image",
				Cmd:              "changed cmd",
				Language:         "changed language",
				Type:             "changed type",
				Default:          true,
				TimeOutInSeconds: 20,
			}

			query := map[string]interface{}{"name": "security_test_update_one"}
			_, errUpdate := huskydbMongoRequestsTest.UpsertOneDBSecurityTest(query, securityTestToUpdate)
			Expect(errUpdate).To(BeNil())

			repository, errGet := huskydbMongoRequestsTest.FindOneDBSecurityTest(query)
			Expect(errGet).To(BeNil())
			Expect(repository).To(Equal(securityTestToUpdate))
		})
	})
})

var _ = Describe("DBAnalysis", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.Analysis{} when not found", func() {
			query := map[string]interface{}{"RID": "test-id"}

			analysis, err := huskydbMongoRequestsTest.FindOneDBAnalysis(query)
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

			query := map[string]interface{}{"RID": "test-id"}

			analysis, errGet := huskydbMongoRequestsTest.FindOneDBAnalysis(query)
			Expect(errGet).To(BeNil())
			Expect(analysis).To(Equal(analysisToInsert))
		})
	})
	Context("When try to find all", func() {
		It("Should return no error and empty types.Analysis{} array when not found", func() {
			query := map[string]interface{}{"status": "status-find-all-not-found"}

			analysis, err := huskydbMongoRequestsTest.FindAllDBAnalysis(query)
			Expect(err).To(BeNil())

			expectedResult := []types.Analysis{}
			Expect(analysis).To(Equal(expectedResult))
		})
		It("Should return no error and []types.Analysis{} correctly", func() {
			analysisToInsert1 := types.Analysis{
				RID:        "test-id-1",
				URL:        "some url",
				Branch:     "some branch",
				Status:     "status-find-all",
				Result:     "some result",
				Containers: []types.Container{},
			}

			errInsert := huskydbMongoRequestsTest.InsertDBAnalysis(analysisToInsert1)
			Expect(errInsert).To(BeNil())

			analysisToInsert2 := types.Analysis{
				RID:        "test-id-2",
				URL:        "some url",
				Branch:     "some branch",
				Status:     "status-find-all",
				Result:     "some result",
				Containers: []types.Container{},
			}

			errInsert = huskydbMongoRequestsTest.InsertDBAnalysis(analysisToInsert2)
			Expect(errInsert).To(BeNil())

			expectedResult := []types.Analysis{
				analysisToInsert1,
				analysisToInsert2,
			}

			query := map[string]interface{}{"status": "status-find-all"}
			analysis, errGet := huskydbMongoRequestsTest.FindAllDBAnalysis(query)
			Expect(errGet).To(BeNil())
			Expect(analysis).To(Equal(expectedResult))
		})
	})
	Context("When update one", func() {
		It("Should return no error and update correctly", func() {
			analysisToInsert := types.Analysis{
				RID:        "test-id-update-one",
				URL:        "some url",
				Branch:     "some branch",
				Status:     "some status",
				Result:     "some result",
				Containers: []types.Container{},
			}

			errInsert := huskydbMongoRequestsTest.InsertDBAnalysis(analysisToInsert)
			Expect(errInsert).To(BeNil())

			analysisToUpdate := types.Analysis{
				RID:        "test-id-update-one",
				URL:        "changed url",
				Branch:     "some branch",
				Status:     "some status",
				Result:     "some result",
				Containers: []types.Container{},
			}

			query := map[string]interface{}{"RID": "test-id-update-one"}
			toUpdate := map[string]interface{}{"repositoryURL": analysisToUpdate.URL}

			errUpdate := huskydbMongoRequestsTest.UpdateOneDBAnalysis(query, toUpdate)
			Expect(errUpdate).To(BeNil())

			analysis, errGet := huskydbMongoRequestsTest.FindOneDBAnalysis(query)
			Expect(errGet).To(BeNil())
			Expect(analysis).To(Equal(analysisToUpdate))
		})
	})
	Context("When update one analysis container", func() {
		It("Should return no error and update correctly", func() {
			analysisToInsert := types.Analysis{
				RID:        "test-id-update-one-analysis-container",
				URL:        "some url",
				Branch:     "some branch",
				Status:     "some status",
				Result:     "some result",
				Containers: []types.Container{},
			}

			errInsert := huskydbMongoRequestsTest.InsertDBAnalysis(analysisToInsert)
			Expect(errInsert).To(BeNil())

			analysisToUpdate := types.Analysis{
				RID:        "test-id-update-one-analysis-container",
				URL:        "changed url",
				Branch:     "some branch",
				Status:     "some status",
				Result:     "some result",
				Containers: []types.Container{},
			}

			query := map[string]interface{}{"RID": "test-id-update-one-analysis-container"}
			toUpdate := map[string]interface{}{"repositoryURL": analysisToUpdate.URL}

			errUpdate := huskydbMongoRequestsTest.UpdateOneDBAnalysisContainer(query, toUpdate)
			Expect(errUpdate).To(BeNil())

			analysis, errGet := huskydbMongoRequestsTest.FindOneDBAnalysis(query)
			Expect(errGet).To(BeNil())
			Expect(analysis).To(Equal(analysisToUpdate))
		})
	})
})

var _ = Describe("DBUser", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.User{} when not found", func() {
			query := map[string]interface{}{"username": "some user name"}

			user, err := huskydbMongoRequestsTest.FindOneDBUser(query)
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

			query := map[string]interface{}{"username": "some user name"}

			user, errGet := huskydbMongoRequestsTest.FindOneDBUser(query)
			Expect(errGet).To(BeNil())
			Expect(user).To(Equal(userToInsert))
		})
	})
	Context("When update one", func() {
		It("Should return no error and update correctly", func() {
			userToInsert := types.User{
				Username:     "update_one_username",
				Password:     "some password",
				Salt:         "some salt",
				Iterations:   1,
				KeyLen:       10,
				HashFunction: "some hash function",
			}

			errInsert := huskydbMongoRequestsTest.InsertDBUser(userToInsert)
			Expect(errInsert).To(BeNil())

			userToUpdate := types.User{
				Username:     "update_one_username",
				Password:     "changed password",
				Salt:         "changed salt",
				Iterations:   2,
				KeyLen:       20,
				HashFunction: "changed hash function",
			}

			query := map[string]interface{}{"username": "update_one_username"}
			errUpdate := huskydbMongoRequestsTest.UpdateOneDBUser(query, userToUpdate)
			Expect(errUpdate).To(BeNil())

			user, errGet := huskydbMongoRequestsTest.FindOneDBUser(query)
			Expect(errGet).To(BeNil())
			Expect(user).To(Equal(userToUpdate))
		})
	})
})

var _ = Describe("DBAccessToken", func() {
	Context("When try to find one", func() {
		It("Should return mgo.ErrNotFound and empty  types.DBToken{} when not found", func() {
			query := map[string]interface{}{"uuid": "some uuid"}

			user, err := huskydbMongoRequestsTest.FindOneDBAccessToken(query)
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

			query := map[string]interface{}{"uuid": "some uuid"}

			user, errGet := huskydbMongoRequestsTest.FindOneDBAccessToken(query)
			Expect(errGet).To(BeNil())
			Expect(user).To(Equal(dbToken))
		})
	})
	Context("When update one", func() {
		It("Should return no error and update correctly", func() {
			dbTokenToInsert := types.DBToken{
				HuskyToken: "some token",
				URL:        "some url",
				IsValid:    true,
				Salt:       "some salt",
				UUID:       "some_update_one_uuid",
			}

			errInsert := huskydbMongoRequestsTest.InsertDBAccessToken(dbTokenToInsert)
			Expect(errInsert).To(BeNil())

			dbTokenToUpdate := types.DBToken{
				HuskyToken: "changed token",
				URL:        "changed url",
				IsValid:    false,
				Salt:       "changed salt",
				UUID:       "some_update_one_uuid",
			}

			query := map[string]interface{}{"uuid": "some_update_one_uuid"}
			errUpdate := huskydbMongoRequestsTest.UpdateOneDBAccessToken(query, dbTokenToUpdate)
			Expect(errUpdate).To(BeNil())

			user, errGet := huskydbMongoRequestsTest.FindOneDBAccessToken(query)
			Expect(errGet).To(BeNil())
			Expect(user).To(Equal(dbTokenToUpdate))
		})
	})
})

var _ = Describe("DockerAPIAddresses", func() {
	Context("When update one", func() {
		It("Should return no error and modify correctly", func() {
			dockerAPIAddressesToInsert := types.DockerAPIAddresses{
				CurrentHostIndex: 0,
				HostList:         []string{"host1", "host2"},
			}

			errInsert := mongoHuskyCI.Conn.Insert(dockerAPIAddressesToInsert, mongoHuskyCI.DockerAPIAddressesCollection)
			Expect(errInsert).To(BeNil())

			updatedDockerAPIAddressesResult, err := huskydbMongoRequestsTest.FindAndModifyDockerAPIAddresses()
			Expect(err).To(BeNil())
			Expect(updatedDockerAPIAddressesResult).To(Equal(dockerAPIAddressesToInsert))

			dockerAPIAddressesFindResult := []types.DockerAPIAddresses{}
			errFind := mongoHuskyCI.Conn.Search(bson.M{}, nil, mongoHuskyCI.DockerAPIAddressesCollection, &dockerAPIAddressesFindResult)
			Expect(errFind).To(BeNil())

			expectedDockerApiAddress := []types.DockerAPIAddresses{
				{
					CurrentHostIndex: 1,
					HostList:         []string{"host1", "host2"},
				},
			}

			Expect(dockerAPIAddressesFindResult).To(Equal(expectedDockerApiAddress))
		})
	})
})

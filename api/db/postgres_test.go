package db_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	. "github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
)

type FakePsql struct {
	expectedRolesResult    []map[string]interface{}
	expectedGetValuesError error
}

func (fP *FakePsql) Connect(address string,
	username string,
	password string,
	dbName string,
	maxOpenConns int,
	maxIdleConns int,
	connLT time.Duration) error {
	return nil
}

func (fP *FakePsql) GetValuesFromDB(query string,
	args ...interface{}) ([]map[string]interface{}, error) {
	return fP.expectedRolesResult, fP.expectedGetValuesError
}

func (fP *FakePsql) WriteInDB(query string, args ...interface{}) (int64, error) {
	return 0, nil
}

type FakeJSON struct {
	expectedJsonValues     []byte
	expectedMarshalError   error
	expectedUnmarshalError error
	expectedRepository     types.Repository
}

func (fJ *FakeJSON) Marshal(v interface{}) ([]byte, error) {
	return fJ.expectedJsonValues, fJ.expectedMarshalError
}

func (fJ *FakeJSON) Unmarshal(data []byte, v interface{}) error {
	if fJ.expectedUnmarshalError == nil {
		newV := v.(*types.Repository)
		newV.URL = fJ.expectedRepository.URL
		newV.Branch = fJ.expectedRepository.Branch
		newV.CreatedAt = fJ.expectedRepository.CreatedAt
	}
	return fJ.expectedUnmarshalError
}

var _ = Describe("Postgres", func() {
	Describe("FindOneDBRepository", func() {
		Context("When key map verification returns false", func() {
			It("Should return an empty Repository with the expected error", func() {
				postgres := PostgresRequests{}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"Wrong Key": "Failed"})
				Expect(repo).To(Equal(types.Repository{}))
				Expect(err).To(Equal(errors.New("Could not find repository URL")))
			})
		})
		Context("When GetValuesFromDB returns an error", func() {
			It("Should return an empty Repository with the same error", func() {
				fakePsql := FakePsql{
					expectedRolesResult:    nil,
					expectedGetValuesError: errors.New("Failed to get values"),
				}
				postgres := PostgresRequests{
					Psql: &fakePsql,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(types.Repository{}))
				Expect(err).To(Equal(fakePsql.expectedGetValuesError))
			})
		})
		Context("When GetValuesFromDB returns more than one entry", func() {
			It("Should return an empty Repository with the same error", func() {
				fakePsql := FakePsql{
					expectedRolesResult: []map[string]interface{}{
						{"repositoryURL": "teste1", "repositoryBranch": "teste1"},
						{"repositoryURL": "teste2", "repositoryBranch": "teste2"},
						{"repositoryURL": "teste3", "repositoryBranch": "teste3"},
					},
					expectedGetValuesError: nil,
				}
				postgres := PostgresRequests{
					Psql: &fakePsql,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(types.Repository{}))
				Expect(err).To(Equal(errors.New("Data returned in a not expected format")))
			})
		})
		Context("When Marshal function returns an error", func() {
			It("Should return an empty Repository with the same error", func() {
				fakePsql := FakePsql{
					expectedRolesResult: []map[string]interface{}{
						{"repositoryURL": "teste", "repositoryBranch": "teste"},
					},
					expectedGetValuesError: nil,
				}
				fakeJSON := FakeJSON{
					expectedMarshalError: errors.New("Failed trying to marshal in a JSON"),
				}
				postgres := PostgresRequests{
					Psql:        &fakePsql,
					JsonHandler: &fakeJSON,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(types.Repository{}))
				Expect(err).To(Equal(fakeJSON.expectedMarshalError))
			})
		})
		Context("When Unmarshal function returns an error", func() {
			It("Should return an empty Repository with the same error", func() {
				fakePsql := FakePsql{
					expectedRolesResult: []map[string]interface{}{
						{"repositoryURL": "teste", "repositoryBranch": "teste"},
					},
					expectedGetValuesError: nil,
				}
				fakeJSON := FakeJSON{
					expectedMarshalError:   nil,
					expectedUnmarshalError: errors.New("Failed trying to unmarshal"),
				}
				postgres := PostgresRequests{
					Psql:        &fakePsql,
					JsonHandler: &fakeJSON,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(types.Repository{}))
				Expect(err).To(Equal(fakeJSON.expectedUnmarshalError))
			})
		})
		Context("When Unmarshal function returns the valid Repository struct", func() {
			It("Should return the expected Repository with a nil error", func() {
				fakePsql := FakePsql{
					expectedRolesResult: []map[string]interface{}{
						{"repositoryURL": "teste", "repositoryBranch": "teste"},
					},
					expectedGetValuesError: nil,
				}
				fakeJSON := FakeJSON{
					expectedMarshalError:   nil,
					expectedUnmarshalError: nil,
					expectedRepository: types.Repository{
						URL:       "teste",
						Branch:    "teste",
						CreatedAt: time.Now(),
					},
				}
				postgres := PostgresRequests{
					Psql:        &fakePsql,
					JsonHandler: &fakeJSON,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(fakeJSON.expectedRepository))
				Expect(err).To(BeNil())
			})
		})
	})
})

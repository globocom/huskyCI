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
	})
})

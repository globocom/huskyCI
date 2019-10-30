package db_test

import (
	"errors"
	"time"

	. "github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeRetriever struct {
	expectedRetrieveError error
	expectedRepository    types.Repository
	expectedSecurityTest  types.SecurityTest
	expectedAnalysis      types.Analysis
	expectedUser          types.User
	expectedDBToken       types.DBToken
}

func (fR *FakeRetriever) Connect(
	address string,
	dbName string,
	username string,
	password string,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration) error {
	return nil
}

func (fR *FakeRetriever) RetrieveFromDB(
	query string, response interface{}, params ...interface{}) error {
	if fR.expectedRetrieveError == nil {
		switch response.(type) {
		case *[]types.Repository:
			newV := response.(*[]types.Repository)
			(*newV) = append((*newV), fR.expectedRepository)
		case *[]types.SecurityTest:
			newV := response.(*[]types.SecurityTest)
			(*newV) = append((*newV), fR.expectedSecurityTest)
		case *[]types.Analysis:
			newV := response.(*[]types.Analysis)
			(*newV) = append((*newV), fR.expectedAnalysis)
		case *[]types.User:
			newV := response.(*[]types.User)
			(*newV) = append((*newV), fR.expectedUser)
		case *[]types.DBToken:
			newV := response.(*[]types.DBToken)
			(*newV) = append((*newV), fR.expectedDBToken)
		}
	}
	return fR.expectedRetrieveError
}

func (fR *FakeRetriever) WriteInDB(query string, args ...interface{}) (int64, error) {
	return int64(0), nil
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
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty Repository with the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(types.Repository{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns the valid Repository struct", func() {
			It("Should return the expected Repository with a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedRepository: types.Repository{
						URL:       "teste",
						Branch:    "teste",
						CreatedAt: time.Now(),
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBRepository(
					map[string]interface{}{"repositoryURL": "teste"})
				Expect(repo).To(Equal(fakeRetriever.expectedRepository))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("FindOneDBSecurityTest", func() {
		Context("When key map verification returns false", func() {
			It("Should return an empty SecurityTest with the expected error", func() {
				postgres := PostgresRequests{}
				repo, err := postgres.FindOneDBSecurityTest(
					map[string]interface{}{"Wrong Key": "Failed"})
				Expect(repo).To(Equal(types.SecurityTest{}))
				Expect(err).To(Equal(errors.New("Could not find securityTest name field")))
			})
		})
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty SecurityTest with the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBSecurityTest(
					map[string]interface{}{"name": "teste"})
				Expect(repo).To(Equal(types.SecurityTest{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns the valid SecurityTest struct", func() {
			It("Should return the expected SecurityTest with a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedSecurityTest: types.SecurityTest{
						Name:     "teste",
						Image:    "teste",
						ImageTag: "teste",
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBSecurityTest(
					map[string]interface{}{"name": "teste"})
				Expect(repo).To(Equal(fakeRetriever.expectedSecurityTest))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("FindOneDBAnalysis", func() {
		Context("When key map verification returns false", func() {
			It("Should return an empty Analysis with the expected error", func() {
				postgres := PostgresRequests{}
				repo, err := postgres.FindOneDBAnalysis(
					map[string]interface{}{"Wrong Key": "Failed"})
				Expect(repo).To(Equal(types.Analysis{}))
				Expect(err).To(Equal(errors.New("Could not find RID field")))
			})
		})
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty Analysis with the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBAnalysis(
					map[string]interface{}{"RID": "teste"})
				Expect(repo).To(Equal(types.Analysis{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns the valid Analysis struct", func() {
			It("Should return the expected Analysis with a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedAnalysis: types.Analysis{
						RID:    "teste",
						URL:    "teste",
						Branch: "teste",
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBAnalysis(
					map[string]interface{}{"RID": "teste"})
				Expect(repo).To(Equal(fakeRetriever.expectedAnalysis))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("FindOneDBUser", func() {
		Context("When key map verification returns false", func() {
			It("Should return an empty User with the expected error", func() {
				postgres := PostgresRequests{}
				repo, err := postgres.FindOneDBUser(
					map[string]interface{}{"Wrong Key": "Failed"})
				Expect(repo).To(Equal(types.User{}))
				Expect(err).To(Equal(errors.New("Could not find user in DB")))
			})
		})
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty User with the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBUser(
					map[string]interface{}{"username": "teste"})
				Expect(repo).To(Equal(types.User{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns the valid User struct", func() {
			It("Should return the expected User with a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedUser: types.User{
						Username:   "teste",
						Password:   "teste",
						Salt:       "teste",
						Iterations: 1,
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBUser(
					map[string]interface{}{"username": "teste"})
				Expect(repo).To(Equal(fakeRetriever.expectedUser))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("FindOneDBAccessToken", func() {
		Context("When key map verification returns false", func() {
			It("Should return an empty DBToken with the expected error", func() {
				postgres := PostgresRequests{}
				repo, err := postgres.FindOneDBAccessToken(
					map[string]interface{}{"Wrong Key": "Failed"})
				Expect(repo).To(Equal(types.DBToken{}))
				Expect(err).To(Equal(errors.New("Could not find uuid parameter")))
			})
		})
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty DBToken with the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBAccessToken(
					map[string]interface{}{"uuid": "teste"})
				Expect(repo).To(Equal(types.DBToken{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns the valid DBToken struct", func() {
			It("Should return the expected DBToken with a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedDBToken: types.DBToken{
						HuskyToken: "teste",
						URL:        "teste",
						IsValid:    true,
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindOneDBAccessToken(
					map[string]interface{}{"uuid": "teste"})
				Expect(repo).To(Equal(fakeRetriever.expectedDBToken))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("ConfigureQuery", func() {
		Context("When a simple query is passed without parameters", func() {
			It("Should return the same query passed and an empty arguments slice", func() {
				expectedQuery := `SELECT * FROM test`
				params := map[string]interface{}{}
				expectedVals := make([]interface{}, 0)
				query, vals := ConfigureQuery(expectedQuery, params)
				Expect(query).To(Equal(expectedQuery))
				Expect(vals).To(Equal(expectedVals))
			})
		})
		Context("When a query is passed with parameters", func() {
			It("Should return the expected final query with the listed arguments", func() {
				params := map[string]interface{}{"teste1": "myTest", "teste2": 1}
				query, vals := ConfigureQuery(`SELECT * FROM test`, params)
				if _, ok := vals[0].(string); ok {
					expectedQuery := `SELECT * FROM test WHERE teste1 = $1 AND teste2 = $2`
					expectedVals := []interface{}{"myTest", 1}
					Expect(query).To(Equal(expectedQuery))
					Expect(vals).To(Equal(expectedVals))
				} else {
					expectedQuery := `SELECT * FROM test WHERE teste2 = $1 AND teste1 = $2`
					expectedVals := []interface{}{1, "myTest"}
					Expect(query).To(Equal(expectedQuery))
					Expect(vals).To(Equal(expectedVals))
				}
			})
		})
		Context("When a query is passed with only one parameter", func() {
			It("Should return the expected final query with just one argument", func() {
				expectedQuery := `SELECT * FROM test WHERE teste1 = $1`
				expectedVals := []interface{}{"myTest"}
				params := map[string]interface{}{"teste1": "myTest"}
				query, vals := ConfigureQuery(`SELECT * FROM test`, params)
				Expect(query).To(Equal(expectedQuery))
				Expect(vals).To(Equal(expectedVals))
			})
		})
	})
	Describe("FindAllDBRepository", func() {
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty array and the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindAllDBRepository(
					map[string]interface{}{"teste": "teste"})
				Expect(repo).To(Equal([]types.Repository{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns a nil error", func() {
			It("Should return the expected RepositoryResponse and a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedRepository: types.Repository{
						URL:       "teste",
						Branch:    "teste",
						CreatedAt: time.Now(),
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				expectedRepositoryArray := []types.Repository{}
				expectedRepositoryArray = append(expectedRepositoryArray,
					fakeRetriever.expectedRepository)
				repos, err := postgres.FindAllDBRepository(
					map[string]interface{}{"teste": "teste"})
				Expect(repos).To(Equal(expectedRepositoryArray))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("FindAllDBSecurityTest", func() {
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty array and the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindAllDBSecurityTest(
					map[string]interface{}{"teste": "teste"})
				Expect(repo).To(Equal([]types.SecurityTest{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns a nil error", func() {
			It("Should return the expected SecurityTestResponse and a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedSecurityTest: types.SecurityTest{
						Name:     "teste",
						Image:    "teste",
						ImageTag: "teste",
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				expectedSecurityTestArray := []types.SecurityTest{}
				expectedSecurityTestArray = append(expectedSecurityTestArray,
					fakeRetriever.expectedSecurityTest)
				repos, err := postgres.FindAllDBSecurityTest(
					map[string]interface{}{"teste": "teste"})
				Expect(repos).To(Equal(expectedSecurityTestArray))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("FindAllDBAnalysis", func() {
		Context("When RetrieveFromDB returns an error", func() {
			It("Should return an empty array and the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: errors.New("Failed to retrieve data"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				repo, err := postgres.FindAllDBAnalysis(
					map[string]interface{}{"teste": "teste"})
				Expect(repo).To(Equal([]types.Analysis{}))
				Expect(err).To(Equal(fakeRetriever.expectedRetrieveError))
			})
		})
		Context("When RetrieveFromDB returns a nil error", func() {
			It("Should return the expected AnalysisResponse and a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedRetrieveError: nil,
					expectedAnalysis: types.Analysis{
						RID:    "teste",
						URL:    "teste",
						Branch: "teste",
					},
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				expectedAnalysisArray := []types.Analysis{}
				expectedAnalysisArray = append(expectedAnalysisArray,
					fakeRetriever.expectedAnalysis)
				repos, err := postgres.FindAllDBAnalysis(
					map[string]interface{}{"teste": "teste"})
				Expect(repos).To(Equal(expectedAnalysisArray))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("ConfigureInsertQuery", func() {
		Context("When an Insert query is passed with some params", func() {
			It("Should return the expected query with the params to be inserted", func() {
				query := `INSERT into test`
				params := map[string]interface{}{
					"test1": "value1",
					"test2": 3,
				}
				finalQuery, values := ConfigureInsertQuery(query, params)
				if _, ok := values[0].(string); ok {
					expectedQuery := `INSERT into test (test1, test2) VALUES ($1, $2)`
					expectedValues := []interface{}{"value1", 3}
					Expect(finalQuery).To(Equal(expectedQuery))
					Expect(values).To(Equal(expectedValues))
				} else {
					expectedQuery := `INSERT into test (test2, test1) VALUES ($1, $2)`
					expectedValues := []interface{}{3, "value1"}
					Expect(finalQuery).To(Equal(expectedQuery))
					Expect(values).To(Equal(expectedValues))
				}
			})
		})
		Context("When an Insert query is passed with one param", func() {
			It("Should return the expected query with the one param to be inserted", func() {
				query := `INSERT into test`
				params := map[string]interface{}{
					"test1": "value1",
				}
				expectedQuery := `INSERT into test (test1) VALUES ($1)`
				expectedValues := []interface{}{"value1"}
				finalQuery, values := ConfigureInsertQuery(query, params)
				Expect(finalQuery).To(Equal(expectedQuery))
				Expect(values).To(Equal(expectedValues))
			})
		})
	})
	Describe("ConfigureUpdateQuery", func() {
		Context("When an update query is passed with a search value and a new value", func() {
			It("Should return the expected query and with the params to be updated", func() {
				query := `UPDATE test`
				searchValues := map[string]interface{}{"id": 1}
				newValues := map[string]interface{}{"teste1": "newVal"}
				expectedQuery := `UPDATE test SET teste1 = $2 WHERE id = $1`
				expectedValues := []interface{}{1, "newVal"}
				finalQuery, values := ConfigureUpdateQuery(query, searchValues, newValues)
				Expect(finalQuery).To(Equal(expectedQuery))
				Expect(values).To(Equal(expectedValues))
			})
		})
		Context("When an update query is passed with a search value and a new value", func() {
			It("Should return the expected query and with the params to be updated", func() {
				query := `UPDATE test`
				searchValues := map[string]interface{}{"id": 1, "id2": 2}
				newValues := map[string]interface{}{"teste1": "newVal", "teste2": 3}
				finalQuery, values := ConfigureUpdateQuery(query, searchValues, newValues)
				if _, ok := values[2].(string); ok {
					if values[0].(int) == 1 {
						expectedQuery := `UPDATE test SET teste1 = $3, teste2 = $4 WHERE id = $1 AND id2 = $2`
						expectedValues := []interface{}{1, 2, "newVal", 3}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					} else {
						expectedQuery := `UPDATE test SET teste1 = $3, teste2 = $4 WHERE id2 = $1 AND id = $2`
						expectedValues := []interface{}{2, 1, "newVal", 3}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					}
				} else {
					if values[0].(int) == 1 {
						expectedQuery := `UPDATE test SET teste2 = $3, teste1 = $4 WHERE id = $1 AND id2 = $2`
						expectedValues := []interface{}{1, 2, 3, "newVal"}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					} else {
						expectedQuery := `UPDATE test SET teste2 = $3, teste1 = $4 WHERE id2 = $1 AND id = $2`
						expectedValues := []interface{}{2, 1, 3, "newVal"}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					}
				}
			})
		})
	})
})

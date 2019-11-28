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
	expectedWriteError    error
	expectedNumberRows    int64
	expectedConnectError  error
	expectedPqArray       interface{}
}

func (fR *FakeRetriever) Connect(
	address string,
	dbName string,
	username string,
	password string,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration) error {
	return fR.expectedConnectError
}

func (fR *FakeRetriever) RetrieveFromDB(
	query string, response interface{}, arrayColumns []string, params ...interface{}) error {
	if fR.expectedRetrieveError == nil {
		switch r := response.(type) {
		case *[]types.Repository:
			// newV := response.(*[]types.Repository)
			(*r) = append((*r), fR.expectedRepository)
		case *[]types.SecurityTest:
			// newV := response.(*[]types.SecurityTest)
			(*r) = append((*r), fR.expectedSecurityTest)
		case *[]types.Analysis:
			// newV := response.(*[]types.Analysis)
			(*r) = append((*r), fR.expectedAnalysis)
		case *[]types.User:
			// newV := response.(*[]types.User)
			(*r) = append((*r), fR.expectedUser)
		case *[]types.DBToken:
			// newV := response.(*[]types.DBToken)
			(*r) = append((*r), fR.expectedDBToken)
		}
	}
	return fR.expectedRetrieveError
}

func (fR *FakeRetriever) WriteInDB(query string, args ...interface{}) (int64, error) {
	return fR.expectedNumberRows, fR.expectedWriteError
}

func (fR *FakeRetriever) PqArray(values []string) interface{} {
	return fR.expectedPqArray
}

type FakeJson struct {
	expectedMarshalData  []byte
	expectedMarshalError error
}

func (fJ *FakeJson) Marshal(v interface{}) ([]byte, error) {
	return fJ.expectedMarshalData, fJ.expectedMarshalError
}

func (fJ *FakeJson) Unmarshal(data []byte, v interface{}) error {
	return nil
}

var _ = Describe("Postgres", func() {

	var (
		securityTest types.SecurityTest
		repository   types.Repository
		analysis     types.Analysis
		user         types.User
		accessToken  types.DBToken
		validParams  map[string]interface{}
		validUpdate  map[string]interface{}
	)

	BeforeEach(func() {
		securityTest = types.SecurityTest{
			Name:  "teste",
			Image: "teste",
			Cmd:   "teste",
		}
		repository = types.Repository{
			URL:       "teste",
			CreatedAt: time.Now(),
		}
		analysis = types.Analysis{
			RID:    "teste",
			URL:    "teste",
			Branch: "teste",
			Status: "teste",
		}
		user = types.User{
			Username: "teste",
			Password: "teste",
			Salt:     "teste",
		}
		accessToken = types.DBToken{
			HuskyToken: "teste",
			URL:        "teste",
			IsValid:    true,
		}
		validParams = map[string]interface{}{"id": "teste"}
		validUpdate = map[string]interface{}{"changeField": "changeValue",
			"containers": []types.Container{}}
	})

	Describe("ConnectDB", func() {
		Context("When Connect returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedConnectError: errors.New("Failed to connect to DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.ConnectDB("", "", "", "", time.Second, 0, 0, 0, 0, time.Second)).To(
					Equal(fakeRetriever.expectedConnectError))
			})
		})
	})

	Describe("FindOneDBRepository", func() {
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
					expectedQuery := `SELECT * FROM test WHERE "teste1" = $1 AND "teste2" = $2`
					expectedVals := []interface{}{"myTest", 1}
					Expect(query).To(Equal(expectedQuery))
					Expect(vals).To(Equal(expectedVals))
				} else {
					expectedQuery := `SELECT * FROM test WHERE "teste2" = $1 AND "teste1" = $2`
					expectedVals := []interface{}{1, "myTest"}
					Expect(query).To(Equal(expectedQuery))
					Expect(vals).To(Equal(expectedVals))
				}
			})
		})
		Context("When a query is passed with only one parameter", func() {
			It("Should return the expected final query with just one argument", func() {
				expectedQuery := `SELECT * FROM test WHERE "teste1" = $1`
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
	Describe("InsertDBRepository", func() {
		Context("When Repository is set with a nil URL", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				expectedError := errors.New("Empty repository data")
				Expect(postgres.InsertDBRepository(types.Repository{})).To(Equal(expectedError))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write data in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(
					postgres.InsertDBRepository(repository)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns zero rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(
					postgres.InsertDBRepository(repository)).To(
					Equal(errors.New("No data was inserted")))
			})
		})
		Context("When WriteInDB returns a number of rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(
					postgres.InsertDBRepository(repository)).To(
					BeNil())
			})
		})
	})
	Describe("InsertDBSecurityTest", func() {
		Context("When an empty SecurityTest is passed as an argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				Expect(
					postgres.InsertDBSecurityTest(types.SecurityTest{})).To(
					Equal(errors.New("Empty SecurityTest data")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write data in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(
					postgres.InsertDBSecurityTest(securityTest)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.InsertDBSecurityTest(securityTest)).To(
					Equal(errors.New("No data was inserted")))
			})
		})
		Context("When WriteInDB returns some rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.InsertDBSecurityTest(securityTest)).To(BeNil())
			})
		})
	})
	Describe("InsertDBAnalysis", func() {
		Context("When an empty Analysis is passed as an argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				Expect(
					postgres.InsertDBAnalysis(types.Analysis{})).To(
					Equal(errors.New("Empty Analysis data")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write data in DB"),
				}
				fakeJson := FakeJson{
					expectedMarshalData: nil,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(
					postgres.InsertDBAnalysis(analysis)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				fakeJson := FakeJson{
					expectedMarshalData: nil,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.InsertDBAnalysis(analysis)).To(
					Equal(errors.New("No data was inserted")))
			})
		})
		Context("When WriteInDB returns some rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				fakeJson := FakeJson{
					expectedMarshalData: nil,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.InsertDBAnalysis(analysis)).To(BeNil())
			})
		})
	})

	Describe("InsertDBUser", func() {
		Context("When an empty User is passed as an argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				Expect(
					postgres.InsertDBUser(types.User{})).To(
					Equal(errors.New("Empty User data")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write data in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(
					postgres.InsertDBUser(user)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.InsertDBUser(user)).To(
					Equal(errors.New("No data was inserted")))
			})
		})
		Context("When WriteInDB returns some rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.InsertDBUser(user)).To(BeNil())
			})
		})
	})

	Describe("InsertDBAccessToken", func() {
		Context("When an empty DBToken is passed as an argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				Expect(
					postgres.InsertDBAccessToken(types.DBToken{})).To(
					Equal(errors.New("Empty DBToken data")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write data in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(
					postgres.InsertDBAccessToken(accessToken)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.InsertDBAccessToken(accessToken)).To(
					Equal(errors.New("No data was inserted")))
			})
		})
		Context("When WriteInDB returns some rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.InsertDBAccessToken(accessToken)).To(BeNil())
			})
		})
	})

	Describe("UpdateOneDBRepository", func() {
		Context("When an empty updateQuery is passed as argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateQuery := map[string]interface{}{}
				Expect(postgres.UpdateOneDBRepository(mapParams, updateQuery)).To(
					Equal(errors.New("Empty fields to be updated")))
			})
		})
		Context("When an empty mapParams is passed as argument", func() {
			It("Should return the expected error for empty mapParams", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateQuery := map[string]interface{}{"teste": "update"}
				Expect(postgres.UpdateOneDBRepository(mapParams, updateQuery)).To(
					Equal(errors.New("Empty fields to search")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBRepository(validParams, validUpdate)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBRepository(validParams, validUpdate)).To(
					Equal(errors.New("No data was updated")))
			})
		})
		Context("When WriteInDB returns a number of rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBRepository(validParams, validUpdate)).To(
					BeNil())
			})
		})
	})

	Describe("UpdateOneDBAnalysis", func() {
		Context("When an empty updateQuery is passed as argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateQuery := map[string]interface{}{}
				Expect(postgres.UpdateOneDBAnalysis(mapParams, updateQuery)).To(
					Equal(errors.New("Empty fields to be updated")))
			})
		})
		Context("When an empty mapParams is passed as argument", func() {
			It("Should return the expected error for empty mapParams", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateQuery := map[string]interface{}{"teste": "update"}
				Expect(postgres.UpdateOneDBAnalysis(mapParams, updateQuery)).To(
					Equal(errors.New("Empty fields to search")))
			})
		})
		Context("When ConfigureAnalysisData returns an error", func() {
			It("Should return the same error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: errors.New("Failed to Marshal data"),
				}
				postgres := PostgresRequests{
					JSONHandler: &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysis(validParams, validUpdate)).To(
					Equal(fakeJson.expectedMarshalError))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: nil,
				}
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysis(validParams, validUpdate)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: nil,
				}
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysis(validParams, validUpdate)).To(
					Equal(errors.New("No data was updated")))
			})
		})
		Context("When WriteInDB returns a number of rows affected", func() {
			It("Should return a nil error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: nil,
				}
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysis(validParams, validUpdate)).To(
					BeNil())
			})
		})
	})

	Describe("UpdateOneDBUser", func() {
		Context("When an empty updateUser is passed as argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateUser := types.User{}
				Expect(postgres.UpdateOneDBUser(mapParams, updateUser)).To(
					Equal(errors.New("Empty fields to be updated")))
			})
		})
		Context("When an empty mapParams is passed as argument", func() {
			It("Should return the expected error for empty mapParams", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				Expect(postgres.UpdateOneDBUser(mapParams, user)).To(
					Equal(errors.New("Empty fields to search")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBUser(validParams, user)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBUser(validParams, user)).To(
					Equal(errors.New("No data was updated")))
			})
		})
		Context("When WriteInDB returns a number of rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBUser(validParams, user)).To(
					BeNil())
			})
		})
	})

	Describe("UpdateOneDBAnalysisContainer", func() {
		Context("When an empty updateQuery is passed as argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateQuery := map[string]interface{}{}
				Expect(postgres.UpdateOneDBAnalysisContainer(mapParams, updateQuery)).To(
					Equal(errors.New("Empty fields to be updated")))
			})
		})
		Context("When an empty mapParams is passed as argument", func() {
			It("Should return the expected error for empty mapParams", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updateQuery := map[string]interface{}{"teste": "update"}
				Expect(postgres.UpdateOneDBAnalysisContainer(mapParams, updateQuery)).To(
					Equal(errors.New("Empty fields to search")))
			})
		})
		Context("When ConfigureAnalysisData returns an error", func() {
			It("Should return the same error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: errors.New("Failed trying to Marshal data"),
				}
				postgres := PostgresRequests{
					JSONHandler: &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysisContainer(validParams, validUpdate)).To(
					Equal(fakeJson.expectedMarshalError))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: nil,
				}
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysisContainer(validParams, validUpdate)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: nil,
				}
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysisContainer(validParams, validUpdate)).To(
					Equal(errors.New("No data was updated")))
			})
		})
		Context("When WriteInDB returns a number of rows affected", func() {
			It("Should return a nil error", func() {
				fakeJson := FakeJson{
					expectedMarshalError: nil,
				}
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
					JSONHandler:   &fakeJson,
				}
				Expect(postgres.UpdateOneDBAnalysisContainer(validParams, validUpdate)).To(
					BeNil())
			})
		})
	})

	Describe("UpdateOneDBAccessToken", func() {
		Context("When an empty updatedAccessToken is passed as argument", func() {
			It("Should return the expected error", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				updatedAccessToken := types.DBToken{}
				Expect(postgres.UpdateOneDBAccessToken(mapParams, updatedAccessToken)).To(
					Equal(errors.New("Empty fields to be updated")))
			})
		})
		Context("When an empty mapParams is passed as argument", func() {
			It("Should return the expected error for empty mapParams", func() {
				postgres := PostgresRequests{}
				mapParams := map[string]interface{}{}
				Expect(postgres.UpdateOneDBAccessToken(mapParams, accessToken)).To(
					Equal(errors.New("Empty fields to search")))
			})
		})
		Context("When WriteInDB returns an error", func() {
			It("Should return the same error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBAccessToken(validParams, accessToken)).To(
					Equal(fakeRetriever.expectedWriteError))
			})
		})
		Context("When WriteInDB returns 0 rows affected", func() {
			It("Should return the expected error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBAccessToken(validParams, accessToken)).To(
					Equal(errors.New("No data was updated")))
			})
		})
		Context("When WriteInDB returns a number of rows affected", func() {
			It("Should return a nil error", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				Expect(postgres.UpdateOneDBAccessToken(validParams, accessToken)).To(
					BeNil())
			})
		})
	})
	Describe("UpsertOneDBSecurityTest", func() {
		Context("When an empty SecurityTest is passed", func() {
			It("Should return the expected error and a nil interface", func() {
				params := map[string]interface{}{}
				postgres := PostgresRequests{}
				status, err := postgres.UpsertOneDBSecurityTest(params, types.SecurityTest{})
				Expect(err).To(Equal(errors.New("Empty fields to be updated")))
				Expect(status).To(BeNil())
			})
		})
		Context("When an empty mapParams is passed", func() {
			It("Should return the expected error", func() {
				params := map[string]interface{}{}
				postgres := PostgresRequests{}
				status, err := postgres.UpsertOneDBSecurityTest(params, securityTest)
				Expect(err).To(Equal(errors.New("Empty fields to search")))
				Expect(status).To(BeNil())
			})
		})
		Context("When DataRetriever returns an error", func() {
			It("Should return the same error with a nil interface", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: errors.New("Failed to write in DB"),
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				status, err := postgres.UpsertOneDBSecurityTest(validParams, securityTest)
				Expect(err).To(Equal(fakeRetriever.expectedWriteError))
				Expect(status).To(BeNil())
			})
		})
		Context("When DataRetriever returns 0 rows affected", func() {
			It("Should return the expected error with a nil interface", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 0,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				status, err := postgres.UpsertOneDBSecurityTest(validParams, securityTest)
				Expect(err).To(Equal(errors.New("No data was updated")))
				Expect(status).To(BeNil())
			})
		})
		Context("When DataRetriever returns a valid number of rows affected", func() {
			It("Should return a nil error and the number of rows affected", func() {
				fakeRetriever := FakeRetriever{
					expectedWriteError: nil,
					expectedNumberRows: 1,
				}
				postgres := PostgresRequests{
					DataRetriever: &fakeRetriever,
				}
				status, err := postgres.UpsertOneDBSecurityTest(validParams, securityTest)
				Expect(err).To(BeNil())
				Expect(status).To(Equal(int64(1)))
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
					expectedQuery := `INSERT into test ("test1", "test2") VALUES ($1, $2)`
					expectedValues := []interface{}{"value1", 3}
					Expect(finalQuery).To(Equal(expectedQuery))
					Expect(values).To(Equal(expectedValues))
				} else {
					expectedQuery := `INSERT into test ("test2", "test1") VALUES ($1, $2)`
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
				expectedQuery := `INSERT into test ("test1") VALUES ($1)`
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
				expectedQuery := `UPDATE test SET "teste1" = $2 WHERE "id" = $1`
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
						expectedQuery := `UPDATE test SET "teste1" = $3, "teste2" = $4 WHERE "id" = $1 AND "id2" = $2`
						expectedValues := []interface{}{1, 2, "newVal", 3}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					} else {
						expectedQuery := `UPDATE test SET "teste1" = $3, "teste2" = $4 WHERE "id2" = $1 AND "id" = $2`
						expectedValues := []interface{}{2, 1, "newVal", 3}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					}
				} else {
					if values[0].(int) == 1 {
						expectedQuery := `UPDATE test SET "teste2" = $3, "teste1" = $4 WHERE "id" = $1 AND "id2" = $2`
						expectedValues := []interface{}{1, 2, 3, "newVal"}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					} else {
						expectedQuery := `UPDATE test SET "teste2" = $3, "teste1" = $4 WHERE "id2" = $1 AND "id" = $2`
						expectedValues := []interface{}{2, 1, 3, "newVal"}
						Expect(finalQuery).To(Equal(expectedQuery))
						Expect(values).To(Equal(expectedValues))
					}
				}
			})
		})
	})
	Describe("ConfigureUpsertQuery", func() {
		Context("When a Insert query is passed with one existed param and one new value", func() {
			It("Should return the expected query with the expected values", func() {
				query := `INSERT into test`
				newValues := map[string]interface{}{"value1": "teste1"}
				columnsConflict := map[string]interface{}{"unique": "existed"}
				expectedQuery := `INSERT into test ("value1") VALUES ($1) ON CONFLICT ("unique") DO UPDATE SET "value1" = EXCLUDED."value1"`
				expectedValues := []interface{}{"teste1"}
				finalQuery, values := ConfigureUpsertQuery(query, columnsConflict, newValues)
				Expect(finalQuery).To(Equal(expectedQuery))
				Expect(values).To(Equal(expectedValues))
			})
		})
	})
	Describe("ConfigureAnalysisData", func() {
		Context("When updatedAnalysis has containers and Marshal function returns an error", func() {
			It("Should return the same error and the map variable passed in the argument", func() {
				myCommits := []string{"author1", "author2"}
				updatedAnalysis := map[string]interface{}{
					"commitAuthors": myCommits,
					"containers":    []types.Container{},
				}
				fakeJson := FakeJson{
					expectedMarshalError: errors.New("Failed trying to Marshal data"),
				}
				fakeRetriever := FakeRetriever{
					expectedPqArray: myCommits,
				}
				postgres := PostgresRequests{
					JSONHandler:   &fakeJson,
					DataRetriever: &fakeRetriever,
				}
				newUpdatedData, err := postgres.ConfigureAnalysisData(updatedAnalysis)
				Expect(err).To(Equal(fakeJson.expectedMarshalError))
				Expect(newUpdatedData).To(Equal(updatedAnalysis))
			})
		})
		Context("When updatedAnalysis has huskyciresults and Marshal function returns an error ", func() {
			It("Should return the same error and the map variable passed in the argument", func() {
				myCommits := []string{"author1", "author2"}
				updatedAnalysis := map[string]interface{}{
					"commitAuthors":  myCommits,
					"huskyciresults": types.HuskyCIResults{},
				}
				fakeJson := FakeJson{
					expectedMarshalError: errors.New("Failed trying to Marshal data"),
				}
				fakeRetriever := FakeRetriever{
					expectedPqArray: myCommits,
				}
				postgres := PostgresRequests{
					JSONHandler:   &fakeJson,
					DataRetriever: &fakeRetriever,
				}
				newUpdatedData, err := postgres.ConfigureAnalysisData(updatedAnalysis)
				Expect(err).To(Equal(fakeJson.expectedMarshalError))
				Expect(newUpdatedData).To(Equal(updatedAnalysis))
			})
		})
		Context("When updatedAnalysis has codes and Marshal function returns an error ", func() {
			It("Should return the same error and the map variable passed in the argument", func() {
				myCommits := []string{"author1", "author2"}
				updatedAnalysis := map[string]interface{}{
					"commitAuthors": myCommits,
					"codes":         []types.Code{},
				}
				fakeJson := FakeJson{
					expectedMarshalError: errors.New("Failed trying to Marshal data"),
				}
				fakeRetriever := FakeRetriever{
					expectedPqArray: myCommits,
				}
				postgres := PostgresRequests{
					JSONHandler:   &fakeJson,
					DataRetriever: &fakeRetriever,
				}
				newUpdatedData, err := postgres.ConfigureAnalysisData(updatedAnalysis)
				Expect(err).To(Equal(fakeJson.expectedMarshalError))
				Expect(newUpdatedData).To(Equal(updatedAnalysis))
			})
		})
		Context("When updatedAnalysis has all expected data with valid JSON returned in Marshal", func() {
			It("Should return an nil error and the map variable passed in the argument", func() {
				myCommits := []string{"author1", "author2"}
				marshalData := []byte("ok")
				updatedAnalysis := map[string]interface{}{
					"commitAuthors":  myCommits,
					"codes":          []types.Code{},
					"containers":     []types.Container{},
					"huskyciresults": types.HuskyCIResults{},
				}
				expectedAnalysis := map[string]interface{}{
					"commitAuthors":  myCommits,
					"codes":          marshalData,
					"containers":     marshalData,
					"huskyciresults": marshalData,
				}
				fakeJson := FakeJson{
					expectedMarshalError: nil,
					expectedMarshalData:  marshalData,
				}
				fakeRetriever := FakeRetriever{
					expectedPqArray: myCommits,
				}
				postgres := PostgresRequests{
					JSONHandler:   &fakeJson,
					DataRetriever: &fakeRetriever,
				}
				newUpdatedData, err := postgres.ConfigureAnalysisData(updatedAnalysis)
				Expect(err).To(BeNil())
				Expect(newUpdatedData).To(Equal(expectedAnalysis))
			})
		})
	})
})

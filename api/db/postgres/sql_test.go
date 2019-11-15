package db_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"

	. "github.com/globocom/huskyCI/api/db/postgres"
)

type FakeSql struct {
	ExpectedConfigDbError   error
	ExpectedCloseDbError    error
	ExpectedConfigError     error
	ExpectedGetRowsError    error
	ExpectedNumRows         int64
	ExpectedExecError       error
	ExpectedColumns         []string
	ExpectedGetColumnsError error
	ExpectedCloseRowsError  error
	ExpectedScanRowError    error
	NumOfRows               int
	ActualRow               int
	ExpectedRecError        error
	ExpectedResult          []map[string]interface{}
}

func (fakeDB *FakeSql) ConfigureDB(address, username, password, dbName string) error {
	return fakeDB.ExpectedConfigDbError
}

func (fakeDB *FakeSql) ConfigurePool(maxOpenConns, maxIdleConns int, connLT time.Duration) {}

func (fakeDB *FakeSql) CloseDB() error {
	return fakeDB.ExpectedCloseDbError
}

func (fakeDB *FakeSql) ConfigureQuery(query string, args ...interface{}) error {
	return fakeDB.ExpectedConfigError
}

func (fakeDB *FakeSql) CloseRows() error {
	return fakeDB.ExpectedCloseRowsError
}

func (fakeDB *FakeSql) GetColumns() ([]string, error) {
	return fakeDB.ExpectedColumns, fakeDB.ExpectedGetColumnsError
}

func (fakeDB *FakeSql) HasNextRow() bool {
	if fakeDB.ActualRow == fakeDB.NumOfRows {
		return false
	}
	fakeDB.ActualRow += 1
	return true
}

func (fakeDB *FakeSql) ScanRow(dest ...interface{}) error {
	for i, val := range dest {
		val = "teste"
		dest[i] = &val
	}
	return fakeDB.ExpectedScanRowError
}

func (fakeDB *FakeSql) GetRowsErr() error {
	return fakeDB.ExpectedRecError
}

func (fakeDB *FakeSql) Exec(query string, args ...interface{}) error {
	return fakeDB.ExpectedExecError
}

func (fakeDB *FakeSql) GetRowsAffected() (int64, error) {
	return fakeDB.ExpectedNumRows, fakeDB.ExpectedGetRowsError
}

var _ = Describe("Sql", func() {
	Describe("Connect", func() {
		Context("When ConfigureDB returns an error", func() {
			It("Should return the expected error", func() {
				fakeDB := FakeSql{
					ExpectedConfigDbError: errors.New("The configuration has failed"),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				Expect(sqlConfig.Connect("test", "test", "test", "test", 1, 1, time.Hour)).To(Equal(fakeDB.ExpectedConfigDbError))
			})
		})
		Context("When ConfigureDB returns a nil error", func() {
			It("Should return a nil error", func() {
				fakeDB := FakeSql{
					ExpectedConfigDbError: nil,
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				Expect(sqlConfig.Connect("test", "test", "test", "test", 1, 1, time.Hour)).To(BeNil())
			})
		})
	})
	Describe("CloseDB", func() {
		Context("When CloseDB returns an error", func() {
			It("Should return the same error", func() {
				closeErrors := []error{
					errors.New("Failed during closing DB"),
					nil,
				}
				for _, closeError := range closeErrors {
					fakeDB := FakeSql{
						ExpectedCloseDbError: closeError,
					}
					sqlConfig := SQLConfig{
						Postgres: &fakeDB,
					}
					if closeError != nil {
						Expect(sqlConfig.CloseDB()).To(Equal(closeError))
					} else {
						Expect(sqlConfig.CloseDB()).To(BeNil())
					}
				}
			})
		})
	})
	Describe("GetValuesFromDB", func() {
		Context("When ConfigureQuery returns an error", func() {
			It("Should return the same error and a nil map", func() {
				fakeDB := FakeSql{
					ExpectedConfigError: errors.New("Error during get values from DB"),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				val, err := sqlConfig.GetValuesFromDB("blabla", "arg1", "arg2")
				Expect(val).To(BeNil())
				Expect(err).To(Equal(fakeDB.ExpectedConfigError))
			})
		})
		Context("When GetColumns returns an error", func() {
			It("Should return the same error and a nil map value", func() {
				fakeDB := FakeSql{
					ExpectedConfigError:     nil,
					ExpectedCloseRowsError:  nil,
					ExpectedGetColumnsError: errors.New("Failed trying to get columns"),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				myRows, err := sqlConfig.GetValuesFromDB("blabla", "somearg")
				Expect(myRows).To(BeNil())
				Expect(err).To(Equal(fakeDB.ExpectedGetColumnsError))
			})
		})
		Context("When ScanRow returns an error", func() {
			It("Should return the same error and a nil map value", func() {
				fakeDB := FakeSql{
					ExpectedConfigError:     nil,
					ExpectedCloseRowsError:  nil,
					ExpectedGetColumnsError: nil,
					NumOfRows:               3,
					ExpectedScanRowError:    errors.New("Failed during row scanning"),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				myRows, err := sqlConfig.GetValuesFromDB("blabla", "somearg")
				Expect(myRows).To(BeNil())
				Expect(err).To(Equal(fakeDB.ExpectedScanRowError))
			})
		})
		Context("When GetRowsErr returns an error", func() {
			It("Should return the same error and a nil map value", func() {
				fakeDB := FakeSql{
					ExpectedConfigError:     nil,
					ExpectedCloseRowsError:  nil,
					ExpectedGetColumnsError: nil,
					ExpectedScanRowError:    nil,
					NumOfRows:               3,
					ExpectedRecError:        errors.New("Error trying to proccess a row"),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				myRows, err := sqlConfig.GetValuesFromDB("blabla", "somearg")
				Expect(fakeDB.ActualRow).To(Equal(fakeDB.NumOfRows))
				Expect(myRows).To(BeNil())
				Expect(err).To(Equal(fakeDB.ExpectedRecError))
			})
		})
		Context("When results are empty", func() {
			It("Should return the expected error and a nil map value", func() {
				fakeDB := FakeSql{
					ExpectedConfigError:     nil,
					ExpectedCloseRowsError:  nil,
					ExpectedGetColumnsError: nil,
					ExpectedScanRowError:    nil,
					ExpectedRecError:        nil,
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				myRows, err := sqlConfig.GetValuesFromDB("blabla", "somearg")
				Expect(fakeDB.ActualRow).To(Equal(fakeDB.NumOfRows))
				Expect(myRows).To(BeNil())
				Expect(err).To(Equal(errors.New("No data found")))
			})
		})
		Context("When results are not empty", func() {
			It("Should return a nil error and the expected map", func() {
				fakeDB := FakeSql{
					ExpectedConfigError:     nil,
					ExpectedCloseRowsError:  nil,
					ExpectedGetColumnsError: nil,
					ExpectedScanRowError:    nil,
					ExpectedRecError:        nil,
					NumOfRows:               3,
					ExpectedColumns:         []string{"column1", "column2"},
					ExpectedResult: []map[string]interface{}{
						map[string]interface{}{"column1": "teste", "column2": "teste"},
						map[string]interface{}{"column1": "teste", "column2": "teste"},
						map[string]interface{}{"column1": "teste", "column2": "teste"},
					},
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				myRows, err := sqlConfig.GetValuesFromDB("blabla", "somearg")
				Expect(fakeDB.ActualRow).To(Equal(fakeDB.NumOfRows))
				Expect(myRows).To(Equal(fakeDB.ExpectedResult))
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("WriteInDB", func() {
		Context("When Exec returns an error", func() {
			It("Should return the same error with 0 value", func() {
				fakeDB := FakeSql{
					ExpectedExecError: errors.New("DB Failed"),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				val, err := sqlConfig.WriteInDB("bla", "myArgs")
				Expect(err).To(Equal(fakeDB.ExpectedExecError))
				Expect(val).To(Equal(int64(0)))
			})
		})
		Context("When Exec returns a nil error", func() {
			It("Should return the expected number of rows and nil error", func() {
				fakeDB := FakeSql{
					ExpectedExecError:    nil,
					ExpectedGetRowsError: nil,
					ExpectedNumRows:      int64(1),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				val, err := sqlConfig.WriteInDB("bla", "myArgs")
				Expect(err).To(BeNil())
				Expect(val).To(Equal(fakeDB.ExpectedNumRows))
			})
			It("Should return the expected GetRows error and the number of rows", func() {
				fakeDB := FakeSql{
					ExpectedExecError:    nil,
					ExpectedGetRowsError: errors.New("Error trying to get rows"),
					ExpectedNumRows:      int64(0),
				}
				sqlConfig := SQLConfig{
					Postgres: &fakeDB,
				}
				val, err := sqlConfig.WriteInDB("bla", "myArgs")
				Expect(err).To(Equal(fakeDB.ExpectedGetRowsError))
				Expect(val).To(Equal(fakeDB.ExpectedNumRows))
			})
		})
	})
})

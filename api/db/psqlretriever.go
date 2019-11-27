package db

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lib/pq"
)

// ConvertStringToSlice parses the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func ConvertStringToSlice(array string) []string {
	unquotedChar := `[^",\\{}\s(NULL)]`
	unquotedValue := fmt.Sprintf("(%s)+", unquotedChar)
	quotedChar := `[^"\\]|\\"|\\\\`
	quotedValue := fmt.Sprintf("\"(%s)*\"", quotedChar)
	arrayValue := fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)
	arrayExp := regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))
	var valueIndex int
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	fmt.Println("The matches:", matches)
	for _, match := range matches {
		s := match[valueIndex]
		fmt.Println("Values:", s)
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		// trim the comma if it is more than on element:
		s = strings.Trim(s, ",")
		results = append(results, s)
	}
	return results
}

// Connect will establish and a connection with Postgres DB based on the received arguments
func (sR *SQLJSONRetrieve) Connect(
	address string,
	dbName string,
	username string,
	password string,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration) error {
	return sR.Psql.Connect(
		address,
		username,
		password,
		dbName,
		maxOpenConns,
		maxIdleConns,
		connMaxLifetime)
}

// RetrieveFromDB will get values from Postgres DB through a SELECT query. The response will
// be returned in response interface. It needs to be a pointer to the struct that will map the
// returned values. The params passed are the ones that will be listed in WHERE argument of the
// query if it exists.
func (sR *SQLJSONRetrieve) RetrieveFromDB(
	query string, response interface{}, arrayColumns []string, params ...interface{}) error {
	values, err := sR.Psql.GetValuesFromDB(query, params...)
	if err != nil {
		return err
	}
	// Verifying if a column with array strings exists. In this case,
	// should convert the entry to a slice type
	for _, column := range arrayColumns {
		if column != "" {
			for _, row := range values {
				convertedRow, ok := row[column].([]uint8)
				if !ok {
					continue
				}
				arrayValue := ConvertStringToSlice(string(convertedRow))
				row[column] = arrayValue
			}
		}
	}
	// Converting all other []uint8 entries to an interface{}.
	// If they are JSON objects, it should be mapped in the struct
	// inside response variable.
	for _, row := range values {
		for k, v := range row {
			if val, ok := v.([]uint8); ok {
				var newJSON interface{}
				if err = sR.JSONHandler.Unmarshal([]byte(string(val)), &newJSON); err != nil {
					continue
				}
				row[k] = newJSON
			}
		}
	}
	jsonValues, err := sR.JSONHandler.Marshal(values)
	if err != nil {
		return err
	}
	if err := sR.JSONHandler.Unmarshal(jsonValues, response); err != nil {
		return err
	}
	return nil
}

// WriteInDB will call Postgres with INSERT and UPDATE queries. The values to be inserted or
// updated are passed in args variable.
func (sR *SQLJSONRetrieve) WriteInDB(query string, args ...interface{}) (int64, error) {
	return sR.Psql.WriteInDB(query, args...)
}

// PqArray will get the array passed as argument and return
// an interface with the right format to store an array in
// Postgres.
func (sR *SQLJSONRetrieve) PqArray(values []string) interface{} {
	return pq.Array(values)
}

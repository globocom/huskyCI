package securitytest

import "errors"

// Store holds all available security tests
type Store struct {
	SecurityTests []SecurityTest
}

// SecurityTest is the struct that stores all data from the security tests to be executed.
type SecurityTest struct {
	Name             string `bson:"name" json:"name"`
	Image            string `bson:"image" json:"image"`
	Command          string `bson:"cmd" json:"cmd"`
	Type             string `bson:"type" json:"type"`
	Language         string `bson:"language" json:"language"`
	Default          bool   `bson:"default" json:"default"`
	TimeOutInSeconds int    `bson:"timeOutSeconds" json:"timeOutSeconds"`
}

// GetByName returns a security test scruct based given a name
func (s *Store) GetByName(securityTestName string) (SecurityTest, error) {
	for _, securityTest := range s.SecurityTests {
		if securityTest.Name == securityTestName {
			return securityTest, nil
		}
	}
	return SecurityTest{}, errors.New("security test not found")
}

// GetAllByLanguage returns all security tests scructs based on a given language
func (s *Store) GetAllByLanguage(securityTestLanguage string) ([]SecurityTest, error) {
	securityTestsFound := []SecurityTest{}
	for _, securityTest := range s.SecurityTests {
		if securityTest.Language == securityTestLanguage {
			securityTestsFound = append(securityTestsFound, securityTest)
		}
	}
	if len(securityTestsFound) > 0 {
		return securityTestsFound, nil
	}
	return []SecurityTest{}, errors.New("no security tests found for this language")
}

// GetAllByType returns all security tests scructs based on a given type
func (s *Store) GetAllByType(securityTestType string) ([]SecurityTest, error) {
	securityTestsFound := []SecurityTest{}
	for _, securityTest := range s.SecurityTests {
		if securityTest.Type == securityTestType {
			securityTestsFound = append(securityTestsFound, securityTest)
		}
	}
	if len(securityTestsFound) > 0 {
		return securityTestsFound, nil
	}
	return []SecurityTest{}, errors.New("no security tests found for this type")
}

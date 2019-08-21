package securitytest

import (
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

var securityTestFunctions = map[string]func(e EnryScan, a *AllScansResult) error{
	"huskyci/gosec":    initGoSec,
	"huskyci/bandit":   initBandit,
	"huskyci/safety":   initSafety,
	"huskyci/brakeman": initBrakeman,
	"huskyci/npmaudit": initNpmaudit,
}

// AllScansResult store all scans results of an Analysis
type AllScansResult struct {
	RID            string
	Status         string
	Containers     []types.Container
	Codes          []Code
	FinalResult    string
	HuskyCIResults types.HuskyCIResults
}

// RunAllScans runs both generic and language security
func RunAllScans(enryScan EnryScan) AllScansResult {

	allScansResult := AllScansResult{}
	allScansResult.Codes = enryScan.FinalOutput.Codes

	if err := runGenericScans(&allScansResult); err != nil {
		return allScansResult
	}

	if err := runLanguageScans(&allScansResult, enryScan); err != nil {
		return allScansResult
	}

	return allScansResult
}

func runGenericScans(allScansResult *AllScansResult) error {

	genericTests, err := getAllDefaultSecurityTests("Generic", "")
	if err != nil {
		return err
	}

	enryScan := EnryScan{}

	for _, genericTest := range genericTests {
		if genericTest.Name != "enry" {
			if err := initSecurityTest(genericTest, allScansResult, enryScan); err != nil {
				return err
			}
		}
	}

	return nil
}

func runLanguageScans(allScansResult *AllScansResult, enryScan EnryScan) error {

	languageTests := []types.SecurityTest{}

	for _, code := range enryScan.FinalOutput.Codes {
		codeTests, err := getAllDefaultSecurityTests("Language", code.Language)
		if err != nil {
			return err
		}
		languageTests = append(languageTests, codeTests...)
	}

	for _, languageTest := range languageTests {
		if err := initSecurityTest(languageTest, allScansResult, enryScan); err != nil {
			return err
		}
	}

	return nil
}

func getAllDefaultSecurityTests(typeOf, language string) ([]types.SecurityTest, error) {
	securityTests := []types.SecurityTest{}
	securityTestQuery := map[string]interface{}{"type": typeOf, "default": true}
	if language != "" {
		securityTestQuery = map[string]interface{}{"language": language, "default": true}
	}
	securityTests, err := db.FindAllDBSecurityTest(securityTestQuery)
	if err != nil {
		log.Error("getAllDefaultSecurityTests", "SECURITYTEST", 2009, err)
		return securityTests, err
	}
	return securityTests, nil
}

func initSecurityTest(securityTest types.SecurityTest, allScansResult *AllScansResult, enryScan EnryScan) error {
	securityTestFunction := securityTestFunctions[securityTest.Image]
	if err := securityTestFunction(enryScan, allScansResult); err != nil {
		return err
	}
	return nil
}

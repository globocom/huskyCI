package securitytest

import (
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
)

// RunAllInfo store all scans results of an Analysis
type RunAllInfo struct {
	RID            string
	Status         string
	Containers     []types.Container
	CommitAuthors  []string
	Codes          []Code
	FinalResult    string
	ErrorFound     error
	HuskyCIResults types.HuskyCIResults
}

// Start runs both generic and language security
func (results *RunAllInfo) Start(enryScan SecTestScanInfo) error {

	results.Codes = enryScan.Codes

	if err := results.runGenericScans(enryScan); err != nil {
		return err
	}

	if err := results.runLanguageScans(enryScan); err != nil {
		return err
	}

	results.setToAnalysis()

	return nil
}

func (results *RunAllInfo) runGenericScans(enryScan SecTestScanInfo) error {

	genericTests, err := getAllDefaultSecurityTests("Generic", "")
	if err != nil {
		return err
	}

	for _, genericTest := range genericTests {
		newGenericScan := SecTestScanInfo{}
		if err := newGenericScan.New(enryScan.RID, enryScan.URL, enryScan.Branch, genericTest.Name); err != nil {
			return err
		}
		if err := newGenericScan.Start(); err != nil {
			return err
		}
		results.Containers = append(results.Containers, newGenericScan.Container)
		if genericTest.Name == "gitauthors" {
			results.CommitAuthors = newGenericScan.CommitAuthors.Authors
		}
	}

	return nil
}

func (results *RunAllInfo) runLanguageScans(enryScan SecTestScanInfo) error {

	languageTests := []types.SecurityTest{}
	for _, code := range enryScan.Codes {
		codeTests, err := getAllDefaultSecurityTests("Language", code.Language)
		if err != nil {
			return err
		}
		languageTests = append(languageTests, codeTests...)
	}

	for _, languageTest := range languageTests {
		newLanguageScan := SecTestScanInfo{}
		if err := newLanguageScan.New(enryScan.RID, enryScan.URL, enryScan.Branch, languageTest.Name); err != nil {
			return err
		}
		if err := newLanguageScan.Start(); err != nil {
			return err
		}
		results.Containers = append(results.Containers, newLanguageScan.Container)
		results.setVulns(newLanguageScan)
	}

	return nil
}

func (results *RunAllInfo) setVulns(securityTestScan SecTestScanInfo) {

	for _, highVuln := range securityTestScan.Vulnerabilities.HighVulns {
		switch securityTestScan.SecurityTestName {
		case "bandit":
			results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.HighVulns = append(results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.HighVulns, highVuln)
		case "brakeman":
			results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.HighVulns = append(results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.HighVulns, highVuln)
		case "safety":
			results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.HighVulns = append(results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.HighVulns, highVuln)
		case "gosec":
			results.HuskyCIResults.GoResults.HuskyCIGosecOutput.HighVulns = append(results.HuskyCIResults.GoResults.HuskyCIGosecOutput.HighVulns, highVuln)
		case "npmaudit":
			results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns = append(results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.HighVulns, highVuln)
		}
	}

	for _, mediumVuln := range securityTestScan.Vulnerabilities.MediumVulns {
		switch securityTestScan.SecurityTestName {
		case "bandit":
			results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.MediumVulns = append(results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.MediumVulns, mediumVuln)
		case "brakeman":
			results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.MediumVulns = append(results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.MediumVulns, mediumVuln)
		case "safety":
			results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.MediumVulns = append(results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.MediumVulns, mediumVuln)
		case "gosec":
			results.HuskyCIResults.GoResults.HuskyCIGosecOutput.MediumVulns = append(results.HuskyCIResults.GoResults.HuskyCIGosecOutput.MediumVulns, mediumVuln)
		case "npmaudit":
			results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns = append(results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.MediumVulns, mediumVuln)
		}
	}

	for _, lowVuln := range securityTestScan.Vulnerabilities.LowVulns {
		switch securityTestScan.SecurityTestName {
		case "bandit":
			results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.LowVulns = append(results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.LowVulns, lowVuln)
		case "brakeman":
			results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.LowVulns = append(results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.LowVulns, lowVuln)
		case "safety":
			results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.LowVulns = append(results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.LowVulns, lowVuln)
		case "gosec":
			results.HuskyCIResults.GoResults.HuskyCIGosecOutput.LowVulns = append(results.HuskyCIResults.GoResults.HuskyCIGosecOutput.LowVulns, lowVuln)
		case "npmaudit":
			results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns = append(results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.LowVulns, lowVuln)
		}
	}

	for _, noSec := range securityTestScan.Vulnerabilities.NoSecVulns {
		switch securityTestScan.SecurityTestName {
		case "bandit":
			results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.NoSecVulns = append(results.HuskyCIResults.PythonResults.HuskyCIBanditOutput.NoSecVulns, noSec)
		case "brakeman":
			results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.NoSecVulns = append(results.HuskyCIResults.RubyResults.HuskyCIBrakemanOutput.NoSecVulns, noSec)
		case "safety":
			results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.NoSecVulns = append(results.HuskyCIResults.PythonResults.HuskyCISafetyOutput.NoSecVulns, noSec)
		case "gosec":
			results.HuskyCIResults.GoResults.HuskyCIGosecOutput.NoSecVulns = append(results.HuskyCIResults.GoResults.HuskyCIGosecOutput.LowVulns, noSec)
		case "npmaudit":
			results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.NoSecVulns = append(results.HuskyCIResults.JavaScriptResults.HuskyCINpmAuditOutput.NoSecVulns, noSec)
		}
	}
}

func (results *RunAllInfo) setToAnalysis() {

	results.Status = "finished"
	results.FinalResult = "success"

	if results.ErrorFound != nil {
		results.Status = "error running"
		results.FinalResult = "error"
		return
	}

	for _, container := range results.Containers {
		switch container.CResult {
		case "warning":
			results.FinalResult = "warning"
		case "failed":
			results.FinalResult = "failed"
			return
		}
	}
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

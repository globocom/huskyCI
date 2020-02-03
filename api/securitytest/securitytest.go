package securitytest

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/globocom/huskyCI/api/container"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/vulnerability"
)

// DefaultSecTestConf hold all information from all securityTests
var (
	DefaultSecTestConf *ViperConfig
	AllGeneric         []*SecurityTest
	GitleaksConfig     *SecurityTest
	BanditConfig       *SecurityTest
	BrakemanConfig     *SecurityTest
	GosecConfig        *SecurityTest
	NpmAuditConfig     *SecurityTest
	YarnAuditConfig    *SecurityTest
	SafetyConfig       *SecurityTest
	SpotBugsConfig     *SecurityTest
)

// ViperConfig is the struct that stores the caller for testing.
type ViperConfig struct {
	Caller ViperInterface
}

func init() {

	DefaultSecTestConf = &ViperConfig{
		Caller: &ViperCalls{},
	}

	// load Viper using api/config.yml
	if err := DefaultSecTestConf.Caller.SetConfigFile("config", "api/"); err != nil {
		fmt.Println("Error reading Viper config: ", err)
	}

	BanditConfig = DefaultSecTestConf.getSecurityTestConfig("bandit")
	BrakemanConfig = DefaultSecTestConf.getSecurityTestConfig("brakeman")
	GitleaksConfig = DefaultSecTestConf.getSecurityTestConfig("gitleaks")
	GosecConfig = DefaultSecTestConf.getSecurityTestConfig("gosec")
	NpmAuditConfig = DefaultSecTestConf.getSecurityTestConfig("npmaudit")
	YarnAuditConfig = DefaultSecTestConf.getSecurityTestConfig("yarnaudit")
	SafetyConfig = DefaultSecTestConf.getSecurityTestConfig("safety")
	SpotBugsConfig = DefaultSecTestConf.getSecurityTestConfig("spotbugs")

	AllGeneric = append(AllGeneric, GitleaksConfig)
}

// SecurityTest is the struct that stores all data from the securityTest
type SecurityTest struct {
	Name            string                        `bson:"name" json:"name"`
	Type            string                        `bson:"type" json:"type"`
	Language        string                        `bson:"language" json:"language"`
	WarningFound    string                        `bson:"warningFound" json:"warningFound"`
	ErrorFound      string                        `bson:"errorFound" json:"errorFound"`
	Info            string                        `bson:"info" json:"info"`
	Result          string                        `bson:"result" json:"result"`
	Container       container.Container           `bson:"container" json:"container"`
	Vulnerabilities []vulnerability.Vulnerability `bson:"vulnerabilities" json:"vulnerabilities"`
}

// Analyze analyzes the result of a securityTest given a container output
func (s *SecurityTest) Analyze() error {

	if err := s.checkContainerOutputSize(); err != nil {
		return err
	}

	if err := s.checkContainerErrorOrWarning(); err != nil {
		if s.ErrorFound != "" {
			return err
		}
	}

	if err := s.checkVulns(); err != nil {
		return err
	}

	return nil
}

func (s *SecurityTest) checkContainerOutputSize() error {

	maxNumCharsOutput := 1000000

	if len(s.Container.Output) > maxNumCharsOutput {
		s.Result = "error"
		s.Info = "huskyCI was not able to analyze container output: it is huge!"
		s.ErrorFound = "Container output is too big."
		return errors.New("container output is huge")
	}

	return nil
}

func (s *SecurityTest) checkContainerErrorOrWarning() error {

	r, _ := regexp.Compile(`ERROR_\w+`) // #nohusky

	errorMessage := r.FindString(s.Container.Output)
	errorInfo, errorFound := log.MsgCode[errorMessage]

	if errorFound {
		s.Result = "error"
		s.Info = errorInfo
		s.ErrorFound = s.Container.Output
		return errors.New("error found in container")
	}

	r, _ = regexp.Compile(`WARNING_\w+`) // #nohusky

	warningMessage := r.FindString(s.Container.Output)
	warningInfo := log.MsgCode[warningMessage]

	if warningInfo != "" {
		s.Result = "warning"
		s.Info = warningInfo
		s.WarningFound = warningInfo
		return errors.New("warning found in container")
	}

	return nil

}

func (s *SecurityTest) checkVulns() error {

	switch s.Name {
	case "bandit":
		return s.analyzeBandit()
	case "brakeman":
		return s.analyzeBrakeman()
	case "gosec":
		return s.analyzeGosec()
	case "npmaudit":
		return s.analyzeNpmaudit()
	case "yarnaudit":
		return s.analyzeYarnaudit()
	case "spotbugs":
		return s.analyzeSpotBugs()
	case "gitleaks":
		return s.analyzeGitleaks()
	case "safety":
		return s.analyzeSafety()
	}

	return errors.New("invalid securityTest")
}

// GetAllGeneric returns a slice of securityTests containing all
// generic securityTests.
func GetAllGeneric() []*SecurityTest {
	return AllGeneric
}

// GetAllByLanguage returns all generic securityTests based on its type (Language or Generic).
// If Language is used, the second argument must be the name of the language.
func GetAllByLanguage(language string) []*SecurityTest {

	var allByLanguage []*SecurityTest

	switch language {
	case "Go":
		allByLanguage = append(allByLanguage, GosecConfig)
	case "Python":
		allByLanguage = append(allByLanguage, BanditConfig)
		allByLanguage = append(allByLanguage, SafetyConfig)
	case "JavaScript":
		allByLanguage = append(allByLanguage, YarnAuditConfig)
		allByLanguage = append(allByLanguage, NpmAuditConfig)
	case "Java":
		allByLanguage = append(allByLanguage, SpotBugsConfig)
	case "Ruby":
		allByLanguage = append(allByLanguage, BrakemanConfig)
	}

	return allByLanguage
}

func (dF ViperConfig) getSecurityTestConfig(securityTestName string) *SecurityTest {
	return &SecurityTest{
		Name:     dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.name", securityTestName)),
		Type:     dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.type", securityTestName)),
		Language: dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.language", securityTestName)),
		Container: container.Container{
			Command: dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.container.command", securityTestName)),
			Image: container.Image{
				CanonicalURL: dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.container.image.canonicalurl", securityTestName)),
				Name:         dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.container.image.name", securityTestName)),
				Tag:          dF.Caller.GetStringFromConfigFile(fmt.Sprintf("%s.container.image.tag", securityTestName)),
			},
		},
	}
}

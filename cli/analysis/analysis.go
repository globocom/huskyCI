package analysis

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/globocom/huskyCI/cli/config"
	"github.com/globocom/huskyCI/cli/errorcli"
	"github.com/globocom/huskyCI/cli/util"
	"github.com/globocom/huskyCI/cli/vulnerability"
	"github.com/google/uuid"
	"github.com/src-d/enry/v2"
)

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID              string                        `bson:"ID" json:"ID"`
	CompressedFile  CompressedFile                `bson:"compressedFile" json:"compressedFile"`
	Errors          []string                      `bson:"errorsFound,omitempty" json:"errorsFound"`
	Languages       []string                      `bson:"languages" json:"languages"`
	StartedAt       time.Time                     `bson:"startedAt" json:"startedAt"`
	FinishedAt      time.Time                     `bson:"finishedAt" json:"finishedAt"`
	Vulnerabilities []vulnerability.Vulnerability `bson:"vulnerabilities" json:"vulnerabilities"`
	Result          Result                        `bson:"result,omitempty" json:"result"`
}

// CompressedFile holds the info from the compressed file
type CompressedFile struct {
	Name string `bson:"name" json:"name"`
	Size string `bson:"size" json:"size"`
}

// Result holds the status and the info of an analysis.
type Result struct {
	Status string `bson:"status" json:"status"`
	Info   string `bson:"info,omitempty" json:"info"`
}

// New returns a new analysis struct
func New() *Analysis {
	return &Analysis{
		ID: uuid.New().String(),
	}
}

// CheckPath checks the given path to check which languages were found and do some others security checks
func (a *Analysis) CheckPath(path string) error {

	fullPath, err := filepath.Abs(path)
	if err != nil {
		errorcli.Handle(err)
	}

	fmt.Println("[HUSKYCI] Scanning your code from", fullPath)

	if err := a.setLanguages(fullPath); err != nil {
		errorcli.Handle(err)
	}

	for language := range a.getAvailableSecurityTests(a.Languages) {
		s := fmt.Sprintf("[HUSKYCI] %s found.", language)
		fmt.Println(s)
	}

	return nil
}

// CompressFiles will compress all files from a given path into a single file named GUID
func (a *Analysis) CompressFiles(path string) error {

	fmt.Println("[HUSKYCI] Compressing your code...")

	if err := a.HouseCleaning(); err != nil {
		// it's ok. maybe the file is not there yet.
		fmt.Print("")
	}

	allFilesAndDirNames, err := util.GetAllAllowedFilesAndDirsFromPath(path)
	if err != nil {
		return err
	}

	zipFilePath, err := util.CompressFiles(allFilesAndDirNames)
	if err != nil {
		return err
	}

	if err := a.setZipSize(zipFilePath); err != nil {
		return err
	}

	fmt.Println("[HUSKYCI] Compressed! ", zipFilePath, a.CompressedFile.Size)

	return nil
}

// SendZip will send the zip file to the huskyCI API to start the analysis
func (a *Analysis) SendZip() error {
	fmt.Println("[HUSKYCI] Sending your code to the huskyCI API at...")
	defer fmt.Println("[HUSKYCI] Sent!")
	return nil
}

// CheckStatus is a worker to check the huskyCI API for the status of the particular analysis
func (a *Analysis) CheckStatus() error {
	fmt.Println("[HUSKYCI] Checking if the analysis has already finished...")
	defer fmt.Println("[HUSKYCI] Checked!")
	return nil
}

// PrintVulns prints all vulnerabilities found after the analysis has been finished
func (a *Analysis) PrintVulns() {
	fmt.Println("[HUSKYCI] Results:")
}

// HouseCleaning will do stuff to clean the $HOME directory.
func (a *Analysis) HouseCleaning() error {

	zipFilePath, err := config.GetHuskyZipFilePath()
	if err != nil {
		return err
	}

	return util.DeleteHuskyFile(zipFilePath)
}

func (a *Analysis) setZipSize(destination string) error {
	friendlySize, err := util.GetZipFriendlySize(destination)
	if err != nil {
		return err
	}
	a.CompressedFile.Size = friendlySize
	return nil
}

func (a *Analysis) setLanguages(pathReceived string) error {
	err := filepath.Walk(pathReceived,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileName := info.Name()
			lang, _ := enry.GetLanguageByExtension(fileName)
			a.Languages = util.AppendIfMissing(a.Languages, lang)
			return nil
		})
	if err != nil {
		return err
	}
	if len(a.Languages) == 0 {
		return errors.New("no languages found")
	}
	return nil
}

// getAvailableSecurityTests returns the huskyCI securityTests available.
// Later on this check can be done using an API endpoint via cache.
func (a *Analysis) getAvailableSecurityTests(languages []string) map[string][]string {

	var list = make(map[string][]string)

	// Language securityTests
	for _, language := range languages {
		switch language {
		case "Go":
			list[language] = []string{"huskyci/gosec"}
		case "Python":
			list[language] = []string{"huskyci/bandit", "huskyci/safety"}
		case "Ruby":
			list[language] = []string{"huskyci/brakeman"}
		case "JavaScript":
			list[language] = []string{"huskyci/npmaudit", "huskyci/yarnaudit"}
		case "Java":
			list[language] = []string{"huskyci/spotbugs"}
		case "HCL":
			list[language] = []string{"huskyci/tfsec"}
		case "C#":
			list[language] = []string{"huskyci/securitycodescan"}
		}
	}

	// Generic securityTests:
	list["Generic"] = []string{"huskyci/gitleaks"}

	return list
}

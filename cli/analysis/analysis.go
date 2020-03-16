package analysis

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/globocom/huskyCI/cli/errorcli"
	"github.com/globocom/huskyCI/cli/vulnerability"
	"github.com/google/uuid"
	"github.com/mholt/archiver"
	"github.com/src-d/enry/v2"
)

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID              string                        `bson:"ID" json:"ID"`
	CompressedFile  CompressedFile                `bson:"compressedFile" json:"compressedFile"`
	Errors          []string                      `bson:"errorsFound,omitempty" json:"errorsFound"`
	Files           []string                      `bson:"files" json:"files"`
	Languages       []string                      `bson:"languages" json:"languages"`
	StartedAt       time.Time                     `bson:"startedAt" json:"startedAt"`
	FinishedAt      time.Time                     `bson:"finishedAt" json:"finishedAt"`
	Vulnerabilities []vulnerability.Vulnerability `bson:"vulnerabilities" json:"vulnerabilities"`
	Result          Result                        `bson:"result,omitempty" json:"result"`
}

// CompressedFile holds the info from the compressed file
type CompressedFile struct {
	Name string `bson:"name" json:"name"`
	Size int64  `bson:"size" json:"size"`
}

// Result holds the status and the info of an analysis.
type Result struct {
	Status string `bson:"status" json:"status"`
	Info   string `bson:"info,omitempty" json:"info"`
}

// New returns a new analysis struct
func New() *Analysis {

	newID := uuid.New().String()
	newZipFileName := fmt.Sprintf("%s.zip", newID)

	return &Analysis{
		ID: newID,
		CompressedFile: CompressedFile{
			Name: newZipFileName,
		},
	}

}

// CheckPath checks the given path to check which languages were found and do some others security checks
func (a *Analysis) CheckPath(path string) error {

	fullPath, err := filepath.Abs(path)
	if err != nil {
		errorcli.Handle(err)
	}

	fmt.Println("[HUSKYCI] Scanning your code from", fullPath)
	defer fmt.Println("[HUSKYCI] Scanned!")

	if err := a.setFiles(fullPath); err != nil {
		errorcli.Handle(err)
	}

	if err := a.setLanguages(); err != nil {
		errorcli.Handle(err)
	}

	fmt.Println(fmt.Sprintf("[HUSKYCI] %d files found", len(a.Files)))
	for language, securityTests := range a.getAvailableSecurityTests(a.Languages) {
		fmt.Println(fmt.Sprintf("[HUSKYCI] %s -> %s", language, securityTests))
	}

	return nil
}

// CompressFiles will compress all files from a given path into a single file named GUID
func (a *Analysis) CompressFiles(path string) error {

	fmt.Println("[HUSKYCI] Compressing your code...")

	// get all files and folders name inside this path
	var allFilesAndDirNames []string
	filesAndDirs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range filesAndDirs {
		allFilesAndDirNames = append(allFilesAndDirNames, f.Name())
	}

	// compress everthing inside the $HOME directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	destination := fmt.Sprintf("%s/%s", home, a.CompressedFile.Name)
	if err := archiver.Archive(allFilesAndDirNames, destination); err != nil {
		return err
	}

	// get size of the archived file
	file, err := os.Open(destination)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	a.CompressedFile.Size = fi.Size()

	friendlySize := byteCountSI(a.CompressedFile.Size)

	fmt.Println("[HUSKYCI] Compressed! ", destination, friendlySize)
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
	// fmt.Println("[HUSKYCI] 🟢🔵🟡🔴")
}

func (a *Analysis) setFiles(pathReceived string) error {
	err := filepath.Walk(pathReceived,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			a.Files = append(a.Files, info.Name())
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (a *Analysis) setLanguages() error {
	if len(a.Files) == 0 {
		return errors.New("no files found")
	}
	for _, file := range a.Files {
		lang, _ := enry.GetLanguageByExtension(file)
		a.Languages = appendIfMissing(a.Languages, lang)
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
		}
	}

	// Generic securityTests:
	list["Generic"] = []string{"huskyci/gitleaks"}

	return list
}

func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/globocom/huskyCI/cli/config"
	"github.com/globocom/huskyCI/cli/errorcli"
	"github.com/mholt/archiver"
)

// GetAllAllowedFilesAndDirsFromPath returns a list of all files and dirs allowed to be zipped
func GetAllAllowedFilesAndDirsFromPath(path string) ([]string, error) {

	var allFilesAndDirNames []string

	filesAndDirs, err := ioutil.ReadDir(path)
	if err != nil {
		return allFilesAndDirNames, err
	}
	for _, file := range filesAndDirs {
		fileName := file.Name()
		if err := checkFileExtension(fileName); err != nil {
			continue
		} else {
			allFilesAndDirNames = append(allFilesAndDirNames, fileName)
		}
	}

	return allFilesAndDirNames, nil
}

// CompressFiles compress all files into a zip and return its full path and an error
func CompressFiles(allFilesAndDirNames []string) (string, error) {

	var fullFilePath string

	fullFilePath, err := config.GetHuskyZipFilePath()
	if err != nil {
		return fullFilePath, err
	}

	if err := archiver.Archive(allFilesAndDirNames, fullFilePath); err != nil {
		return fullFilePath, err
	}

	return fullFilePath, nil
}

// GetZipFriendlySize returns the size of a friendly zip file size based on its destination
func GetZipFriendlySize(destination string) (string, error) {

	var friendlySize string

	file, err := os.Open(destination) // #nosec -> this destination is always "$HOME/.huskyci/compressed-code.zip"
	if err != nil {
		return friendlySize, err
	}

	fi, err := file.Stat()
	if err != nil {
		return friendlySize, err
	}

	if err := file.Close(); err != nil {
		return friendlySize, err
	}

	friendlySize = byteCountSI(fi.Size())
	return friendlySize, nil
}

// DeleteHuskyFile will delete the huskyCI file present at "$HOME/.huskyci/compressed-code.zip"
func DeleteHuskyFile(destination string) error {
	return os.Remove(destination)
}

func checkFileExtension(file string) error {
	extensionFound := filepath.Ext(file)
	switch extensionFound {
	case "":
		return nil
	case ".jpg", ".png", ".gif", ".webp", ".tiff", ".psd", ".raw", ".bmp", ".heif", ".indd", ".jpeg", ".svg", ".ai", ".eps", ".pdf":
		return errorcli.ErrInvalidExtension
	case ".webm", ".mpg", ".mp2", ".mpeg", ".mpe", ".mpv", ".ogg", ".mp4", ".m4p", ".m4v", ".avi", ".wmv", ".mov", ".qt", ".flv", ".swf", ".avchd":
		return errorcli.ErrInvalidExtension
	default:
		return nil
	}
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

// AppendIfMissing will append an item in a slice if it is missing
func AppendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

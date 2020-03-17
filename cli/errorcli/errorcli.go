package errorcli

import (
	"errors"
	"fmt"
	"os"
)

var (
	// InvalidExtension occurs when an extension is a image and video one
	InvalidExtension = errors.New("invalid extension")
)

// Handle prints the error message in the cli format
func Handle(errorFound error) {
	fmt.Println("[HUSKYCI] ❌ Error found!", errorFound)
	os.Exit(1)
}

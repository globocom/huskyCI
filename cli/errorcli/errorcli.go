package errorcli

import (
	"errors"
	"fmt"
	"os"
)

var (
	// ErrInvalidExtension occurs when an extension is a image and video one
	ErrInvalidExtension = errors.New("invalid extension")
)

// Handle prints the error message in the cli format
func Handle(errorFound error) {
	fmt.Println("[HUSKYCI] ❌ Error found!", errorFound)
	os.Exit(1)
}

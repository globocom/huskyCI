package errorcli

import (
	"fmt"
	"os"
)

// Handle prints the error message in the cli format
func Handle(errorFound error) {
	fmt.Println("[HUSKYCI] ❌ Error found!", errorFound)
	os.Exit(1)
}

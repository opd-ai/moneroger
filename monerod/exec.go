// Package monerod provides functionality for managing the Monero daemon (monerod) process.
// It handles daemon location, execution, and lifecycle management.
package monerod

import (
	"fmt"
	"path/filepath"

	"github.com/opd-ai/moneroger/util"
)

// MoneroDPath searches for the monerod executable in the system path.
// It looks in the following locations in order:
// 1. Directory containing the current executable
// 2. Current working directory
// 3. System PATH directories
//
// Returns:
//   - string: The full path to the monerod executable if found
//   - error: An error if the executable cannot be found
//
// The function checks each directory in the search path for an executable
// named "monerod". On Windows, it will also check for "monerod.exe".
//
// Errors:
//   - Returns an error if monerod is not found in any search location
//
// Example:
//
//	path, err := MoneroDPath()
//	if err != nil {
//	    log.Fatal("monerod not found:", err)
//	}
//	fmt.Println("Found monerod at:", path)
//
// Related:
//   - util.Path() for search path generation
//   - util.FileExists() for file checking
func MoneroDPath() (string, error) {
	paths := util.Path()
	for _, path := range paths {
		monerodPath := filepath.Join(path, "monerod")
		if util.FileExists(monerodPath) {
			return monerodPath, nil
		}
	}
	return "", fmt.Errorf("Monero daemon(monerod) not found")
}

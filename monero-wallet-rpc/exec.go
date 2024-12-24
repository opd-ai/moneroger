// Package monerowalletrpc provides functionality for managing the Monero wallet RPC service.
// It handles service location, execution, and lifecycle management.
package monerowalletrpc

import (
	"fmt"
	"path/filepath"

	"github.com/opd-ai/moneroger/util"
)

// MoneroWalletRPCPath searches for the monero-wallet-rpc executable in the system path.
// It looks in the following locations in order:
// 1. Directory containing the current executable
// 2. Current working directory
// 3. System PATH directories
//
// Returns:
//   - string: The full path to the monero-wallet-rpc executable if found
//   - error: An error if the executable cannot be found
//
// Search behavior:
//   - Uses util.Path() to get search directories
//   - Checks each directory for the executable
//   - Returns first match found
//   - Case-sensitive on Unix-like systems
//   - Checks for .exe extension on Windows automatically
//
// Errors:
//   - Returns descriptive error if executable is not found in any location
//
// Example:
//
//	path, err := MoneroWalletRPCPath()
//	if err != nil {
//	    log.Fatal("monero-wallet-rpc not found:", err)
//	}
//	fmt.Println("Found wallet RPC at:", path)
//
// Related:
//   - util.Path() for search path generation
//   - util.FileExists() for file checking
//   - github.com/opd-ai/moneroger/monerod.MoneroDPath() for daemon executable
func MoneroWalletRPCPath() (string, error) {
	paths := util.Path()
	for _, path := range paths {
		moneroWalletRPCPath := filepath.Join(path, "monero-wallet-rpc")
		if util.FileExists(moneroWalletRPCPath) {
			return moneroWalletRPCPath, nil
		}
	}
	return "", fmt.Errorf("Monero wallet RPC(monero-wallet-rpc) not found")
}

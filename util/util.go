// Package util provides utility functions for the moneroger library,
// including file operations, path management, password generation,
// and network port handling.
package util

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	moneroconst "github.com/opd-ai/moneroger/const"
	"github.com/sethvargo/go-password/password"
)

// FileExists checks if a file exists at the specified path.
//
// Parameters:
//   - path: The file path to check (string)
//
// Returns:
//   - bool: true if the file exists and is accessible, false otherwise
//
// Note: This function returns false for both non-existent files
// and files that exist but are inaccessible due to permissions.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// Path returns a slice of directories to search for executables.
// It combines:
// - The directory containing the current executable
// - The current working directory
// - The system PATH environment variable
//
// Returns:
//   - []string: Slice of directory paths to search
//
// Note: Logs errors but continues execution if unable to determine
// executable or working directory paths.
func Path() []string {
	path := os.Getenv("PATH")
	elements := []string{}

	// Get executable directory
	me, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path:", err)
	} else {
		meDir := filepath.Dir(me)
		elements = append(elements, meDir)
	}

	// Get working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory:", err)
	} else {
		elements = append(elements, workDir)
	}

	// Add system PATH elements
	if path != "" {
		elements = append(elements, strings.Split(path, ":")...)
	}

	return elements
}

// SecurePassword generates a cryptographically secure random password.
//
// Returns:
//   - string: A 20-character password with random number of digits (0-19)
//
// Panics:
//   - If the password generation fails (should be extremely rare)
//
// Uses:
//   - github.com/sethvargo/go-password/password for generation
func SecurePassword() string {
	rand.Seed(time.Now().UnixNano())
	digs := rand.Intn(19)
	res, err := password.Generate(20, digs, 0, false, false)
	if err != nil {
		panic(err)
	}
	return res
}

// IsPortInUse checks if a TCP port is currently in use on localhost.
//
// Parameters:
//   - port: Port number to check (int)
//
// Returns:
//   - bool: true if the port is in use, false otherwise
//
// Note: This function attempts a TCP connection with a 1-second timeout.
// A successful connection indicates the port is in use.
func IsPortInUse(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// WaitForPort waits for a TCP port to become available.
//
// Parameters:
//   - ctx: Context for cancellation
//   - port: Port number to wait for (int)
//
// Returns:
//   - error: nil if port becomes available, error otherwise
//
// Errors:
//   - Context cancellation error if context is cancelled
//   - Timeout error if port doesn't become available within DefaultStartupTimeout
//
// Related:
//   - moneroconst.DefaultStartupTimeout
//   - IsPortInUse function
func WaitForPort(ctx context.Context, port int) error {
	deadline := time.Now().Add(moneroconst.DefaultStartupTimeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if IsPortInUse(port) {
				return nil
			}
			time.Sleep(time.Second)
		}
	}
	return fmt.Errorf("timeout waiting for port %d", port)
}

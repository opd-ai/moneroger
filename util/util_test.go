package util

import (
	"context"
	"net"
	"os"
	"testing"
)

// TestFileExists verifies the FileExists function correctly identifies
// existing and non-existing files
func TestFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", tmpFile.Name(), true},
		{"non-existing file", "/path/to/nonexistent/file", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.path); got != tt.expected {
				t.Errorf("FileExists(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

// TestPath verifies the Path function includes required directories
func TestPath(t *testing.T) {
	paths := Path()

	// Should contain at least working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, p := range paths {
		if p == wd {
			found = true
			break
		}
	}

	if !found {
		t.Error("Path() should include working directory")
	}
}

// TestSecurePassword verifies password generation requirements
func TestSecurePassword(t *testing.T) {
	p1 := SecurePassword()
	p2 := SecurePassword()

	// Check length
	if len(p1) != 20 {
		t.Errorf("SecurePassword() length = %d, want 20", len(p1))
	}

	// Check uniqueness
	if p1 == p2 {
		t.Error("SecurePassword() should generate unique passwords")
	}
}

// TestIsPortInUse verifies port availability checking
func TestIsPortInUse(t *testing.T) {
	// Start a listener on a random port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	tests := []struct {
		name     string
		port     int
		expected bool
	}{
		{"used port", port, true},
		{"unused port", 0, false}, // Port 0 should never be in use
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPortInUse(tt.port); got != tt.expected {
				t.Errorf("IsPortInUse(%d) = %v, want %v", tt.port, got, tt.expected)
			}
		})
	}
}

// TestWaitForPort verifies port waiting behavior
func TestWaitForPort(t *testing.T) {
	// Test immediate availability
	t.Run("port immediately available", func(t *testing.T) {
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()

		port := listener.Addr().(*net.TCPAddr).Port
		ctx := context.Background()

		if err := WaitForPort(ctx, port); err != nil {
			t.Errorf("WaitForPort() error = %v", err)
		}
	})

	// Test context cancellation
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		if err := WaitForPort(ctx, 0); err == nil {
			t.Error("WaitForPort() should return error on cancelled context")
		}
	})
}

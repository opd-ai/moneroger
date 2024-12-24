package monerod

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opd-ai/moneroger/util"
)

// TestMoneroDPath tests the daemon executable search functionality
func TestMoneroDPath(t *testing.T) {
	// Create a temporary directory with test executable
	tmpDir, err := os.MkdirTemp("", "monerod-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock monerod executable
	mockPath := filepath.Join(tmpDir, "monerod")
	f, err := os.Create(mockPath)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	os.Chmod(mockPath, 0755)

	// Temporarily modify PATH
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	// Test finding the executable
	path, err := MoneroDPath()
	if err != nil {
		t.Errorf("MoneroDPath() error = %v", err)
	}
	if path != mockPath {
		t.Errorf("MoneroDPath() = %v, want %v", path, mockPath)
	}
}

// TestRPCCredentials tests credential management
func TestRPCCredentials(t *testing.T) {
	t.Run("default username", func(t *testing.T) {
		d := &MoneroDaemon{}
		if user := d.RPCUser(); user != "gouser" {
			t.Errorf("RPCUser() = %v, want gouser", user)
		}
	})

	t.Run("custom username", func(t *testing.T) {
		d := &MoneroDaemon{rpcUser: "custom"}
		if user := d.RPCUser(); user != "custom" {
			t.Errorf("RPCUser() = %v, want custom", user)
		}
	})

	t.Run("password generation", func(t *testing.T) {
		d := &MoneroDaemon{}
		pass1 := d.RPCPass()
		if pass1 == "" {
			t.Error("RPCPass() returned empty password")
		}

		// Password should be consistent
		pass2 := d.RPCPass()
		if pass1 != pass2 {
			t.Error("RPCPass() returned inconsistent passwords")
		}
	})
}

// TestNewMoneroDaemon tests daemon creation and configuration
func TestNewMoneroDaemon(t *testing.T) {
	// Create temporary data directory
	dataDir, err := os.MkdirTemp("", "monero-data-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dataDir)

	tests := []struct {
		name    string
		config  util.Config
		wantErr bool
	}{
		{
			name: "basic mainnet config",
			config: util.Config{
				DataDir:    dataDir,
				MoneroPort: 18081,
			},
			wantErr: false,
		},
		{
			name: "testnet config",
			config: util.Config{
				DataDir:    dataDir,
				MoneroPort: 28081,
				TestNet:    true,
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: util.Config{
				DataDir:    dataDir,
				MoneroPort: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			daemon, err := NewMoneroDaemon(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoneroDaemon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if daemon.dataDir != tt.config.DataDir {
					t.Errorf("dataDir = %v, want %v", daemon.dataDir, tt.config.DataDir)
				}
				if daemon.testnet != tt.config.TestNet {
					t.Errorf("testnet = %v, want %v", daemon.testnet, tt.config.TestNet)
				}
				// Clean up
				_ = daemon.Shutdown(ctx)
			}
		})
	}
}

// TestShutdown tests daemon shutdown behavior
func TestShutdown(t *testing.T) {
	t.Run("nil process", func(t *testing.T) {
		d := &MoneroDaemon{}
		ctx := context.Background()
		if err := d.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	})
}

// TestDefaultConstants verifies constant values
func TestDefaultConstants(t *testing.T) {
	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"default port", defaultMonerodPort, 18081},
		{"startup timeout", defaultStartupTimeout, 30 * time.Second},
		{"shutdown timeout", defaultShutdownTimeout, 10 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

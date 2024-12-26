package monerowalletrpc

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opd-ai/moneroger/monerod"
	"github.com/opd-ai/moneroger/util"
)

// createTestFile creates a temporary file for testing
func createTestFile(t *testing.T, prefix string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

// TestValidateConfig tests configuration validation
func TestValidateConfig(t *testing.T) {
	// Create a temporary wallet file
	walletFile := createTestFile(t, "wallet-*")
	defer os.Remove(walletFile)

	tests := []struct {
		name    string
		config  util.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: util.Config{
				WalletFile: walletFile,
				WalletPort: 18082,
			},
			wantErr: false,
		},
		{
			name: "empty wallet file",
			config: util.Config{
				WalletFile: "",
				WalletPort: 18082,
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: util.Config{
				WalletFile: walletFile,
				WalletPort: -1,
			},
			wantErr: true,
		},
		{
			name: "non-existent wallet file",
			config: util.Config{
				WalletFile: "/nonexistent/wallet",
				WalletPort: 18082,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestWalletRPCCredentials tests RPC credential management
func TestWalletRPCCredentials(t *testing.T) {
	t.Run("default username", func(t *testing.T) {
		w := &WalletRPC{}
		if user := w.WalletRPCUser(); user != "gouser" {
			t.Errorf("WalletRPCUser() = %v, want gouser", user)
		}
	})

	t.Run("custom username", func(t *testing.T) {
		w := &WalletRPC{rpcUser: "custom"}
		if user := w.WalletRPCUser(); user != "custom" {
			t.Errorf("WalletRPCUser() = %v, want custom", user)
		}
	})

	t.Run("password generation", func(t *testing.T) {
		w := &WalletRPC{}
		pass1 := w.WalletRPCPass()
		if pass1 == "" {
			t.Error("WalletRPCPass() returned empty password")
		}

		// Password should be consistent
		pass2 := w.WalletRPCPass()
		if pass1 != pass2 {
			t.Error("WalletRPCPass() returned inconsistent passwords")
		}
	})
}

// TestWalletState tests wallet state string representations
func TestWalletState(t *testing.T) {
	tests := []struct {
		state WalletState
		want  string
	}{
		{WalletStateUnknown, "unknown"},
		{WalletStateStarting, "starting"},
		{WalletStateRunning, "running"},
		{WalletStateStopping, "stopping"},
		{WalletStateStopped, "stopped"},
		{WalletState(99), "unknown"}, // Invalid state
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("WalletState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMoneroWalletRPCPath tests executable path resolution
func TestMoneroWalletRPCPath(t *testing.T) {
	// Create a temporary directory with mock executable
	tmpDir, err := os.MkdirTemp("", "wallet-rpc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock executable
	mockPath := filepath.Join(tmpDir, "monero-wallet-rpc")
	f, err := os.Create(mockPath)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	os.Chmod(mockPath, 0o755)

	// Temporarily modify PATH
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	path, err := MoneroWalletRPCPath()
	if err != nil {
		t.Errorf("MoneroWalletRPCPath() error = %v", err)
	}
	if path != mockPath {
		t.Errorf("MoneroWalletRPCPath() = %v, want %v", path, mockPath)
	}
}

// TestWalletRPCShutdown tests shutdown behavior
func TestWalletRPCShutdown(t *testing.T) {
	t.Run("nil process", func(t *testing.T) {
		w := &WalletRPC{}
		ctx := context.Background()
		if err := w.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	})

	t.Run("shutdown timeout", func(t *testing.T) {
		w := &WalletRPC{}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		if err := w.Shutdown(ctx); err != nil {
			// Expect timeout or clean shutdown
			t.Log("Expected shutdown behavior:", err)
		}
	})
}

// MockDaemon creates a mock daemon for testing
func MockDaemon(t *testing.T) *monerod.MoneroDaemon {
	t.Helper()
	return &monerod.MoneroDaemon{}
}

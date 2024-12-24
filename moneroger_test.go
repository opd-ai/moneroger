package moneroger

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/opd-ai/moneroger/util"
)

// createTestConfig creates a test configuration with temporary directories
func createTestConfig(t *testing.T) util.Config {
	t.Helper()

	// Create temporary directory for data
	dataDir, err := os.MkdirTemp("", "moneroger-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dataDir) })

	// Create temporary wallet file
	walletFile, err := os.CreateTemp(dataDir, "wallet-*.keys")
	if err != nil {
		t.Fatal(err)
	}
	walletFile.Close()

	return util.Config{
		DataDir:    dataDir,
		WalletFile: walletFile.Name(),
		MoneroPort: 18081,
		WalletPort: 18082,
		TestNet:    true, // Use testnet for tests
	}
}

// TestNewMoneroger tests the creation of a new Moneroger instance
func TestNewMoneroger(t *testing.T) {
	tests := []struct {
		name    string
		config  util.Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  createTestConfig(t),
			wantErr: false,
		},
		{
			name: "invalid daemon port",
			config: util.Config{
				DataDir:    "testdata",
				WalletFile: "testdata/wallet.keys",
				MoneroPort: -1, // Invalid port
				WalletPort: 18082,
			},
			wantErr: true,
		},
		{
			name: "missing wallet file",
			config: util.Config{
				DataDir:    "testdata",
				WalletFile: "", // Missing wallet file
				MoneroPort: 18081,
				WalletPort: 18082,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoneroger(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoneroger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if m == nil {
					t.Error("NewMoneroger() returned nil manager without error")
				}
				// Clean up
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = m.Shutdown(ctx)
			}
		})
	}
}

// TestStartupShutdownSequence tests the proper ordering of operations
func TestStartupShutdownSequence(t *testing.T) {
	config := createTestConfig(t)
	m, err := NewMoneroger(config)
	if err != nil {
		t.Fatal(err)
	}

	// Test startup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Start(ctx); err != nil {
		t.Errorf("Start() error = %v", err)
	}

	// Test shutdown
	if err := m.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

// TestContextCancellation verifies context handling
func TestContextCancellation(t *testing.T) {
	config := createTestConfig(t)
	m, err := NewMoneroger(config)
	if err != nil {
		t.Fatal(err)
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Start should respect context cancellation
	if err := m.Start(ctx); err == nil {
		t.Error("Start() should fail with cancelled context")
	}

	// Shutdown should still work with new context
	ctx = context.Background()
	if err := m.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

// TestConcurrentAccess verifies thread safety
func TestConcurrentAccess(t *testing.T) {
	config := createTestConfig(t)
	m, err := NewMoneroger(config)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- m.Start(ctx)
	}()

	// Attempt immediate shutdown
	if err := m.Shutdown(ctx); err != nil {
		t.Errorf("Concurrent Shutdown() error = %v", err)
	}

	// Check start result
	if err := <-errChan; err != nil {
		t.Errorf("Concurrent Start() error = %v", err)
	}
}

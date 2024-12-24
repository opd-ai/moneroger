package moneroconst

import (
	"testing"
	"time"
)

// yeah they're pointless but bitches love tests

// TestDefaultPorts verifies the standard Monero daemon port assignments
func TestDefaultPorts(t *testing.T) {
	// Test monerod port value
	if DefaultMonerodPort != 18081 {
		t.Errorf("DefaultMonerodPort = %d, want 18081", DefaultMonerodPort)
	}

	// Test wallet RPC port value
	if DefaultWalletRPCPort != 18082 {
		t.Errorf("DefaultWalletRPCPort = %d, want 18082", DefaultWalletRPCPort)
	}

	// Verify ports are different
	if DefaultMonerodPort == DefaultWalletRPCPort {
		t.Error("Monerod and Wallet RPC ports should be different")
	}
}

// TestTimeoutValues verifies the timeout durations
func TestTimeoutValues(t *testing.T) {
	// Test startup timeout duration
	expectedStartup := 30 * time.Second
	if DefaultStartupTimeout != expectedStartup {
		t.Errorf("DefaultStartupTimeout = %v, want %v", DefaultStartupTimeout, expectedStartup)
	}

	// Test shutdown timeout duration
	expectedShutdown := 10 * time.Second
	if DefaultShutdownTimeout != expectedShutdown {
		t.Errorf("DefaultShutdownTimeout = %v, want %v", DefaultShutdownTimeout, expectedShutdown)
	}

	// Verify startup timeout is longer than shutdown timeout
	if DefaultStartupTimeout <= DefaultShutdownTimeout {
		t.Error("Startup timeout should be longer than shutdown timeout")
	}
}

// TestPortRange verifies ports are in valid range
func TestPortRange(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"monerod port", DefaultMonerodPort},
		{"wallet RPC port", DefaultWalletRPCPort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.port <= 0 || tt.port > 65535 {
				t.Errorf("Port %d is outside valid range (1-65535)", tt.port)
			}
		})
	}
}

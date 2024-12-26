// Package moneroconst provides default configuration constants for the moneroger library.
// It defines standard ports and timeouts used by both monerod and monero-wallet-rpc daemons.
package moneroconst

import (
	"time"
)

// Default configurations for Monero daemons
const (
	// DefaultMonerodPort is the standard RPC port for monerod daemon (18081)
	// This port is used for communication between the daemon and wallet
	DefaultMonerodPort = 18081

	// DefaultWalletRPCPort is the standard RPC port for monero-wallet-rpc daemon (18082)
	// This port is used by applications to communicate with the wallet
	DefaultWalletRPCPort = 18083

	// DefaultStartupTimeout defines how long to wait for daemons to start (30 seconds)
	// If a daemon doesn't respond within this time, startup is considered failed
	DefaultStartupTimeout = 30 * time.Second

	// DefaultShutdownTimeout defines how long to wait for graceful shutdown (10 seconds)
	// After this timeout, the process will be forcefully terminated
	DefaultShutdownTimeout = 10 * time.Second
)

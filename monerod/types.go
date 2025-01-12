// Package monerod provides functionality for managing Monero daemon processes.
// It handles daemon configuration, lifecycle management, and RPC communication.
package monerod

import (
	"os/exec"
	"time"

	"github.com/opd-ai/moneroger/util"
)

// Default configurations for the Monero daemon
const (
	// defaultMonerodPort is the default RPC port for monerod (mainnet)
	defaultMonerodPort = 18081

	// defaultStartupTimeout is the maximum time to wait for daemon startup
	defaultStartupTimeout = 30 * time.Second

	// defaultShutdownTimeout is the maximum time to wait for graceful shutdown
	defaultShutdownTimeout = 10 * time.Second
)

// MoneroDaemon represents a running monerod instance and manages its lifecycle.
// It provides access to daemon configuration and process control.
//
// Fields:
//   - cmd: Command instance for process management
//   - dataDir: Directory for blockchain and configuration data
//   - rpcPort: Port number for RPC interface
//   - rpcUser: Username for RPC authentication
//   - rpcPass: Password for RPC authentication
//   - testnet: Boolean flag for testnet operation
//   - process: Reference to the running daemon process
//
// The daemon can be configured for either mainnet or testnet operation,
// with appropriate default ports and network settings applied automatically.
type MoneroDaemon struct {
	cmd           *exec.Cmd
	dataDir       string
	rpcPort       int
	rpcUser       string
	rpcPass       string
	testnet       bool
	useRemoteNode bool
}

// RPCPort returns the configured RPC port for the daemon.
// If no port was explicitly set, returns the default port (18081 for mainnet).
//
// Returns:
//   - int: The RPC port number
//
// Related:
//   - defaultMonerodPort constant
func (m *MoneroDaemon) RPCPort() int {
	return m.rpcPort
}

// RPCUser returns the RPC authentication username.
// If no username was set, initializes it to the default "gouser".
//
// Returns:
//   - string: The RPC username
//
// Note: This method has side effects - it will set a default username
// if one hasn't been configured.
func (m *MoneroDaemon) RPCUser() string {
	if m.rpcUser == "" {
		m.rpcUser = "gouser"
	}
	return m.rpcUser
}

// RPCPass returns the RPC authentication password.
// If no password was set, generates a secure random password using util.SecurePassword().
//
// Returns:
//   - string: The RPC password
//
// Note: This method has side effects - it will generate and set a secure
// password if one hasn't been configured.
//
// Related:
//   - util.SecurePassword() for password generation
func (m *MoneroDaemon) RPCPass() string {
	if m.rpcPass == "" {
		m.rpcPass = util.SecurePassword()
	}
	return m.rpcPass
}

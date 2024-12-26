// Package monerowalletrpc provides functionality for managing Monero wallet RPC services.
package monerowalletrpc

import (
	"os"
	"os/exec"

	"github.com/opd-ai/moneroger/monerod"
	"github.com/opd-ai/moneroger/util"
)

// WalletRPC represents a running monero-wallet-rpc instance and manages its lifecycle.
// It handles RPC configuration, process management, and daemon communication.
//
// Fields:
//   - cmd: Command instance for process management
//   - walletFile: Path to the wallet file (.keys)
//   - rpcPort: Port number for RPC interface
//   - rpcUser: Username for RPC authentication
//   - rpcPass: Password for RPC authentication
//   - daemon: Reference to associated monerod instance
//   - process: Reference to the running wallet RPC process
//
// The WalletRPC instance maintains connection settings and process state,
// coordinating with the Monero daemon for blockchain access.
type WalletRPC struct {
	cmd        *exec.Cmd
	walletDir  string
	rpcPort    int
	rpcUser    string
	rpcPass    string
	walletPass string
	daemon     *monerod.MoneroDaemon
	process    *os.Process
}

// WalletState represents the current operational state of the wallet RPC service.
// It provides a type-safe enumeration of possible wallet states.
type WalletState uint8

// Wallet state constants define the possible states of a wallet RPC service.
const (
	WalletStateUnknown  WalletState = iota // Initial or unknown state
	WalletStateStarting                    // Service is starting up
	WalletStateRunning                     // Service is operational
	WalletStateStopping                    // Service is shutting down
	WalletStateStopped                     // Service has stopped
)

// String returns a human-readable representation of the wallet state.
// This implements the Stringer interface for WalletState.
//
// Returns:
//   - string: A lowercase string description of the current state
//
// States:
//   - "starting": Wallet is initializing
//   - "running": Wallet is operational
//   - "stopping": Wallet is shutting down
//   - "stopped": Wallet has terminated
//   - "unknown": State cannot be determined
func (s WalletState) String() string {
	switch s {
	case WalletStateStarting:
		return "starting"
	case WalletStateRunning:
		return "running"
	case WalletStateStopping:
		return "stopping"
	case WalletStateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// WalletRPCPort returns the configured RPC port for the wallet service.
//
// Returns:
//   - int: The RPC port number
func (m *WalletRPC) WalletRPCPort() int {
	return m.rpcPort
}

// WalletRPCUser returns the RPC authentication username.
// If no username was set, initializes it to the default "gouser".
//
// Returns:
//   - string: The RPC username
//
// Note: This method has side effects - it will set a default username
// if one hasn't been configured.
func (m *WalletRPC) WalletRPCUser() string {
	if m.rpcUser == "" {
		m.rpcUser = "gouser"
	}
	return m.rpcUser
}

// WalletRPCPass returns the RPC authentication password.
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
func (m *WalletRPC) WalletRPCPass() string {
	if m.rpcPass == "" {
		m.rpcPass = util.SecurePassword()
	}
	return m.rpcPass
}

func (m *WalletRPC) WalletPass() string {
	if m.walletPass == "" {
		m.walletPass = "changeme"
	}
	return m.rpcPass
}

func (m *WalletRPC) SetWalletPass(pass string) error {
	m.walletPass = pass
	return nil
}

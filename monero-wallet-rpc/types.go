package monerowalletrpc

import (
	"os"
	"os/exec"

	"github.com/opd-ai/moneroger/monerod"
	"github.com/opd-ai/moneroger/util"
)

// WalletRPC represents a running monero-wallet-rpc instance
type WalletRPC struct {
	cmd        *exec.Cmd
	walletFile string
	rpcPort    int
	rpcUser    string
	rpcPass    string
	daemon     *monerod.MoneroDaemon
	process    *os.Process
}

// WalletState represents the current state of the wallet
type WalletState uint8

const (
	WalletStateUnknown WalletState = iota
	WalletStateStarting
	WalletStateRunning
	WalletStateStopping
	WalletStateStopped
)

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

func (m *WalletRPC) WalletRPCPort() int {
	return m.rpcPort
}

func (m *WalletRPC) WalletRPCUser() string {
	if m.rpcUser == "" {
		m.rpcUser = "gouser"
	}
	return m.rpcUser
}

func (m *WalletRPC) WalletRPCPass() string {
	if m.rpcPass == "" {
		m.rpcPass = util.SecurePassword()
	}
	return m.rpcPass
}

package monerod

import (
	"os"
	"os/exec"
	"time"

	"github.com/opd-ai/moneroger/util"
)

// Default configurations
const (
	defaultMonerodPort     = 18081
	defaultStartupTimeout  = 30 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

// MoneroDaemon represents a running monerod instance
type MoneroDaemon struct {
	cmd     *exec.Cmd
	dataDir string
	rpcPort int
	rpcUser string
	rpcPass string
	testnet bool
	process *os.Process
}

func (m *MoneroDaemon) RPCPort() int {
	return m.rpcPort
}

func (m *MoneroDaemon) RPCUser() string {
	if m.rpcUser == "" {
		m.rpcUser = "gouser"
	}
	return m.rpcUser
}

func (m *MoneroDaemon) RPCPass() string {
	if m.rpcPass == "" {
		m.rpcPass = util.SecurePassword()
	}
	return m.rpcPass
}

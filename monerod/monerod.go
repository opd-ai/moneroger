// Package monerod provides functionality for managing Monero daemon processes.
package monerod

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/opd-ai/moneroger/errors"
	"github.com/opd-ai/moneroger/util"
)

// NewMoneroDaemon creates or connects to a Monero daemon instance.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - config: Configuration settings for the daemon including:
//   - DataDir: Directory for blockchain and wallet data
//   - MoneroPort: RPC port number
//   - TestNet: Boolean flag for testnet operation
//
// Returns:
//   - *MoneroDaemon: Pointer to the daemon instance
//   - error: Any error encountered during startup
//
// The function will:
// 1. Check if a daemon is already running on the specified port
// 2. If running, return a connection to the existing daemon
// 3. If not running, start a new daemon process
//
// Errors:
//   - Process spawn failures
//   - Port binding issues
//   - Context cancellation
//
// Related:
//   - util.Config for configuration options
//   - util.IsPortInUse for port checking
func NewMoneroDaemon(ctx context.Context, config util.Config) (*MoneroDaemon, error) {
	// Check if daemon is already running
	if util.IsPortInUse(config.MoneroPort) {
		return &MoneroDaemon{
			rpcPort: config.MoneroPort,
			dataDir: config.DataDir,
			testnet: config.TestNet,
		}, nil
	}

	daemon := &MoneroDaemon{
		dataDir: config.DataDir,
		rpcPort: config.MoneroPort,
		testnet: config.TestNet,
	}

	if err := daemon.Start(ctx); err != nil {
		return nil, errors.E(
			errors.OpStart,
			errors.ComponentMonerod,
			errors.KindProcess,
			err,
		)
	}

	return daemon, nil
}

// Start launches the monerod process with appropriate configuration.
// This is an internal method used by NewMoneroDaemon.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error encountered during startup
//
// The method will:
// 1. Configure daemon arguments
// 2. Launch the monerod process
// 3. Wait for RPC port availability
//
// Related:
//   - MoneroDPath for executable location
//   - util.WaitForPort for startup confirmation
func (m *MoneroDaemon) Start(ctx context.Context) error {
	args := []string{
		"--data-dir", m.dataDir,
		"--rpc-bind-port", fmt.Sprintf("%d", m.RPCPort()),
		"--rpc-login", fmt.Sprintf("%s:%s", m.RPCUser(), m.RPCPass()),
		"--non-interactive",
	}

	if m.testnet {
		args = append(args, "--testnet")
	}
	moneroD, err := MoneroDPath()
	if err != nil {
		return errors.E(
			errors.OpProcessSpawn,
			errors.ComponentMonerod,
			errors.KindProcess,
			err,
		)
	}
	cmd := exec.CommandContext(ctx, moneroD, args...)
	if err := cmd.Start(); err != nil {
		return errors.E(
			errors.OpProcessSpawn,
			errors.ComponentMonerod,
			errors.KindProcess,
			err,
		)
	}

	m.cmd = cmd
	m.process = cmd.Process

	// Wait for RPC to become available
	if err := util.WaitForPort(ctx, m.RPCPort()); err != nil {
		return errors.E(
			errors.OpPortBinding,
			errors.ComponentMonerod,
			errors.KindNetwork,
			err,
		)
	}

	return nil
}

// Shutdown gracefully stops the Monero daemon.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error encountered during shutdown
//
// The method sends an interrupt signal (SIGINT) to the daemon process,
// allowing it to clean up and shut down gracefully. If the process
// isn't running, the method returns nil.
//
// Errors:
//   - Signal delivery failures
//   - Context cancellation
func (m *MoneroDaemon) Shutdown(ctx context.Context) error {
	if m.process != nil {
		if err := m.process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("failed to send interrupt to monerod: %w", err)
		}
	}
	return nil
}

func (m *MoneroDaemon) PID() string {
	if m.cmd != nil {
		if m.cmd.Process != nil {
			return fmt.Sprintf("%d", m.cmd.Process.Pid)
		}
	}
	return "-1"
}

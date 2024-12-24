package monerod

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/opd-ai/moneroger/errors"
	"github.com/opd-ai/moneroger/util"
)

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

	if err := daemon.start(ctx); err != nil {
		return nil, errors.E(
			errors.OpStart,
			errors.ComponentMonerod,
			errors.KindProcess,
			err,
		)
	}

	return daemon, nil
}

func (m *MoneroDaemon) start(ctx context.Context) error {
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

func (m *MoneroDaemon) Shutdown(ctx context.Context) error {
	if m.process != nil {
		if err := m.process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("failed to send interrupt to monerod: %w", err)
		}
	}
	return nil
}

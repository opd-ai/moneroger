package monerowalletrpc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/opd-ai/moneroger/errors"
	"github.com/opd-ai/moneroger/monerod"
	"github.com/opd-ai/moneroger/util"
)

// Common operations for wallet-rpc
const (
	opStart          = errors.Op("WalletRPC.Start")
	opShutdown       = errors.Op("WalletRPC.Shutdown")
	opValidateConfig = errors.Op("WalletRPC.ValidateConfig")
	opCheckHealth    = errors.Op("WalletRPC.CheckHealth")
)

// NewWalletRPC creates and starts a new WalletRPC instance
func NewWalletRPC(ctx context.Context, config util.Config, daemon *monerod.MoneroDaemon) (*WalletRPC, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	wallet := &WalletRPC{
		walletFile: config.WalletFile,
		rpcPort:    config.WalletPort,
		daemon:     daemon,
	}

	if err := wallet.start(ctx); err != nil {
		return nil, err
	}

	return wallet, nil
}

func validateConfig(config util.Config) error {
	if config.WalletFile == "" {
		return errors.E(
			opValidateConfig,
			errors.ComponentWalletRPC,
			errors.KindConfig,
			fmt.Errorf("wallet file path cannot be empty"),
		)
	}

	if config.WalletPort <= 0 {
		return errors.E(
			opValidateConfig,
			errors.ComponentWalletRPC,
			errors.KindConfig,
			fmt.Errorf("invalid wallet RPC port: %d", config.WalletPort),
		)
	}

	// Check if wallet file exists
	if _, err := os.Stat(config.WalletFile); os.IsNotExist(err) {
		return errors.E(
			opValidateConfig,
			errors.ComponentWalletRPC,
			errors.KindSystem,
			fmt.Errorf("wallet file does not exist: %s", config.WalletFile),
		)
	}

	return nil
}

// start starts the wallet-rpc process
func (w *WalletRPC) start(ctx context.Context) error {
	// Check if port is already in use
	if util.IsPortInUse(w.WalletRPCPort()) {
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindNetwork,
			fmt.Errorf("port %d is already in use", w.WalletRPCPort()),
		)
	}

	args := []string{
		"--wallet-file", w.walletFile,
		"--rpc-bind-port", fmt.Sprintf("%d", w.WalletRPCPort()),
		"--daemon-address", fmt.Sprintf("http://localhost:%d", w.daemon.RPCPort()),
		"--daemon-login", fmt.Sprintf("%s:%s", w.daemon.RPCUser(), w.daemon.RPCPass()),
		"--rpc-login", fmt.Sprintf("%s:%s", w.WalletRPCUser(), w.WalletRPCPass()),
	}
	moneroWalletRPC, err := MoneroWalletRPCPath()
	if err != nil {
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindProcess,
			fmt.Errorf("failed to start wallet-rpc process: %w", err),
		)
	}
	cmd := exec.CommandContext(ctx, moneroWalletRPC, args...)

	// Start the process
	if err := cmd.Start(); err != nil {
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindProcess,
			fmt.Errorf("failed to start wallet-rpc process: %w", err),
		)
	}

	w.cmd = cmd
	w.process = cmd.Process

	// Wait for RPC to become available with timeout
	if err := util.WaitForPort(ctx, w.WalletRPCPort()); err != nil {
		// Try to clean up the process if port binding fails
		_ = w.Shutdown(ctx)
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindTimeout,
			fmt.Errorf("wallet-rpc failed to bind to port %d: %w", w.WalletRPCPort(), err),
		)
	}

	// Verify the wallet is responding correctly
	if err := w.checkHealth(ctx); err != nil {
		_ = w.Shutdown(ctx)
		return err
	}

	return nil
}

// Shutdown gracefully stops the wallet-rpc daemon
func (w *WalletRPC) Shutdown(ctx context.Context) error {
	if w.process == nil {
		return nil
	}

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Send interrupt signal
	if err := w.process.Signal(os.Interrupt); err != nil {
		return errors.E(
			opShutdown,
			errors.ComponentWalletRPC,
			errors.KindProcess,
			fmt.Errorf("failed to send interrupt signal: %w", err),
		)
	}

	// Wait for process to exit
	done := make(chan error, 1)
	go func() {
		_, err := w.process.Wait()
		done <- err
	}()

	select {
	case <-ctx.Done():
		return errors.E(
			opShutdown,
			errors.ComponentWalletRPC,
			errors.KindTimeout,
			fmt.Errorf("shutdown timed out"),
		)
	case err := <-done:
		if err != nil {
			return errors.E(
				opShutdown,
				errors.ComponentWalletRPC,
				errors.KindProcess,
				fmt.Errorf("error during shutdown: %w", err),
			)
		}
	}

	w.process = nil
	w.cmd = nil
	return nil
}

// checkHealth verifies the wallet-rpc is responding correctly
func (w *WalletRPC) checkHealth(ctx context.Context) error {
	// TODO: Implement actual health check using RPC call
	// For now, just check if the port is still open
	if !util.IsPortInUse(w.WalletRPCPort()) {
		return errors.E(
			opCheckHealth,
			errors.ComponentWalletRPC,
			errors.KindNetwork,
			fmt.Errorf("wallet-rpc is not responding on port %d", w.WalletRPCPort()),
		)
	}
	return nil
}

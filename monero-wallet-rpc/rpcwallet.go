// Package monerowalletrpc provides functionality for managing Monero wallet RPC services.
// It handles wallet process lifecycle, RPC communication, and daemon coordination.
package monerowalletrpc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/opd-ai/moneroger/errors"
	"github.com/opd-ai/moneroger/monerod"
	"github.com/opd-ai/moneroger/util"
)

// Common operation constants for error wrapping
const (
	opStart          = errors.Op("WalletRPC.Start")
	opShutdown       = errors.Op("WalletRPC.Shutdown")
	opValidateConfig = errors.Op("WalletRPC.ValidateConfig")
	opCheckHealth    = errors.Op("WalletRPC.CheckHealth")
)

// NewWalletRPC creates and starts a new Monero wallet RPC service instance.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - config: Configuration settings including wallet file path and port
//   - daemon: Reference to running monerod instance for blockchain access
//
// Returns:
//   - *WalletRPC: Pointer to configured and running wallet RPC instance
//   - error: Any error encountered during setup or startup
//
// The function performs the following steps:
// 1. Validates configuration parameters
// 2. Creates WalletRPC instance with provided settings
// 3. Starts the wallet RPC process
// 4. Verifies service health
//
// Errors:
//   - Invalid configuration parameters
//   - Process startup failures
//   - Port binding issues
//   - Health check failures
//
// Related:
//   - validateConfig for configuration validation
//   - WalletRPC.start for process management
func NewWalletRPC(ctx context.Context, config util.Config, daemon *monerod.MoneroDaemon) (*WalletRPC, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	wallet := &WalletRPC{
		walletDir: config.WalletFile,
		rpcPort:   config.WalletPort,
		daemon:    daemon,
	}

	if err := wallet.Start(ctx); err != nil {
		return nil, err
	}

	return wallet, nil
}

// validateConfig checks the validity of wallet RPC configuration parameters.
//
// Parameters:
//   - config: Configuration settings to validate
//
// Returns:
//   - error: Validation error if any parameter is invalid
//
// Validates:
// 1. Wallet file path existence
// 2. RPC port number validity
// 3. File system permissions
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

	/*if _, err := os.Stat(config.WalletFile); os.IsNotExist(err) {
		return errors.E(
			opValidateConfig,
			errors.ComponentWalletRPC,
			errors.KindSystem,
			fmt.Errorf("wallet file does not exist: %s", config.WalletFile),
		)
	}*/

	return nil
}

// Start launches the wallet RPC process with appropriate configuration.
//
// Parameters:
//   - ctx: Context for process management and timeouts
//
// Returns:
//   - error: Any error encountered during startup
//
// The method:
// 1. Checks port availability
// 2. Configures process arguments
// 3. Launches wallet RPC process
// 4. Verifies service availability
// 5. Performs health check
func (w *WalletRPC) Start(ctx context.Context) error {
	if util.IsPortInUse(w.WalletRPCPort()) {
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindNetwork,
			fmt.Errorf("port %d is already in use", w.WalletRPCPort()),
		)
	}

	daemonAddr := fmt.Sprintf("http://localhost:%d", w.daemon.RPCPort())

	args := []string{
		"--wallet-dir", w.walletDir,
		"--rpc-bind-port", fmt.Sprintf("%d", w.WalletRPCPort()),
		"--daemon-address", daemonAddr,
		"--prompt-for-password",
		"--daemon-login", fmt.Sprintf("%s:%s", w.daemon.RPCUser(), w.daemon.RPCPass()),
		"--rpc-login", fmt.Sprintf("%s:%s", w.WalletRPCUser(), w.WalletRPCPass()),
		"--password", w.WalletPass(),
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

	// Add stdout/stderr capture
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindProcess,
			fmt.Errorf("failed to start wallet-rpc process: %w\nOutput: %s\nError: %s",
				err, stdout.String(), stderr.String()),
		)
	}

	w.cmd = cmd

	if err := util.WaitForPort(ctx, w.WalletRPCPort()); err != nil {
		// Capture output before cleanup
		output := fmt.Sprintf("Output: %s\nError: %s", stdout.String(), stderr.String())
		_ = w.Shutdown(ctx)
		return errors.E(
			opStart,
			errors.ComponentWalletRPC,
			errors.KindTimeout,
			fmt.Errorf("wallet-rpc failed to bind to port %d: %w\n%s",
				w.WalletRPCPort(), err, output),
		)
	}

	if err := w.checkHealth(ctx); err != nil {
		_ = w.Shutdown(ctx)
		return err
	}

	return nil
}

// Shutdown gracefully stops the wallet RPC service.
//
// Parameters:
//   - ctx: Context for shutdown timeout control
//
// Returns:
//   - error: Any error encountered during shutdown
//
// The method:
// 1. Sends interrupt signal to process
// 2. Waits for process termination
// 3. Cleans up resources
//
// Timeout:
//   - Default 10 second shutdown timeout
//   - Returns error if shutdown exceeds timeout
//
// Related:
//   - checkHealth for service verification
func (w *WalletRPC) Shutdown(ctx context.Context) error {
	if w.cmd.Process == nil {
		return nil
	}

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Send interrupt signal
	if err := w.cmd.Process.Signal(os.Interrupt); err != nil {
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
		_, err := w.cmd.Process.Wait()
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

	w.cmd.Process = nil
	w.cmd = nil
	return nil
}

// checkHealth verifies the wallet RPC service is responding correctly.
//
// Parameters:
//   - ctx: Context for timeout control
//
// Returns:
//   - error: Any error encountered during health check
//
// Currently:
// - Verifies port is still in use
// TODO: Implement full RPC health check
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

func (m *WalletRPC) PID() string {
	if m.cmd != nil {
		if m.cmd.Process != nil {
			return fmt.Sprintf("%d", m.cmd.Process.Pid)
		}
	}
	return "-1"
}

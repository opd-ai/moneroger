package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opd-ai/moneroger"
	"github.com/opd-ai/moneroger/util"
)

// verifyExecutables checks if required Monero executables are available
func verifyExecutables() error {
	executables := []string{"monerod", "monero-wallet-rpc"}
	for _, exe := range executables {
		_, err := exec.LookPath(exe)
		if err != nil {
			return fmt.Errorf("%s not found in PATH: %w", exe, err)
		}
	}
	return nil
}

func main() {
	// Command line flags for configuration
	var (
		dataDir    = flag.String("datadir", "", "Directory for blockchain data and wallet files")
		walletDir  = flag.String("wallet", "", "Path to wallet file (directory)")
		moneroPort = flag.Int("daemon-port", 18081, "Port for Monero daemon RPC")
		walletPort = flag.Int("wallet-port", 18083, "Port for wallet RPC")
		testnet    = flag.Bool("testnet", false, "Use testnet instead of mainnet")
		debug      = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	// Enable debug logging if requested
	if *debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	}

	// Verify Monero executables are available
	if err := verifyExecutables(); err != nil {
		log.Fatalf("Prerequisite check failed: %v", err)
	}

	// Validate command line arguments
	if *dataDir == "" {
		log.Fatal("--datadir is required")
	}
	if *walletDir == "" {
		*walletDir = *dataDir
	}

	// Convert paths to absolute
	absDataDir, err := filepath.Abs(*dataDir)
	if err != nil {
		log.Fatalf("Failed to resolve data directory path: %v", err)
	}
	absWalletFile, err := filepath.Abs(*walletDir)
	if err != nil {
		log.Fatalf("Failed to resolve wallet file path: %v", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll(absDataDir, 0o755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Create configuration
	config := util.RecommendConfig(absDataDir)
	config.WalletFile = absWalletFile
	config.MoneroPort = *moneroPort
	config.WalletPort = *walletPort
	config.TestNet = *testnet

	if *debug {
		log.Printf("Using configuration: %+v", config)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Moneroger with increased timeout for debugging
	log.Printf("Initializing Monero services (testnet: %v)...", *testnet)

	manager, err := moneroger.NewMoneroger(config)
	if err != nil {
		log.Fatalf("Failed to initialize Moneroger: %v", err)
	}
	log.Printf("Monero services initialized: monerod: %s, monero-wallet-rpc %s", manager.MoneroDaemonPID(), manager.RPCWalletPID())
	defer manager.Shutdown(ctx)

	// Handle graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-signalChan
	log.Printf("Received signal %v, initiating shutdown...", sig)

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown services
	if err := manager.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Shutdown complete")
}

Project Path: moneroger.go

I'd like you to create a tests file `tests.go` which prostvides adequate testing of the following Go programs.
The tests should be simple.
They should test only one thing at a time.
The tests should be readable.
The tests should not be excessive.

Source Tree: 
```
moneroger.go

```

`/home/user/go/src/github.com/opd-ai/moneroger/moneroger.go`:

```go
// Package moneroger provides high-level management of Monero services,
// coordinating both the Monero daemon (monerod) and wallet RPC service
// in a single unified interface.
package moneroger

import (
	"context"

	monerowalletrpc "github.com/opd-ai/moneroger/monero-wallet-rpc"
	"github.com/opd-ai/moneroger/monerod"
	"github.com/opd-ai/moneroger/util"
)

// Moneroger coordinates Monero daemon and wallet RPC services.
// It manages the lifecycle of both services ensuring proper startup
// and shutdown order.
//
// Fields:
//   - monerod: The Monero daemon instance
//   - monerowalletrpc: The wallet RPC service instance
//
// The Moneroger instance maintains references to both services
// and handles their coordination. It ensures the daemon is available
// before starting the wallet service, and handles graceful shutdown
// in the correct order.
type Moneroger struct {
	monerod         monerod.MoneroDaemon
	monerowalletrpc monerowalletrpc.WalletRPC
}

// NewMoneroger creates a new instance managing both Monero services.
//
// Parameters:
//   - config: Configuration settings for both services including:
//     DataDir: Base directory for blockchain and wallet data
//     WalletFile: Path to wallet file
//     MoneroPort: Daemon RPC port
//     WalletPort: Wallet RPC port
//     TestNet: Network selection flag
//
// Returns:
//   - *Moneroger: Configured manager instance
//   - error: Any error during setup
//
// The function:
// 1. Starts the Monero daemon
// 2. Starts the wallet RPC service
// 3. Returns a manager coordinating both services
//
// Errors:
//   - Daemon startup failures
//   - Wallet service startup failures
//   - Configuration validation errors
//
// Related:
//   - monerod.NewMoneroDaemon
//   - monerowalletrpc.NewWalletRPC
//   - util.Config
func NewMoneroger(config util.Config) (*Moneroger, error) {
	ctx := context.Background()

	// Start Monero daemon
	daemon, err := monerod.NewMoneroDaemon(ctx, config)
	if err != nil {
		return nil, err
	}

	// Start wallet RPC service
	wallet, err := monerowalletrpc.NewWalletRPC(ctx, config, daemon)
	if err != nil {
		return nil, err
	}

	return &Moneroger{
		monerod:         *daemon,
		monerowalletrpc: *wallet,
	}, nil
}

// start initializes both Monero services in the correct order.
// This is an internal method used by NewMoneroger.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error during startup sequence
//
// The method:
// 1. Starts the Monero daemon
// 2. Waits for daemon availability
// 3. Starts the wallet RPC service
//
// Related:
//   - MoneroDaemon.Start
//   - WalletRPC.Start
func (m *Moneroger) Start(ctx context.Context) error {
	if err := m.monerod.Start(ctx); err != nil {
		return err
	}
	return m.monerowalletrpc.Start(ctx)
}

// Shutdown gracefully stops both Monero services in the correct order.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error during shutdown sequence
//
// The method:
// 1. Stops the wallet RPC service first
// 2. Stops the Monero daemon after
//
// This order ensures proper cleanup and prevents wallet
// errors due to daemon unavailability.
//
// Related:
//   - WalletRPC.Shutdown
//   - MoneroDaemon.Shutdown
func (m *Moneroger) Shutdown(ctx context.Context) error {
	if err := m.monerowalletrpc.Shutdown(ctx); err != nil {
		return err
	}
	return m.monerod.Shutdown(ctx)
}

```  


Rely on standard test packages.
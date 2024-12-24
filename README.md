# Moneroger

Moneroger is a Go library that provides robust process management and coordination for Monero daemons (monerod) and wallet RPC services. It handles process lifecycle, configuration, and health monitoring with proper error handling and graceful shutdown support.

## Features

- 🚀 **Automated Process Management**
  - Start/stop monerod and monero-wallet-rpc processes
  - Automatic executable discovery in system PATH
  - Health monitoring and port availability checks
  - Graceful shutdown handling

- 🔒 **Security First**
  - Automatic secure RPC credential generation
  - Proper authentication between components
  - Safe process handling and cleanup

- ⚙️ **Flexible Configuration**
  - Support for both mainnet and testnet
  - Configurable data directories and ports
  - Timeout controls for operations
  - Custom RPC credentials

- 🛠️ **Developer Friendly**
  - Structured error handling with context
  - Clear component separation
  - Comprehensive testing
  - Well-documented API

## Installation

```bash
go get github.com/opd-ai/moneroger
```

Requires Go 1.21 or later.

## Prerequisites

- Monero daemon (`monerod`) installed and in system PATH
- Monero wallet RPC (`monero-wallet-rpc`) installed and in system PATH
- Write permissions for data directory

## Usage

### Basic Example

```go
package main

import (
    "context"
    "log"

    "github.com/opd-ai/moneroger/monerod"
    "github.com/opd-ai/moneroger/monerowalletrpc"
    "github.com/opd-ai/moneroger/util"
)

func main() {
    ctx := context.Background()
    
    // Configure services
    config := util.Config{
        DataDir:    "/path/to/monero/data",
        WalletFile: "/path/to/wallet.keys",
        MoneroPort: 18081,
        WalletPort: 18082,
        TestNet:    false,
    }

    // Start Monero daemon
    daemon, err := monerod.NewMoneroDaemon(ctx, config)
    if err != nil {
        log.Fatal(err)
    }
    defer daemon.Shutdown(ctx)

    // Start wallet RPC service
    wallet, err := monerowalletrpc.NewWalletRPC(ctx, config, daemon)
    if err != nil {
        log.Fatal(err)
    }
    defer wallet.Shutdown(ctx)

    // Services are now running...
}
```

### Configuration Options

```go
type Config struct {
    // Base directory for blockchain data and wallet files
    DataDir string

    // Path to the Monero wallet file (.keys)
    WalletFile string

    // TCP port for monerod RPC service (default: 18081)
    MoneroPort int

    // TCP port for wallet RPC service (default: 18082)
    WalletPort int

    // Run on testnet instead of mainnet
    TestNet bool
}
```

## Error Handling

The library provides structured error handling with categorized errors:

```go
switch errors.GetKind(err) {
case errors.KindNetwork:
    // Handle network-related errors
case errors.KindProcess:
    // Handle process management errors
case errors.KindConfig:
    // Handle configuration errors
case errors.KindTimeout:
    // Handle timeout errors
case errors.KindSystem:
    // Handle system-level errors
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

For verbose output:

```bash
go test -v ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests and ensure they pass (`go test ./...`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development Requirements

- Go 1.21+
- `gofumpt` for code formatting
- Access to Monero executables for integration testing

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Monero Project](https://www.getmonero.org/) for the core Monero software
- [go-password](https://github.com/sethvargo/go-password) for secure password generation

## Status

This project is under active development. API may change before reaching v1.0.0.
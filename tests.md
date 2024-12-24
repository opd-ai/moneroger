Project Path: moneroger

I'd like you to create a tests file `tests.go` which prostvides adequate testing of the following Go programs.
The tests should be simple.
They should test only one thing at a time.
The tests should be readable.
The tests should not be excessive.

Source Tree: 
```
moneroger
├── util
│   ├── util.go
│   └── config.go
├── LICENSE
├── errors
│   ├── types.go
│   └── errors.go
├── go.mod
├── README.md
├── const
│   └── const.go
├── Makefile
├── monerod
│   ├── monerod.go
│   └── types.go
├── monero-wallet-rpc
│   ├── rpcwallet.go
│   └── types.go
└── go.sum

```

`/home/user/go/src/github.com/opd-ai/moneroger/util/util.go`:

```go
// Package util provides utility functions for the moneroger library,
// including file operations, path management, password generation,
// and network port handling.
package util

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	moneroconst "github.com/opd-ai/moneroger/const"
	"github.com/sethvargo/go-password/password"
)

// FileExists checks if a file exists at the specified path.
//
// Parameters:
//   - path: The file path to check (string)
//
// Returns:
//   - bool: true if the file exists and is accessible, false otherwise
//
// Note: This function returns false for both non-existent files
// and files that exist but are inaccessible due to permissions.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// Path returns a slice of directories to search for executables.
// It combines:
// - The directory containing the current executable
// - The current working directory
// - The system PATH environment variable
//
// Returns:
//   - []string: Slice of directory paths to search
//
// Note: Logs errors but continues execution if unable to determine
// executable or working directory paths.
func Path() []string {
	path := os.Getenv("PATH")
	elements := []string{}

	// Get executable directory
	me, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path:", err)
	} else {
		meDir := filepath.Dir(me)
		elements = append(elements, meDir)
	}

	// Get working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory:", err)
	} else {
		elements = append(elements, workDir)
	}

	// Add system PATH elements
	if path != "" {
		elements = append(elements, strings.Split(path, ":")...)
	}

	return elements
}

// SecurePassword generates a cryptographically secure random password.
//
// Returns:
//   - string: A 20-character password with random number of digits (0-19)
//
// Panics:
//   - If the password generation fails (should be extremely rare)
//
// Uses:
//   - github.com/sethvargo/go-password/password for generation
func SecurePassword() string {
	rand.Seed(time.Now().UnixNano())
	digs := rand.Intn(19)
	res, err := password.Generate(20, digs, 0, false, false)
	if err != nil {
		panic(err)
	}
	return res
}

// IsPortInUse checks if a TCP port is currently in use on localhost.
//
// Parameters:
//   - port: Port number to check (int)
//
// Returns:
//   - bool: true if the port is in use, false otherwise
//
// Note: This function attempts a TCP connection with a 1-second timeout.
// A successful connection indicates the port is in use.
func IsPortInUse(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// WaitForPort waits for a TCP port to become available.
//
// Parameters:
//   - ctx: Context for cancellation
//   - port: Port number to wait for (int)
//
// Returns:
//   - error: nil if port becomes available, error otherwise
//
// Errors:
//   - Context cancellation error if context is cancelled
//   - Timeout error if port doesn't become available within DefaultStartupTimeout
//
// Related:
//   - moneroconst.DefaultStartupTimeout
//   - IsPortInUse function
func WaitForPort(ctx context.Context, port int) error {
	deadline := time.Now().Add(moneroconst.DefaultStartupTimeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if IsPortInUse(port) {
				return nil
			}
			time.Sleep(time.Second)
		}
	}
	return fmt.Errorf("timeout waiting for port %d", port)
}

```  

`/home/user/go/src/github.com/opd-ai/moneroger/util/config.go`:

```go
package util

// Config holds the configuration parameters for both monerod and monero-wallet-rpc daemons.
// It provides all necessary settings for initializing and running the Monero services.
//
// Fields:
//
//   - DataDir: Base directory for blockchain data and wallet files
//     Must be writable by the process
//
//   - WalletFile: Path to the Monero wallet file (.keys file)
//     Can be absolute or relative to DataDir
//
//   - MoneroPort: TCP port for monerod RPC service
//     Default: 18081 (mainnet), 28081 (testnet)
//     Must be available and accessible
//
//   - WalletPort: TCP port for monero-wallet-rpc service
//     Default: 18082 (mainnet), 28082 (testnet)
//     Must be available and accessible
//
//   - TestNet: Flag to run services on Monero testnet
//     true = testnet, false = mainnet
//
// Usage:
//
//	config := &Config{
//	    DataDir:    "/path/to/monero/data",
//	    WalletFile: "wallet.keys",
//	    MoneroPort: 18081,
//	    WalletPort: 18082,
//	    TestNet:    false,
//	}
//
// Related:
//   - moneroconst.DefaultMonerodPort
//   - moneroconst.DefaultWalletRPCPort
//   - util.IsPortInUse() for port validation
//   - util.FileExists() for path validation
type Config struct {
	// DataDir is the base directory for blockchain data and wallet files
	DataDir string

	// WalletFile is the path to the Monero wallet file
	WalletFile string

	// MoneroPort is the TCP port for monerod RPC service
	MoneroPort int

	// WalletPort is the TCP port for monero-wallet-rpc service
	WalletPort int

	// TestNet determines whether to run on testnet (true) or mainnet (false)
	TestNet bool
}

```  

`/home/user/go/src/github.com/opd-ai/moneroger/LICENSE`:

```
MIT License

Copyright (c) 2024 opdai

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```  

`/home/user/go/src/github.com/opd-ai/moneroger/errors/types.go`:

```go
// Package errors provides a structured error handling system for the moneroger library.
// It enables detailed error tracking with operation context, error categorization,
// and proper error wrapping support following Go 1.13+ error conventions.
package errors

import "fmt"

// Op represents an operation in the moneroger library.
// It is typically the name of the method where the error occurred.
// Example: "WalletRPC.Start" or "MoneroDaemon.Shutdown"
type Op string

// Kind represents the category of error that occurred.
// This enables programmatic error handling and appropriate error responses.
type Kind uint8

const (
	// KindUnknown represents errors that don't fall into other categories
	KindUnknown Kind = iota

	// KindNetwork represents network-related errors such as:
	// - Port binding failures
	// - Connection timeouts
	// - RPC communication issues
	KindNetwork

	// KindProcess represents process management errors such as:
	// - Process start failures
	// - Unexpected process termination
	// - Signal handling issues
	KindProcess

	// KindConfig represents configuration-related errors such as:
	// - Invalid port numbers
	// - Missing required paths
	// - Invalid parameter combinations
	KindConfig

	// KindTimeout represents timeout-related errors such as:
	// - Daemon startup timeouts
	// - Operation deadlines exceeded
	// - Response waiting timeouts
	KindTimeout

	// KindSystem represents system-level errors such as:
	// - File permission issues
	// - Resource exhaustion
	// - System call failures
	KindSystem
)

// Error is the fundamental error type for the moneroger library.
// It implements the error interface and provides detailed context about where
// and how an error occurred.
//
// Fields:
//   - Op: The operation where the error occurred
//   - Component: The system component (e.g., "monerod", "wallet-rpc")
//   - Kind: The category of error
//   - Err: The underlying error (if any)
//
// Usage:
//
//	return &Error{
//	    Op: "WalletRPC.Start",
//	    Component: "wallet-rpc",
//	    Kind: KindNetwork,
//	    Err: err,
//	}
type Error struct {
	// Op identifies the operation that failed
	Op Op

	// Component identifies the system component where the error occurred
	Component string

	// Kind identifies the category of error
	Kind Kind

	// Err is the underlying error that triggered this error, if any
	Err error
}

// Error implements the error interface, providing a formatted error message.
// The message format is: "component: operation: kind: underlying error"
// If there is no underlying error, it omits the last part.
//
// Returns:
//   - A formatted string containing error details
func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s: %s", e.Component, e.Op, e.Kind)
	}
	return fmt.Sprintf("%s: %s: %s: %v", e.Component, e.Op, e.Kind, e.Err)
}

// Unwrap implements the errors unwrapping interface from Go 1.13+.
// This allows the error to work with errors.Is and errors.As.
//
// Returns:
//   - The underlying error, or nil if none exists
//
// Related: https://golang.org/pkg/errors/#Unwrap
func (e *Error) Unwrap() error {
	return e.Err
}

// String converts a Kind to its string representation.
// This method is used for error message formatting and logging.
//
// Returns:
//   - A human-readable string describing the error kind
//
// Example:
//
//	KindNetwork.String() returns "network error"
func (k Kind) String() string {
	switch k {
	case KindNetwork:
		return "network error"
	case KindProcess:
		return "process error"
	case KindConfig:
		return "configuration error"
	case KindTimeout:
		return "timeout error"
	case KindSystem:
		return "system error"
	default:
		return "unknown error"
	}
}

```  

`/home/user/go/src/github.com/opd-ai/moneroger/errors/errors.go`:

```go
package errors

import "errors"

// Component names identify the major components of the moneroger system.
// These constants are used to provide context in error messages and logging.
const (
	// ComponentMonerod identifies the Monero daemon component
	ComponentMonerod = "monerod"

	// ComponentWalletRPC identifies the Monero wallet RPC component
	ComponentWalletRPC = "wallet-rpc"

	// ComponentUtil identifies the utility functions component
	ComponentUtil = "util"
)

// Common operations represent standard actions performed across components.
// These constants are used to identify where in the operation lifecycle an error occurred.
const (
	// OpStart represents daemon or service startup operations
	OpStart Op = "Start"

	// OpShutdown represents graceful shutdown operations
	OpShutdown Op = "Shutdown"

	// OpHealthCheck represents health monitoring operations
	OpHealthCheck Op = "HealthCheck"

	// OpPortBinding represents network port binding operations
	OpPortBinding Op = "PortBinding"

	// OpProcessSpawn represents process creation operations
	OpProcessSpawn Op = "ProcessSpawn"
)

// E creates a new Error from a variable number of arguments.
// It constructs an Error by examining the type of each argument and setting
// the corresponding field in the Error struct.
//
// Parameters:
//   - args: Variable number of arguments that can be of types:
//   - Op: Sets the operation field
//   - string: Sets the component field
//   - Kind: Sets the error kind
//   - *Error: Copies the error as the underlying error
//   - error: Sets as the underlying error
//
// Returns:
//   - error: A new Error instance with fields set based on the provided arguments
//
// Example:
//
//	E(OpStart, ComponentMonerod, KindNetwork, err)
func E(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch a := arg.(type) {
		case Op:
			e.Op = a
		case string:
			e.Component = a
		case Kind:
			e.Kind = a
		case *Error:
			copy := *a
			e.Err = &copy
		case error:
			e.Err = a
		}
	}
	return e
}

// Is reports whether two errors are of the same type and kind.
// This function implements part of the Go 1.13+ error comparison interface.
//
// Parameters:
//   - err: The error being checked
//   - target: The error to compare against
//
// Returns:
//   - bool: true if both errors are Error types and have the same Kind
//
// Example:
//
//	Is(err, &Error{Kind: KindNetwork}) // checks if err is a network error
//
// Related: https://golang.org/pkg/errors/#Is
func Is(err, target error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Kind == t.Kind
}

// GetKind extracts the Kind from an error.
// It safely handles both Error types and regular errors.
//
// Parameters:
//   - err: The error to examine, can be any error type
//
// Returns:
//   - Kind: The kind of error, or KindUnknown if not an Error type
//
// Example:
//
//	kind := GetKind(err)
//	if kind == KindNetwork {
//	    // handle network error
//	}
//
// Related:
//   - errors.As: https://golang.org/pkg/errors/#As
//   - Kind type in types.go
func GetKind(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return KindUnknown
}

```  

`/home/user/go/src/github.com/opd-ai/moneroger/go.mod`:

```mod
module github.com/opd-ai/moneroger

go 1.21.3

require github.com/sethvargo/go-password v0.3.1

```  

`/home/user/go/src/github.com/opd-ai/moneroger/README.md`:

```md
# moneroger
Manages monerod and monero-wallet-rpc for go applications

```  

`/home/user/go/src/github.com/opd-ai/moneroger/const/const.go`:

```go
// Package moneroconst provides default configuration constants for the moneroger library.
// It defines standard ports and timeouts used by both monerod and monero-wallet-rpc daemons.
package moneroconst

import (
	"time"
)

// Default configurations for Monero daemons
const (
	// DefaultMonerodPort is the standard RPC port for monerod daemon (18081)
	// This port is used for communication between the daemon and wallet
	DefaultMonerodPort = 18081

	// DefaultWalletRPCPort is the standard RPC port for monero-wallet-rpc daemon (18082)
	// This port is used by applications to communicate with the wallet
	DefaultWalletRPCPort = 18082

	// DefaultStartupTimeout defines how long to wait for daemons to start (30 seconds)
	// If a daemon doesn't respond within this time, startup is considered failed
	DefaultStartupTimeout = 30 * time.Second

	// DefaultShutdownTimeout defines how long to wait for graceful shutdown (10 seconds)
	// After this timeout, the process will be forcefully terminated
	DefaultShutdownTimeout = 10 * time.Second
)

```  

`/home/user/go/src/github.com/opd-ai/moneroger/Makefile`:

```

fmt:
	find . -name '*.go' -exec gofumpt -w -s -extra {} \;

doc:/
	find . -type d -exec code2prompt --template ~/code2prompt/templates/write-a-test.hbs --output {}/tests.md {} \;
```  

`/home/user/go/src/github.com/opd-ai/moneroger/monerod/monerod.go`:

```go
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

// start launches the monerod process with appropriate configuration.
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

```  

`/home/user/go/src/github.com/opd-ai/moneroger/monerod/types.go`:

```go
// Package monerod provides functionality for managing Monero daemon processes.
// It handles daemon configuration, lifecycle management, and RPC communication.
package monerod

import (
	"os"
	"os/exec"
	"time"

	"github.com/opd-ai/moneroger/util"
)

// Default configurations for the Monero daemon
const (
	// defaultMonerodPort is the default RPC port for monerod (mainnet)
	defaultMonerodPort = 18081

	// defaultStartupTimeout is the maximum time to wait for daemon startup
	defaultStartupTimeout = 30 * time.Second

	// defaultShutdownTimeout is the maximum time to wait for graceful shutdown
	defaultShutdownTimeout = 10 * time.Second
)

// MoneroDaemon represents a running monerod instance and manages its lifecycle.
// It provides access to daemon configuration and process control.
//
// Fields:
//   - cmd: Command instance for process management
//   - dataDir: Directory for blockchain and configuration data
//   - rpcPort: Port number for RPC interface
//   - rpcUser: Username for RPC authentication
//   - rpcPass: Password for RPC authentication
//   - testnet: Boolean flag for testnet operation
//   - process: Reference to the running daemon process
//
// The daemon can be configured for either mainnet or testnet operation,
// with appropriate default ports and network settings applied automatically.
type MoneroDaemon struct {
	cmd     *exec.Cmd
	dataDir string
	rpcPort int
	rpcUser string
	rpcPass string
	testnet bool
	process *os.Process
}

// RPCPort returns the configured RPC port for the daemon.
// If no port was explicitly set, returns the default port (18081 for mainnet).
//
// Returns:
//   - int: The RPC port number
//
// Related:
//   - defaultMonerodPort constant
func (m *MoneroDaemon) RPCPort() int {
	return m.rpcPort
}

// RPCUser returns the RPC authentication username.
// If no username was set, initializes it to the default "gouser".
//
// Returns:
//   - string: The RPC username
//
// Note: This method has side effects - it will set a default username
// if one hasn't been configured.
func (m *MoneroDaemon) RPCUser() string {
	if m.rpcUser == "" {
		m.rpcUser = "gouser"
	}
	return m.rpcUser
}

// RPCPass returns the RPC authentication password.
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
func (m *MoneroDaemon) RPCPass() string {
	if m.rpcPass == "" {
		m.rpcPass = util.SecurePassword()
	}
	return m.rpcPass
}

```  

`/home/user/go/src/github.com/opd-ai/moneroger/monero-wallet-rpc/rpcwallet.go`:

```go
// Package monerowalletrpc provides functionality for managing Monero wallet RPC services.
// It handles wallet process lifecycle, RPC communication, and daemon coordination.
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
		walletFile: config.WalletFile,
		rpcPort:    config.WalletPort,
		daemon:     daemon,
	}

	if err := wallet.start(ctx); err != nil {
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

// start launches the wallet RPC process with appropriate configuration.
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

```  

`/home/user/go/src/github.com/opd-ai/moneroger/monero-wallet-rpc/types.go`:

```go
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
	walletFile string
	rpcPort    int
	rpcUser    string
	rpcPass    string
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

```  

`/home/user/go/src/github.com/opd-ai/moneroger/go.sum`:

```sum
github.com/sethvargo/go-password v0.3.1 h1:WqrLTjo7X6AcVYfC6R7GtSyuUQR9hGyAj/f1PYQZCJU=
github.com/sethvargo/go-password v0.3.1/go.mod h1:rXofC1zT54N7R8K/h1WDUdkf9BOx5OptoxrMBcrXzvs=

```  


Rely on standard test packages.
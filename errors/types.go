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

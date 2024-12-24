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

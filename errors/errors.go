package errors

import "errors"

// Component names
const (
	ComponentMonerod   = "monerod"
	ComponentWalletRPC = "wallet-rpc"
	ComponentUtil      = "util"
)

// Common operations
const (
	OpStart        Op = "Start"
	OpShutdown     Op = "Shutdown"
	OpHealthCheck  Op = "HealthCheck"
	OpPortBinding  Op = "PortBinding"
	OpProcessSpawn Op = "ProcessSpawn"
)

// E creates a new Error
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

// Is reports whether target is of the same type and kind as err
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

// GetKind extracts the Kind from an error
func GetKind(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return KindUnknown
}

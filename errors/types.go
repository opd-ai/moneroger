package errors

import "fmt"

// Op represents an operation in the moneroger library
type Op string

// Kind represents the category of error
type Kind uint8

const (
	// Error kinds
	KindUnknown Kind = iota
	KindNetwork      // Network-related errors
	KindProcess      // Process management errors
	KindConfig       // Configuration errors
	KindTimeout      // Timeout errors
	KindSystem       // System-level errors
)

// Error is the fundamental error type for the moneroger library
type Error struct {
	// The operation being performed, usually the method name
	Op Op

	// The component where the error occurred
	Component string

	// The specific kind of error
	Kind Kind

	// The underlying error
	Err error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s: %s", e.Component, e.Op, e.Kind)
	}
	return fmt.Sprintf("%s: %s: %s: %v", e.Component, e.Op, e.Kind, e.Err)
}

// Unwrap implements the errors unwrapping interface
func (e *Error) Unwrap() error {
	return e.Err
}

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

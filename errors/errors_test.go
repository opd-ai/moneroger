package errors

import (
	"fmt"
	"testing"
)

// TestErrorString verifies the Error.Error() string formatting
func TestErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "complete error",
			err: &Error{
				Op:        OpStart,
				Component: ComponentMonerod,
				Kind:      KindNetwork,
				Err:       fmt.Errorf("connection failed"),
			},
			expected: "monerod: Start: network error: connection failed",
		},
		{
			name: "error without underlying error",
			err: &Error{
				Op:        OpShutdown,
				Component: ComponentWalletRPC,
				Kind:      KindProcess,
			},
			expected: "wallet-rpc: Shutdown: process error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestE verifies the error construction helper function
func TestE(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	err := E(OpStart, ComponentMonerod, KindNetwork, baseErr)

	e, ok := err.(*Error)
	if !ok {
		t.Fatal("E() should return *Error")
	}

	if e.Op != OpStart {
		t.Errorf("Op = %v, want %v", e.Op, OpStart)
	}
	if e.Component != ComponentMonerod {
		t.Errorf("Component = %v, want %v", e.Component, ComponentMonerod)
	}
	if e.Kind != KindNetwork {
		t.Errorf("Kind = %v, want %v", e.Kind, KindNetwork)
	}
	if e.Err != baseErr {
		t.Errorf("Err = %v, want %v", e.Err, baseErr)
	}
}

// TestIs verifies error type comparison
func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "matching kinds",
			err:      &Error{Kind: KindNetwork},
			target:   &Error{Kind: KindNetwork},
			expected: true,
		},
		{
			name:     "different kinds",
			err:      &Error{Kind: KindNetwork},
			target:   &Error{Kind: KindProcess},
			expected: false,
		},
		{
			name:     "non-Error type",
			err:      fmt.Errorf("regular error"),
			target:   &Error{Kind: KindNetwork},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.target); got != tt.expected {
				t.Errorf("Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestGetKind verifies error kind extraction
func TestGetKind(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected Kind
	}{
		{
			name:     "network error",
			err:      &Error{Kind: KindNetwork},
			expected: KindNetwork,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			expected: KindUnknown,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: KindUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKind(tt.err); got != tt.expected {
				t.Errorf("GetKind() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestUnwrap verifies error unwrapping
func TestUnwrap(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	err := &Error{
		Kind: KindNetwork,
		Err:  baseErr,
	}

	if got := err.Unwrap(); got != baseErr {
		t.Errorf("Unwrap() = %v, want %v", got, baseErr)
	}
}

// TestKindString verifies Kind.String() conversions
func TestKindString(t *testing.T) {
	tests := []struct {
		kind     Kind
		expected string
	}{
		{KindNetwork, "network error"},
		{KindProcess, "process error"},
		{KindConfig, "configuration error"},
		{KindTimeout, "timeout error"},
		{KindSystem, "system error"},
		{KindUnknown, "unknown error"},
		{Kind(99), "unknown error"}, // Invalid kind
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

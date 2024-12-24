package moneroconst

import (
	"time"
)

// Default configurations
const (
	DefaultMonerodPort     = 18081
	DefaultWalletRPCPort   = 18082
	DefaultStartupTimeout  = 30 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
)

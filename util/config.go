package util

import (
	"log"
	"math"
	"path/filepath"

	"github.com/ricochet2200/go-disk-usage/du"
)

var TwoHundredFiftyGigabytes uint64 = uint64(250 * math.Pow(10, 9))

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
//		config := &Config{
//		    DataDir:    "/path/to/monero/data",
//		    WalletFile: "wallet.keys",
//		    MoneroPort: 18081,
//		    WalletPort: 18082,
//		    TestNet:    false,
//	        RemoteNode: "",
//		}
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
	// RemoteNode instructs the monero-wallet-rpc client to use a remote port
	RemoteNode string
}

func RecommendConfig(dataDir string) (config Config) {
	config.DataDir = dataDir
	if !DirExists(config.DataDir) {
		usage := du.NewDiskUsage(config.DataDir)
		if usage.Available() > TwoHundredFiftyGigabytes {
			log.Println("Greater than 250GB available space detected, full node functionality enabled")
		}
		config.RemoteNode = ""
	} else {
		config.RemoteNode = pickDefaultRemoteNode()
	}
	config.TestNet = false
	config.WalletFile = filepath.Join(config.DataDir, "wallet")
	config.MoneroPort = 18081
	config.WalletPort = 18083
	return
}

func pickDefaultRemoteNode() string {
	return "not enabled yet"
}

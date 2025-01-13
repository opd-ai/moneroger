package util

import (
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/spf13/viper"
)

var TwoHundredFiftyGigabytes uint64 = uint64(250 * math.Pow(10, 9))

func pickDefaultRemoteNode() string {
	return "not enabled yet"
}

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

// RecommendConfig generates a recommended Monero configuration based on the provided data directory.
// If no data directory is specified, it creates one in the current working directory under "moneroger".
// It also checks available disk space to determine if full node functionality should be enabled.
//
// Parameters:
//   - dataDir: String path to desired data directory. If empty, defaults to ./moneroger
//
// Returns:
//   - Config: A Config struct with recommended settings:
//   - DataDir: Absolute path to data directory
//   - WalletFile: Set to "wallet" in the data directory
//   - MoneroPort: Default 18081
//   - WalletPort: Default 18083
//   - TestNet: Set to false (mainnet)
//   - RemoteNode: Empty string if enough disk space (>250GB), otherwise a remote node address
//
// Panics:
//   - If unable to get current working directory
//
// Related:
//   - TwoHundredFiftyGigabytes constant for space requirement
//   - DirExists() for directory validation
//   - pickDefaultRemoteNode() for remote node selection
func RecommendConfig(dataDir string) (config Config) {
	if dataDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dataDir = filepath.Join(wd, "moneroger")
	}
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

// LoadConfig reads and parses the configuration file at the specified path.
// It returns the parsed configuration and any error encountered.
//
// Parameters:
//   - path: File path to the YAML configuration file
//
// Returns:
//   - *Config: Parsed configuration structure
//   - error: Any error encountered during loading or parsing
//
// The function will return an error if:
//   - The configuration file cannot be read
//   - The YAML is invalid
//   - Required fields are missing
//
// Related types:
//   - Config: The configuration structure
//   - viper.Viper: Underlying configuration parser
func LoadConfig(path string) (*Config, error) {
	// Set the configuration file path and type
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Parse into Config structure
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

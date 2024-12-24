package util

// Config holds the configuration for both daemons
type Config struct {
	DataDir    string
	WalletFile string
	MoneroPort int
	WalletPort int
	TestNet    bool
}

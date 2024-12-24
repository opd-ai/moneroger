package util

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	moneroconst "github.com/opd-ai/moneroger/const"
	"github.com/sethvargo/go-password/password"
)

func SecurePassword() string {
	rand.Seed(time.Now().UnixNano())
	digs := rand.Intn(19)
	res, err := password.Generate(20, digs, 0, false, false)
	if err != nil {
		panic(err)
	}
	return res
}

// Helper function to check if a port is in use
func IsPortInUse(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Helper function to wait for a port to become available
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

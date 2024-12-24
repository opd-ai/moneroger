package util

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	moneroconst "github.com/opd-ai/moneroger/const"
	"github.com/sethvargo/go-password/password"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func Path() []string {
	path := os.Getenv("PATH")
	if path == "" {

	}
	elements := []string{}
	me, err := os.Executable()
	if err != nil {
		log.Println("this should probably be impossible but OK", err)
	}
	meDir := filepath.Dir(me)
	elements = append(elements, meDir)
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("this should also be impossible, but OK", err)
	}
	elements = append(elements, workDir)
	elements = append(elements, strings.Split(path, ":")...)
	return elements
}

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

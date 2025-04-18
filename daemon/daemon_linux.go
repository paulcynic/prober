//go:build linux
// +build linux


package daemon

import (
	"os"
	"path/filepath"
	"strconv"
)

func processExists(pid int) bool {
	_, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid)))
	return err == nil // err is nil if file exists
}

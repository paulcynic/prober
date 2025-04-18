//go:build !linux && !windows
// +build !linux,!windows


package daemon

import "golang.org/x/sys/unix"

func processExists(pid int) bool {
	// OS X & BSD systems don't have a proc filesystem.
	// Use kill -0 pid to judge if the process exists.
	err := unix.Kill(pid, 0)
	return err == nil
}

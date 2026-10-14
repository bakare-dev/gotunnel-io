//go:build !windows

package main

import (
	"os"
	"syscall"
)

func run(exe string, args []string) {
	err := syscall.Exec(exe, append([]string{exe}, args...), os.Environ())
	fail("installed GoTunnel but could not start it: %v", err)
}

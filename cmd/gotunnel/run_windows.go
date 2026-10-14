//go:build windows

package main

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
)

func run(exe string, args []string) {
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	cmd := exec.Command(exe, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}
	if err != nil {
		fail("installed GoTunnel but could not start it: %v", err)
	}
}

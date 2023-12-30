package main

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

var binaryName = "shutdown"
var sudo = "sudo"

type ShutdownHandler interface {
	// ShutDownDelayed initiates a shutdown in the given minutes.
	// If 0, the shutdown is triggered immediately and cannot be canceled anymore.
	ShutDownDelayed(delayMinutes int) error
	// CancelShutdown cancels a pending shutdown, provided the shutdown is still cancelable.
	CancelShutdown() error
}

type DryRunShutdown struct {
	Logger *slog.Logger
}

func (d *DryRunShutdown) ShutDownDelayed(delayMinutes int) error {
	d.Logger.Info("Dry-Run: Shutting down", "args", strings.Join(getDelayedShutdownArgs(false, delayMinutes), ", "))
	return nil
}

func (d *DryRunShutdown) CancelShutdown() error {
	d.Logger.Info("Dry-Run: Canceling shutdown", "args", strings.Join(getCancelShutdownArgs(false), ", "))
	return nil
}

type ExecutableShutdown struct {
	Logger   *slog.Logger
	SkipSudo bool
}

func (e *ExecutableShutdown) ShutDownDelayed(delayMinutes int) error {
	cmd := exec.Command(e.getBinary(), getDelayedShutdownArgs(e.SkipSudo, delayMinutes)...)
	return e.run(cmd)
}

func (e *ExecutableShutdown) CancelShutdown() error {
	cmd := exec.Command(e.getBinary(), getCancelShutdownArgs(e.SkipSudo)...)
	return e.run(cmd)
}

func (e *ExecutableShutdown) run(cmd *exec.Cmd) error {
	e.Logger.Debug("Running command", "path", cmd.Path, "args", strings.Join(cmd.Args, ", "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not run command: %w: %s", err, out)
	}
	e.Logger.Debug("Ran command", "output", string(out))
	return nil
}

func getDelayedShutdownArgs(skipSudo bool, delayMinutes int) []string {
	args := []string{"--poweroff", fmt.Sprintf("+%d", delayMinutes), "Shutdown initiated over Web-Powercycle service"}
	if skipSudo {
		return args
	}
	return append([]string{binaryName}, args...)
}

func getCancelShutdownArgs(skipSudo bool) []string {
	args := []string{"-c"}
	if skipSudo {
		return args
	}
	return append([]string{binaryName}, args...)
}

func (e *ExecutableShutdown) getBinary() string {
	if e.SkipSudo {
		return binaryName
	}
	return sudo
}

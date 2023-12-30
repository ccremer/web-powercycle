package main

import (
	"github.com/urfave/cli/v2"
)

func newLogLevelFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name: "log-level", Aliases: []string{"v"}, EnvVars: []string{"LOG_LEVEL"},
		Usage:       "logging verbosity",
		DefaultText: "info [debug | warn | error | disabled]",
		Value:       "info",
	}
}

func newDryRunFlag(dest *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "dry-run",
		Usage:       "Don't actually shut down system, only log. Useful for testing",
		Destination: dest,
	}
}

func newSkipSudoFlag(dest *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "skip-sudo",
		Usage:       "Don't invoke 'shutdown' with sudo",
		Hidden:      true,
		Destination: dest,
	}
}

func newAuthUserFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "auth-user",
		Usage:       "User name for basic auth",
		Required:    true,
		EnvVars:     []string{"AUTH_USER"},
		Destination: dest,
	}
}

func newAuthPassFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "auth-pass",
		Usage:       "Password for basic auth",
		Required:    true,
		EnvVars:     []string{"AUTH_PASS"},
		Destination: dest,
	}
}

package main

import (
	"github.com/urfave/cli/v2"
)

func newLogLevelFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name: "log-level", Aliases: []string{"v"}, EnvVars: []string{"LOG_LEVEL"},
		Usage:       "logging verbosity",
		DefaultText: "info [debug | warn | error]",
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

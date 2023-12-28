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

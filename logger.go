package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var logger *slog.Logger

func init() {
	// Remove `-v` short option from --version flag
	cli.VersionFlag.(*cli.BoolFlag).Aliases = nil
}

// LogMetadata prints various metadata to the root logger.
// It prints version, architecture and current user ID and returns nil.
func LogMetadata(c *cli.Context) error {
	logger.Info("Starting up "+appName,
		"version", version,
		"date", date,
		"commit", commit,
		"go_os", runtime.GOOS,
		"go_arch", runtime.GOARCH,
		"uid", os.Getuid(),
		"gid", os.Getgid(),
	)
	return nil
}

func setupLogging(c *cli.Context) error {
	level := c.String(newLogLevelFlag().Name)

	logLevel := pterm.LogLevelInfo
	switch level {
	case "debug":
		logLevel = pterm.LogLevelDebug
	case "warn":
		logLevel = pterm.LogLevelWarn
	case "error":
		logLevel = pterm.LogLevelError
	}

	backend := pterm.DefaultLogger
	logHandler := pterm.NewSlogHandler(backend.WithLevel(logLevel))

	logger = slog.New(logHandler)
	return nil
}

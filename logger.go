package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var slogger *slog.Logger
var stdLogger *log.Logger

func init() {
	// Remove `-v` short option from --version flag
	cli.VersionFlag.(*cli.BoolFlag).Aliases = nil
}

// LogMetadata prints various metadata to the root slogger.
// It prints version, architecture and current user ID and returns nil.
func LogMetadata(c *cli.Context) error {
	slogger.Info("Starting up "+appName,
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

	ptermLevel := pterm.LogLevelInfo
	slogLevel := slog.LevelInfo
	switch level {
	case "debug":
		ptermLevel = pterm.LogLevelDebug
		slogLevel = slog.LevelDebug
	case "warn":
		ptermLevel = pterm.LogLevelWarn
		slogLevel = slog.LevelWarn
	case "error":
		ptermLevel = pterm.LogLevelError
		slogLevel = slog.LevelError
	case "disabled":
		ptermLevel = pterm.LogLevelDisabled
		slogLevel = slog.LevelError
	}

	backend := pterm.DefaultLogger
	logHandler := pterm.NewSlogHandler(backend.WithLevel(ptermLevel))

	slogger = slog.New(logHandler)
	slog.SetDefault(slogger)
	stdLogger = slog.NewLogLogger(slogger.Handler(), slogLevel)
	if ptermLevel == pterm.LogLevelDisabled {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(stdLogger.Writer())
	}
	return nil
}

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

func newCertFilePathFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "cert-file-path",
		Usage:       "Path to the TLS certificate file",
		Value:       "/etc/web-powercycle/cert.crt",
		EnvVars:     []string{"CERT_FILE_PATH"},
		Destination: dest,
	}
}

func newCertKeyPathFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "cert-key-path",
		Usage:       "Path to the TLS certificate key",
		Value:       "/etc/web-powercycle/cert.key",
		EnvVars:     []string{"CERT_KEY_PATH"},
		Destination: dest,
	}
}

func newListenAddressFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "listen-address",
		Usage:       "Address (port) to listen on. Don't forget to prefix with ':' if listening on 0.0.0.0 or ::0",
		Value:       ":7443",
		EnvVars:     []string{"LISTEN_ADDRESS"},
		Destination: dest,
	}
}

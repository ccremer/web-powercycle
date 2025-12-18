package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var (
	// These will be populated by Goreleaser
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")

	appName     = "web-powercycle"
	appLongName = "Shut down Linux over web interface"
)

func main() {
	app := NewApp()
	err := app.Run(os.Args)
	if err != nil {
		pterm.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

func NewApp() *cli.App {
	webCommand := WebCommand{}
	app := &cli.App{
		Name:    appName,
		Usage:   appLongName,
		Version: fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date),

		Before: before(setupLogging),
		Flags: []cli.Flag{
			newLogLevelFlag(),
			newDryRunFlag(&webCommand.DryRunMode),
			newSkipSudoFlag(&webCommand.SkipSudo),
			newAuthUserFlag(&webCommand.AuthUser),
			newAuthPassFlag(&webCommand.AuthPass),
			newCertFilePathFlag(&webCommand.CertFilePath),
			newCertKeyPathFlag(&webCommand.CertKeyPath),
			newListenAddressFlag(&webCommand.ListenAddress),
			newInsecureHttpFlag(&webCommand.InsecureHttp),
		},
		Action: actions(LogMetadata, webCommand.StartWeb),
	}
	return app
}

func before(actions ...cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		for _, fn := range actions {
			if err := fn(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func actions(actions ...cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		for _, action := range actions {
			if err := action(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

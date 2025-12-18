package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
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
	app := NewRootCommand()
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		pterm.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

func NewRootCommand() *cli.Command {
	webCommand := WebCommand{}
	app := &cli.Command{
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

func before(actions ...func(context.Context, *cli.Command) (context.Context, error)) cli.BeforeFunc {
	return func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		var err error
		for _, fn := range actions {
			if ctx, err = fn(ctx, cmd); err != nil {
				return ctx, err
			}
		}
		return ctx, nil
	}
}

func actions(actions ...cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		for _, action := range actions {
			if err := action(ctx, cmd); err != nil {
				return err
			}
		}
		return nil
	}
}

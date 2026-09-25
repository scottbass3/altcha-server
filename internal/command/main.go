package command

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func Main(version string, commands ...*cli.Command) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := &cli.App{
		Version:  version,
		Name:     "altcha-server",
		Usage:    "create challenges and validate solutions for altcha captcha",
		Commands: commands,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.RunContext(ctx, os.Args); err != nil {
		logger.Error(ctx, err.Error())
		stop()
		os.Exit(1)
	}
}

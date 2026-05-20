package command

import (
	"context"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func Main(version string, commands ...*cli.Command) {
	ctx := context.Background()

	app := &cli.App{
		Version:  version,
		Name:     "altcha-server",
		Usage:    "create challenges and validate solutions for altcha captcha",
		Commands: commands,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.RunContext(ctx, os.Args); err  != nil {
		os.Exit(1)
	}
}

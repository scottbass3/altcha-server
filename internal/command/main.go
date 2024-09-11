package command

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func Main(commands ...*cli.Command) {
	ctx := context.Background()

	app := &cli.App {
		Version:	"1",
		Name: 		"altcha-server",
		Usage:		"create challenges and validate solutions for atlcha captcha",
		Commands: commands,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:		"debug",
				EnvVars:	[]string{"ALTCHA_DEBUG"},
				Value: 		false,
			},
		},
	}

	app.ExitErrHandler = func (ctx *cli.Context, err error) {
		if err == nil {
			return
		}

		debug := ctx.Bool("debug")

		if !debug {
			fmt.Printf("[ERROR] %v\n", err)
		} else {
			fmt.Printf("%+v", err)
		}
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.RunContext(ctx, os.Args); err  != nil {
		os.Exit(1)
	}
}
package command

import (
	"fmt"

	"forge.cadoles.com/cadoles/altcha-server/internal/api"
	"forge.cadoles.com/cadoles/altcha-server/internal/command/common"
	"forge.cadoles.com/cadoles/altcha-server/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/urfave/cli/v2"
)

func RunCommand() *cli.Command {
	flags := common.Flags()
	
	return &cli.Command{
		Name:	"run",
		Usage:	"run the atlcha api server",
		Flags:	flags,
		Action:	func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				fmt.Printf("%+v\n", err)
			}

			api.NewServer(cfg).Run(ctx.Context)
			return nil
		},
	}
}
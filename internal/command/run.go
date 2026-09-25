package command

import (
	"github.com/caarlos0/env/v11"
	"github.com/scottbass3/altcha-server/internal/api"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/urfave/cli/v2"
)

func RunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "run the altcha api server",
		Action: func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				return err
			}
			server, err := api.NewServer(cfg)
			if err != nil {
				return err
			}
			return server.Run(ctx.Context)
		},
	}
}

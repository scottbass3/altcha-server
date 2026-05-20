package command

import (
	"fmt"

	"github.com/scottbass3/altcha-server/internal/command/common"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "generate a challenge",
		Action: func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			client, err := common.NewClientFromConfig(cfg)
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			challenge, err := client.Generate()
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			fmt.Printf("%+v\n", challenge)
			return nil
		},
	}
}

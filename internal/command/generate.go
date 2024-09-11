package command

import (
	"fmt"

	"forge.cadoles.com/cadoles/altcha-server/internal/client"
	"forge.cadoles.com/cadoles/altcha-server/internal/command/common"
	"forge.cadoles.com/cadoles/altcha-server/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func GenerateCommand() *cli.Command {
	flags := common.Flags()

	return &cli.Command{
		Name:	"generate",
		Usage:	"generate a challenge",
		Flags:	flags,
		Action:	func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				fmt.Printf("%+v\n", err)
			}
			
			c := client.NewClient(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, cfg.Expire, cfg.CheckExpire)
			
			challenge, err := c.Generate()
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			fmt.Printf("%+v\n", challenge)

			return nil
		},
	}
}
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

func SolveCommand() *cli.Command {
	flags := common.Flags()

	return &cli.Command{
		Name:		"solve",
		Usage:		"solve the challenge and return the solution",
		Flags:		flags,
		Args: 		true,
		ArgsUsage:	"[CHALLENGE] [SALT]",
		Action:	func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				fmt.Printf("%+v\n", err)
			}

			challenge := ctx.Args().Get(0)
			salt := ctx.Args().Get(1)

			c := client.NewClient(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, salt, cfg.Expire, cfg.CheckExpire)

			
			solution, err := c.Solve(challenge)
			
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return nil
			}
			
			fmt.Printf("%+v\n", solution)

			return nil
		},
	}
}
package command

import (
	"fmt"
	"time"

	"github.com/scottbass3/altcha-server/internal/client"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func SolveCommand() *cli.Command {
	return &cli.Command{
		Name:		"solve",
		Usage:		"solve the challenge and return the solution",
		Args: 		true,
		ArgsUsage:	"[CHALLENGE] [SALT]",
		Action:	func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			challenge := ctx.Args().Get(0)
			salt := ctx.Args().Get(1)

			expirationDuration, err := time.ParseDuration(cfg.Expire)
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			client, err := client.New(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, salt, expirationDuration, cfg.CheckExpire)
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}
			
			solution, err := client.Solve(challenge)
			
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}
			
			fmt.Printf("%+v\n", solution)

			return nil
		},
	}
}

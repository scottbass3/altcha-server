package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/scottbass3/altcha-server/internal/client"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/altcha-org/altcha-lib-go"
	"github.com/caarlos0/env/v11"
	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func VerifyCommand() *cli.Command {
	return &cli.Command{
		Name:	"verify",
		Usage:	"verify the solution",
		Args:	true,
		ArgsUsage: "[challenge] [salt] [signature] [solution]",
		Action: func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			challenge := ctx.Args().Get(0)
			salt := ctx.Args().Get(1)
			signature := ctx.Args().Get(2)
			solution, _ := strconv.ParseInt(ctx.Args().Get(3), 10, 64)
			
			expirationDuration, err := time.ParseDuration(cfg.Expire+"s")
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			client, err := client.New(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, expirationDuration, cfg.CheckExpire)
			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			payload := altcha.Payload{
				Algorithm:	cfg.Algorithm,
				Challenge:	challenge,
				Number:		solution,
				Salt:		salt,
				Signature:	signature,
			}

			verified, err := client.VerifySolution(payload)

			if err != nil {
				logger.Error(ctx.Context, err.Error())
				return err
			}

			fmt.Print(verified)

			return nil
		},
	}
}

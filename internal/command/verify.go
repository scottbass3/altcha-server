package command

import (
	"fmt"
	"strconv"

	"forge.cadoles.com/cadoles/altcha-server/internal/client"
	"forge.cadoles.com/cadoles/altcha-server/internal/command/common"
	"forge.cadoles.com/cadoles/altcha-server/internal/config"
	"github.com/altcha-org/altcha-lib-go"
	"github.com/caarlos0/env/v11"
	"github.com/urfave/cli/v2"
)

func VerifyCommand() *cli.Command {
	flags := common.Flags()

	return &cli.Command{
		Name:	"verify",
		Usage:	"verify the solution",
		Flags:	flags,
		Args:	true,
		ArgsUsage: "[challenge] [salt] [signature] [solution]",
		Action: func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				fmt.Printf("%+v\n", err)
			}

			challenge := ctx.Args().Get(0)
			salt := ctx.Args().Get(1)
			signature := ctx.Args().Get(2)
			solution, _ := strconv.ParseInt(ctx.Args().Get(3), 10, 64)
			
			c := client.NewClient(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, cfg.Expire, cfg.CheckExpire)

			payload := altcha.Payload{
				Algorithm:	cfg.Algorithm,
				Challenge:	challenge,
				Number:		solution,
				Salt:		salt,
				Signature:	signature,
			}

			verified, err := c.VerifySolution(payload)

			if err != nil {
				return err
			}

			fmt.Print(verified)

			return nil
		},
	}
}
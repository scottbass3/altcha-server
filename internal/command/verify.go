package command

import (
	"fmt"
	"strconv"

	"github.com/altcha-org/altcha-lib-go"
	"github.com/caarlos0/env/v11"
	"github.com/scottbass3/altcha-server/internal/command/common"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/urfave/cli/v2"
)

func VerifyCommand() *cli.Command {
	return &cli.Command{
		Name:      "verify",
		Usage:     "verify the solution",
		Args:      true,
		ArgsUsage: "[challenge] [salt] [signature] [solution]",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() < 4 {
				return fmt.Errorf("usage: verify [challenge] [salt] [signature] [solution]")
			}

			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				return err
			}

			challenge := ctx.Args().Get(0)
			salt := ctx.Args().Get(1)
			signature := ctx.Args().Get(2)
			solution, err := strconv.ParseInt(ctx.Args().Get(3), 10, 64)
			if err != nil {
				return fmt.Errorf("invalid solution %q: must be an integer", ctx.Args().Get(3))
			}

			client, err := common.NewClientFromConfig(cfg)
			if err != nil {
				return err
			}

			payload := altcha.Payload{
				Algorithm: cfg.Algorithm,
				Challenge: challenge,
				Number:    solution,
				Salt:      salt,
				Signature: signature,
			}

			verified, err := client.VerifySolution(payload)
			if err != nil {
				return err
			}

			fmt.Print(verified)
			return nil
		},
	}
}

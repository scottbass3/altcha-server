package command

import (
	"errors"
	"fmt"

	"github.com/altcha-org/altcha-lib-go"
	"github.com/caarlos0/env/v11"
	"github.com/scottbass3/altcha-server/internal/client"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/urfave/cli/v2"
)

func SolveCommand() *cli.Command {
	return &cli.Command{
		Name:      "solve",
		Usage:     "solve the challenge and return the solution",
		Args:      true,
		ArgsUsage: "[CHALLENGE] [SALT]",
		Action: func(ctx *cli.Context) error {
			cfg := config.Config{}
			if err := env.Parse(&cfg); err != nil {
				return err
			}

			challenge := ctx.Args().Get(0)
			if salt := ctx.Args().Get(1); salt != "" {
				cfg.Salt = salt
			}

			algorithm, err := client.ParseAlgorithm(cfg.Algorithm)
			if err != nil {
				return err
			}

			solution, err := altcha.SolveChallenge(challenge, cfg.Salt, algorithm, int(cfg.MaxNumber), 0, ctx.Context.Done())
			if err != nil {
				return err
			}
			if solution == nil {
				return errors.New("no solution found")
			}

			fmt.Printf("%+v\n", solution)
			return nil
		},
	}
}

package db

import (
	cli "github.com/urfave/cli/v2"
	"github.com/toncek345/webshop/config"
)

var fixture *cli.Command = &cli.Command{
	Name:  "fixture",
	Usage: "loads fixtures for given env",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "env",
			Value: "development",
			Usage: "environment mode to run application",
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.New(config.Environment(ctx.String("env")))
		if err != nil {
			return err
		}

		fixture, err := cfg.Fixture()
		if err != nil {
			return err
		}

		if err := fixture.Load(); err != nil {
			return err
		}

		return nil
	},
}

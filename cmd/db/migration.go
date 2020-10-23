package db

import (
	"github.com/toncek345/webshop/config"
	cli "github.com/urfave/cli/v2"
)

var up *cli.Command = &cli.Command{
	Name:  "up",
	Usage: "runs pending migrations",
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

		m, err := cfg.MigrationObj()
		if err != nil {
			return err
		}

		if err := m.Up(); err != nil {
			return err
		}

		return nil
	},
}

var down *cli.Command = &cli.Command{
	Name:  "down",
	Usage: "removes recent migration",
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

		m, err := cfg.MigrationObj()
		if err != nil {
			return err
		}

		if err := m.Down(); err != nil {
			return err
		}

		return nil
	},
}

var drop *cli.Command = &cli.Command{
	Name:  "drop",
	Usage: "drops database",
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

		m, err := cfg.MigrationObj()
		if err != nil {
			return err
		}

		if err := m.Drop(); err != nil {
			return err
		}

		return nil
	},
}

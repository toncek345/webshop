package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/toncek345/webshop/config"
	cli "github.com/urfave/cli/v2"
)

func getCfg(env string) (cfg config.Config, err error) {
	switch env {
	case "development":
		cfg, err = config.New(config.DevEnv)
	case "production":
		cfg, err = config.New(config.ProdEnv)
	}

	return
}

func RegisterCommand(app *cli.App) {
	app.Commands = append(
		app.Commands,
		&cli.Command{
			Name:  "migrate",
			Usage: "migration actions",
			Subcommands: []*cli.Command{
				{
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
						cfg, err := getCfg(ctx.String("env"))
						if err != nil {
							return err
						}

						m, err := migrate.New(
							"file://migrations",
							cfg.PGSQLConnString())
						if err != nil {
							return err
						}

						if err := m.Up(); err != nil {
							return err
						}

						return nil
					},
				},
				{
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
						cfg, err := getCfg(ctx.String("env"))
						if err != nil {
							return err
						}

						m, err := migrate.New(
							"file://migrations",
							cfg.PGSQLConnString())
						if err != nil {
							return err

						}

						v, _, err := m.Version()
						if err != nil && err != migrate.ErrNilVersion {
							return err
						}

						v--
						if err := m.Migrate(v); err != nil {
							return err
						}

						return nil
					},
				},
				{
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
						cfg, err := getCfg(ctx.String("env"))
						if err != nil {
							return err
						}

						m, err := migrate.New(
							"file://migrations",
							cfg.PGSQLConnString())
						if err != nil {
							return err
						}

						if err := m.Drop(); err != nil {
							return err
						}

						return nil
					},
				},
			},
		},
	)
}

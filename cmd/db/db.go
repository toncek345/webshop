package db

import (
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
			Name:  "db",
			Usage: "db actions",
			Subcommands: []*cli.Command{
				up, down, drop,
				fixture,
			},
		},
	)
}

package db

import (
	cli "github.com/urfave/cli/v2"
)

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

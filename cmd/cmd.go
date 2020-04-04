package cmd

import (
	"github.com/toncek345/webshop/cmd/migrate"
	"github.com/toncek345/webshop/cmd/web"
	cli "github.com/urfave/cli/v2"
)

func RegisterCmds() *cli.App {
	app := &cli.App{}

	web.RegisterCommand(app)
	migrate.RegisterCommand(app)

	return app
}

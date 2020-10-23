package cmd

import (
	"github.com/toncek345/webshop/cmd/db"
	"github.com/toncek345/webshop/cmd/web"
	cli "github.com/urfave/cli/v2"
)

func RegisterCmds() *cli.App {
	app := &cli.App{}

	web.RegisterCommand(app)
	db.RegisterCommand(app)

	return app
}

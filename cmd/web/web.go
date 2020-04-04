package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/senko/clog"
	"github.com/toncek345/webshop/api"
	"github.com/toncek345/webshop/config"
	"github.com/toncek345/webshop/models"
	cli "github.com/urfave/cli/v2"
)

// TODO: fixtures
// dbInit = flag.Bool(
// 	"dbInit",
// 	false,
// 	"set to true when db is done and tables need to be created, run only once")

// TODO: fixtures
// if *dbInit {
// 	if err := models.InitDb(); err != nil {
// 		panic(err)
// 	}
// 	clog.Info("db init completed")
// 	os.Exit(0)
// }

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
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "env",
					Value: "development",
					Usage: "environment mode to run application",
				},
			},
			Name:  "web",
			Usage: "starts webserver",
			Action: func(ctx *cli.Context) error {
				clog.Setup(clog.DEBUG, false)
				config, err := getCfg(ctx.String("env"))
				if err != nil {
					return err
				}

				sqlConn, err := config.SqlxConn()
				if err != nil {
					return err
				}

				models, err := models.New(sqlConn)
				if err != nil {
					return err
				}

				webApp := api.New(models, config.StaticPath)

				addr := fmt.Sprintf(":%s", strconv.FormatInt(config.WebPort, 10))
				server := http.Server{
					Handler:      webApp.Router(),
					Addr:         addr,
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
				}

				clog.Debug("webshop running")
				clog.Errorf("%v", server.ListenAndServe())
				return nil
			},
		},
	)
}

func main() {

}

package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/senko/clog"
	"github.com/toncek345/webshop/config"
	"github.com/toncek345/webshop/internal/api"
	cli "github.com/urfave/cli/v2"
)

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
				config, err := config.New(config.Environment(ctx.String("env")))
				if err != nil {
					return err
				}

				db, err := config.SqlxConn()
				if err != nil {
					return err
				}

				storage, err := config.Storage()
				if err != nil {
					return err
				}

				r := chi.NewRouter()
				r.Mount("/", api.New(db, storage))
				// static folder serves only images and other non front static files
				r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
					fs := http.StripPrefix("/api/static", http.FileServer(http.Dir(config.StaticPath)))
					fs.ServeHTTP(w, r)
				})

				addr := fmt.Sprintf(":%s", strconv.FormatInt(config.WebPort, 10))
				server := http.Server{
					Handler:      r,
					Addr:         addr,
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
				}

				clog.Debugf("webshop running on %s", addr)
				clog.Errorf("%v", server.ListenAndServe())
				return nil
			},
		},
	)
}

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

				db, err := config.SqlxConn()
				if err != nil {
					return err
				}

				r := chi.NewRouter()
				r.Mount("/", api.New(db, config.DiskStoragePath))
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

func main() {

}

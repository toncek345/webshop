package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/senko/clog"
	"github.com/toncek345/webshop/api"
	"github.com/toncek345/webshop/config"
	"github.com/toncek345/webshop/models"
)

var (
	environment = flag.String(
		"env",
		"development",
		"development, production running mode",
	)

	// TODO: fixtures
	// dbInit = flag.Bool(
	// 	"dbInit",
	// 	false,
	// 	"set to true when db is done and tables need to be created, run only once")
)

func main() {
	flag.Parse()
	clog.Setup(clog.DEBUG, false)

	// TODO: fixtures
	// if *dbInit {
	// 	if err := models.InitDb(); err != nil {
	// 		panic(err)
	// 	}
	// 	clog.Info("db init completed")
	// 	os.Exit(0)
	// }

	var env config.Environment
	switch *environment {
	case "development":
		env = config.DevEnv
	case "production":
		env = config.ProdEnv
	}

	config, err := config.New(env)
	if err != nil {
		panic(err)
	}

	sqlConn, err := config.SqlxConn()
	if err != nil {
		panic(err)
	}

	models, err := models.New(sqlConn)
	if err != nil {
		panic(err)
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
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/senko/clog"
	"github.com/toncek345/webshop/models"
	"github.com/toncek345/webshop/web"
)

var (
	portNo = flag.Int64(
		"port",
		9000,
		"listening port number")

	pathToStatic = flag.String(
		"static",
		"./static/",
		"full path to static folder, add trailing / => static is currently used only for saving pictures")

	dbConnectionString = flag.String(
		"dbString",
		"user=postgres dbname=webshop sslmode=disable",
		"database connection string, currenty only postgres supported")

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

	sqlConn, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}

	models, err := models.New(sqlConn)
	if err != nil {
		panic(err)
	}

	webApp := web.New(models, *pathToStatic)

	addr := fmt.Sprintf(":%s", strconv.FormatInt(*portNo, 10))
	server := http.Server{
		Handler:      webApp.Router(),
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	clog.Debug("webshop running")
	clog.Errorf("%v", server.ListenAndServe())
}

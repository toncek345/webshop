package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webshop/models"
	"webshop/urls"

	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/senko/clog"
)

func main() {
	clog.Setup(clog.DEBUG, false)

	// setup prog input
	portNo := flag.Int64("port", 9000, "listening port number")
	pathToStatic := flag.String("static", "./static/", "full path to static folder, add trailing / => static is currently used only for saving pictures")
	dbConnectionString := flag.String("dbString",
		"user=postgres password=postgres dbname=webshopGo sslmode=disable",
		"database connection string, currenty only postgres supported")
	dbInit := flag.Bool("dbInit", false,
		"set to true when db is done and tables need to be created, run only once")
	flag.Parse()

	err := models.DbConnect(*dbConnectionString)
	if err != nil {
		panic(err)
	}

	if *dbInit {
		if err := models.InitDb(); err != nil {
			panic(err)
		}
		clog.Info("db init completed")
		os.Exit(0)
	}

	r := mux.NewRouter()
	urls.SetUrls(r, *pathToStatic)

	addr := fmt.Sprintf("0.0.0.0:%s", strconv.FormatInt(*portNo, 10))
	server := http.Server{
		Addr: addr,
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	clog.Debug("webshop running")
	clog.Errorf("%v", server.ListenAndServe())
}

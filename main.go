package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/senko/clog"
	"github.com/toncek345/webshop/models"
	"github.com/toncek345/webshop/urls"
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

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	urls.SetUrls(r, *pathToStatic)

	addr := fmt.Sprintf("0.0.0.0:%s", strconv.FormatInt(*portNo, 10))
	server := http.Server{
		Handler:      r,
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	clog.Debug("webshop running")
	clog.Errorf("%v", server.ListenAndServe())
}

package main

import (
	"net/http"
	"time"
	"webshop/urls"

	"github.com/gorilla/mux"
	"github.com/senko/clog"

	_ "github.com/lib/pq"
)

func main() {
	clog.Setup(clog.DEBUG, false)
	pathToStatic := "static"

	r := mux.NewRouter()
	urls.SetUrls(r)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(pathToStatic)))).Name("static")

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	clog.Errorf("%v", server.ListenAndServe())
}

package urls

import (
	"net/http"
	"webshop/models"

	"github.com/gorilla/mux"
	"github.com/senko/clog"
)

func newsUrls(r *mux.Router) {
	r.HandleFunc("/news", getNews).Methods("GET")

	r.HandleFunc("/news",
		authenticationRequired(
			createNews)).Methods("POST")

	r.HandleFunc("/news/{id:[0-9]+}",
		authenticationRequired(
			deleteNews)).Methods("DELETE")

	r.HandleFunc("/news/{id:[0-9]+}",
		authenticationRequired(
			updateNews)).Methods("PUT")
}

func getNews(w http.ResponseWriter, r *http.Request) {
	n, err := models.GetAllNews()
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, n)
}

func createNews(w http.ResponseWriter, r *http.Request) {
	var obj models.News

	err := decode(r, &obj)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err = models.CreateNews(obj)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func deleteNews(w http.ResponseWriter, r *http.Request) {
	id, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err = models.DeleteNewsById(id)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func updateNews(w http.ResponseWriter, r *http.Request) {
	var n models.News

	id, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err = decode(r, &n)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err = models.UpdateNewsById(id, n)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

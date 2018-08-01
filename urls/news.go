package urls

import (
	"net/http"
	"webshop/models"

	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
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
	type newNews struct {
		News  models.News
		Image string // base64 encoded image
	}
	var obj newNews

	err := decode(r, &obj)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	if obj.Image == "" {
		clog.Warningf("empty image for news")
		respond(w, r, http.StatusBadRequest, "error: empty image")
		return
	}

	binaryImage, err := base64.StdEncoding.DecodeString(obj.Image)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	imageFilename := fmt.Sprintf("news-%s.jpg", uuid.NewV4().String())
	ioutil.WriteFile(staticFolderPath+imageFilename, binaryImage, os.ModePerm)
	obj.News.ImagePath = imageFilename

	err = models.CreateNews(obj.News)
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

	news, err := models.DeleteNewsById(id)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	os.Remove(staticFolderPath + news.ImagePath)

	respond(w, r, http.StatusOK, nil)
}

func updateNews(w http.ResponseWriter, r *http.Request) {
	type newNews struct {
		News  models.News
		Image string // base64 encoded image
	}
	var obj newNews

	id, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err = decode(r, &obj)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	n, err := models.GetNewsById(id)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	data, err := base64.StdEncoding.DecodeString(obj.Image)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	ioutil.WriteFile(staticFolderPath+n.ImagePath, data, os.ModePerm)

	err = models.UpdateNewsById(id, n)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

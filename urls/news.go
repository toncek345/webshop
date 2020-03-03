package urls

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/senko/clog"
	"github.com/toncek345/webshop/models"
)

func newsUrls(r chi.Router) {
	r.Route("/news", func(r chi.Router) {
		r.Get("/", getNews)

		r.Group(func(r chi.Router) {
			r.Use(authenticationRequired)

			r.Post("/", createNews)
			r.Delete("/{id:[0-9]+}", deleteNews)
			r.Put("/{id:[0-9]+}", updateNews)
		})
	})
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

package api

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

func (app *App) newsRouter(r chi.Router) {
	r.Route("/news", func(r chi.Router) {
		r.Get("/", app.getNews)

		r.Group(func(r chi.Router) {
			r.Use(app.authenticationRequired)

			r.Post("/", app.createNews)
			r.Delete("/{id:[0-9]+}", app.deleteNews)
			r.Put("/{id:[0-9]+}", app.updateNews)
		})
	})
}

func (app *App) getNews(w http.ResponseWriter, r *http.Request) {
	n, err := app.models.News.Get()
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, n)
}

func (app *App) createNews(w http.ResponseWriter, r *http.Request) {
	type newNews struct {
		News  models.News
		Image string // base64 encoded image
	}
	var obj newNews

	err := app.JSONDecode(r, &obj)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	if obj.Image == "" {
		clog.Warningf("empty image for news")
		app.JSONRespond(w, r, http.StatusBadRequest, "error: empty image")
		return
	}

	binaryImage, err := base64.StdEncoding.DecodeString(obj.Image)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	imageFilename := fmt.Sprintf("news-%s.jpg", uuid.NewV4().String())
	ioutil.WriteFile(app.staticFolderPath+imageFilename, binaryImage, os.ModePerm)
	obj.News.ImageKey = imageFilename

	err = app.models.News.CreateNews(obj.News)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

func (app *App) deleteNews(w http.ResponseWriter, r *http.Request) {
	id, err := parseUrlVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	news, err := app.models.News.GetByID(id)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := app.models.News.DeleteByID(id); err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	os.Remove(app.staticFolderPath + news.ImageKey)
	app.JSONRespond(w, r, http.StatusOK, nil)
}

func (app *App) updateNews(w http.ResponseWriter, r *http.Request) {
	type newNews struct {
		News  models.News
		Image string // base64 encoded image
	}
	var obj newNews

	id, err := parseUrlVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	err = app.JSONDecode(r, &obj)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	n, err := app.models.News.GetByID(id)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	data, err := base64.StdEncoding.DecodeString(obj.Image)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	ioutil.WriteFile(app.staticFolderPath+n.ImageKey, data, os.ModePerm)

	err = app.models.News.UpdateByID(id, n)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

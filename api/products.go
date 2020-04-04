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

func (app *App) productRouter(r chi.Router) {

	r.Route("/product", func(r chi.Router) {
		r.Get("/", app.getProducts)

		r.Group(func(r chi.Router) {
			r.Use(app.authenticationRequired)

			r.Post("/", app.createProduct)
			r.Delete("/{id:[0-9]+}", app.productDelete)
			r.Put("/{id:[0-9]+}", app.productUpdate)

			r.Route("/image", func(r chi.Router) {
				r.Post("/{id:[0-9]+}", app.createImage)
				r.Delete("/{id:[0-9]+}", app.deleteImage)
			})
		})
	})
}

func (app *App) createImage(w http.ResponseWriter, r *http.Request) {
	var base64Image string

	if err := app.JSONDecode(r, &base64Image); err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	productId, err := parseUrlVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	binaryImage, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	filename := fmt.Sprintf("product-%s.jpg", uuid.NewV4().String())
	ioutil.WriteFile(app.staticFolderPath+filename, binaryImage, os.ModePerm)

	err = app.models.Products.InsertImages(productId, []string{filename})
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

func (app *App) deleteImage(w http.ResponseWriter, r *http.Request) {
	imageId, err := parseUrlVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	imgName, err := app.models.Products.DeleteImage(imageId)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	err = os.Remove(app.staticFolderPath + imgName)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	type newProduct struct {
		Product models.Product
		Images  []string // base64 encoded images
	}

	var obj newProduct

	if err := app.JSONDecode(r, &obj); err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	productId, err := app.models.Products.Create(obj.Product)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	if len(obj.Images) > 0 {
		binaryImages := make([][]byte, len(obj.Images))
		imageNames := []string{}

		for i, _ := range obj.Images {
			data, err := base64.StdEncoding.DecodeString(obj.Images[i])
			if err != nil {
				clog.Warningf("%s", err)
				app.JSONRespond(w, r, http.StatusInternalServerError, err)
				return
			}

			filename := fmt.Sprintf("product-%s.jpg", uuid.NewV4().String())
			ioutil.WriteFile(app.staticFolderPath+filename, data, os.ModePerm)
			imageNames = append(imageNames, filename)

			binaryImages[i] = data
		}

		if err := app.models.Products.InsertImages(productId, imageNames); err != nil {
			clog.Warningf("%s", err)
			app.JSONRespond(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.models.Products.Get()
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, products)
}

func (app *App) productDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUrlVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	if err = app.models.Products.DeleteByID(id); err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

func (app *App) productUpdate(w http.ResponseWriter, r *http.Request) {
	var n models.Product

	id, err := parseUrlVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	err = app.JSONDecode(r, &n)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusBadRequest, err)
		return
	}

	err = app.models.Products.UpdateByID(id, n)
	if err != nil {
		clog.Warningf("%s", err)
		app.JSONRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	app.JSONRespond(w, r, http.StatusOK, nil)
}

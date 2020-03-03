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

func productUrls(r chi.Router) {

	r.Route("/product", func(r chi.Router) {
		r.Get("/", getProducts)

		r.Group(func(r chi.Router) {
			r.Use(authenticationRequired)

			r.Post("/", createProduct)
			r.Delete("/{id:[0-9]+}", productDelete)
			r.Put("/{id:[0-9]+}", productUpdate)

			r.Route("/image", func(r chi.Router) {
				r.Post("/{id:[0-9]+}", createImage)
				r.Delete("/{id:[0-9]+}", deleteImage)
			})
		})
	})
}

func createImage(w http.ResponseWriter, r *http.Request) {
	var base64Image string

	if err := decode(r, &base64Image); err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	productId, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	binaryImage, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	filename := fmt.Sprintf("product-%s.jpg", uuid.NewV4().String())
	ioutil.WriteFile(staticFolderPath+filename, binaryImage, os.ModePerm)

	err = models.InsertImages(productId, []string{filename})
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	imageId, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	imgName, err := models.DeleteImage(imageId)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	err = os.Remove(staticFolderPath + imgName)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	type newProduct struct {
		Product models.Product
		Images  []string // base64 encoded images
	}

	var obj newProduct

	if err := decode(r, &obj); err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	productId, err := models.CreateProduct(obj.Product)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if len(obj.Images) > 0 {
		binaryImages := make([][]byte, len(obj.Images))
		imageNames := []string{}

		for i, _ := range obj.Images {
			data, err := base64.StdEncoding.DecodeString(obj.Images[i])
			if err != nil {
				clog.Warningf("%s", err)
				respond(w, r, http.StatusInternalServerError, err)
				return
			}

			filename := fmt.Sprintf("product-%s.jpg", uuid.NewV4().String())
			ioutil.WriteFile(staticFolderPath+filename, data, os.ModePerm)
			imageNames = append(imageNames, filename)

			binaryImages[i] = data
		}

		if err := models.InsertImages(productId, imageNames); err != nil {
			clog.Warningf("%s", err)
			respond(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	respond(w, r, http.StatusOK, nil)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetAllProducts()
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, products)
}

func productDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	_, err = models.DeleteProductById(id)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func productUpdate(w http.ResponseWriter, r *http.Request) {
	var n models.Product

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

	err = models.UpdateProductById(id, n)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

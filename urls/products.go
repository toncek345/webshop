package urls

import (
	"net/http"
	"webshop/models"

	"github.com/gorilla/mux"
	"github.com/senko/clog"
)

func productUrls(r *mux.Router) {
	r.HandleFunc("/product", logRoute(
		getProducts)).Methods("GET")

	r.HandleFunc("/product", logRoute(
		authenticationRequired(
			createProduct))).Methods("POST")

	r.HandleFunc("/product/{id:[0-9]+}",
		logRoute(
			authenticationRequired(
				productDelete))).Methods("DELETE")

	r.HandleFunc("/product/{id:[0-9]+}",
		logRoute(
			authenticationRequired(
				productUpdate))).Methods("PUT")

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

func createProduct(w http.ResponseWriter, r *http.Request) {
	var obj models.Product

	if err := decode(r, &obj); err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err := models.CreateProduct(obj)
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func productDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseMuxVarsInt(r, "id")
	if err != nil {
		clog.Warningf("%s", err)
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	err = models.DeleteProductById(id)
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

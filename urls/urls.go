package urls

// package file that only loads individual urls and provides functions to be
// used outside of urls package

import (
	"net/http"

	"webshop/models"

	"encoding/json"

	"io/ioutil"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/senko/clog"
)

func init() {

}

func SetUrls(r *mux.Router) {
	// / [get]
	// /products [get, post]
	// /products/{id} [get, del, put]
	// /news [get, post]
	// /news/{id} [get, del, put]
	// /admin [get]
	// /admin/login [post]
	// /admin/logout [post]

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// serve static page
		w.Write([]byte("helo homepage"))
	}).Methods("GET")
}

func returnJsonResp(obj interface{}, w *http.ResponseWriter) {
	data, err := json.Marshal(obj)
	if err != nil {
		clog.Warningf("error marshalling obj: %s", err)
	}

	(*w).Write(data)
}

func parseId(r *http.Request, name string) int {
	id, err := strconv.ParseInt(mux.Vars(r)[name], 10, 32)
	if err != nil {
		clog.Warningf("id parse error: %s", err)
	}
	return int(id)
}

func adminUrls(r *mux.Router) {
	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		// serves admin page
	}).Methods("GET")

	r.HandleFunc("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		// login admin -> token in header
	}).Methods("POST")

	r.HandleFunc("/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		// logout admin -> token in header
	}).Methods("POST")
}

func newsUrls(r *mux.Router) {
	r.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		// gets all news
		returnJsonResp(models.GetAllNews(), &w)
	}).Methods("GET")

	r.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		// create news with auth -> admin only
	}).Methods("POST")

	r.HandleFunc("/news/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		// get one article
		// TODO: ako get all vrati sve koji ce mi ovaj path onda?
	}).Methods("GET")

	r.HandleFunc("/news/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		// delete, update one article with auth -> admin only
	}).Methods("DELETE", "PUT")
}

func productUrls(r *mux.Router) {
	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		// gets all products
		prod := models.GetAllProducts()
		returnJsonResp(prod, &w)
	}).Methods("GET")

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		// create product with auth -> admin only
		// TODO: user needs to be logged in
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			clog.Warningf("/products:[%s] error reading: %s", r.Method, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		obj := models.Product{}
		json.Unmarshal(body, &obj)
		if err != nil {
			clog.Warningf("/products:[%s] error unmarshalling", r.Method, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ok := models.CreateProduct(obj)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/products/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		// gets one product
		// TODO: ako get all vrati sve za sta mi treba onda ovaj path?
		id := parseId(r, "id")

		obj, err := models.GetProductById(id)
		if err != nil {
			clog.Warningf("/products/[id:%d][%s], get by id error: %s",
				id, r.Method, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.Marshal(obj)
	}).Methods("GET")

	r.HandleFunc("/products/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		// del, put with auth -> admin only
		// TODO: auth
		id := parseId(r, "id")

		switch r.Method {
		case "DELETE":
			err := models.DeleteProductById(id)
			if err != nil {
				clog.Warningf("/products/[id:%d][%s] error deleting product: %s",
					id, r.Method, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "PUT":
			// TODO:
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE", "PUT")
}

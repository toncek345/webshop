// Package urls provides all application routes with handlers.
// 		/ [get]
// 		/product [get, post]
// 		/product/{id} [get, del, put]

//		/product/{id}/image [post]
//		/image/{id} [delete]

// 		/news [get, post]
// 		/news/{id} [get, del, put]

// 		/admin [get]
// 		/user/login [post]
// 		/user/logout [post]
package urls

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/senko/clog"
)

var staticFolderPath string

func SetUrls(r chi.Router, staticFolder string) {
	staticFolderPath = staticFolder

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("all good"))
		})

		newsUrls(r)
		adminUrls(r)
		productUrls(r)

		// static folder serves only images and other non front static files
		r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
			fs := http.StripPrefix("/api/v1/static", http.FileServer(http.Dir(staticFolderPath)))
			fs.ServeHTTP(w, r)
		})
	})
}

// Gets auth uuid from header which is sent in x-auth.
func getAuthHeader(r *http.Request) string {
	return r.Header.Get("x-auth")
}

// Decodes from given request to object that is given.
// Call this function with pointer to object that needs to be decoded.
func decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

// Generic respond func.
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		clog.Errorf("error marshalling obj: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		clog.Errorf("error returning data: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Parses from mux.vars[name] to integer.
func parseMuxVarsInt(r *http.Request, name string) (int, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, name), 10, 32)
	if err != nil {
		clog.Warningf("id parse error: %s", err)
	}
	return int(id), nil
}

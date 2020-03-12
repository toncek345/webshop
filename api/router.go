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
package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/toncek345/webshop/models"
)

type App struct {
	staticFolderPath string

	models models.Models
}

func New(
	models models.Models,
	staticFolderPath string) App {

	return App{
		models:           models,
		staticFolderPath: staticFolderPath,
	}
}

func (app *App) Router() chi.Router {
	r := chi.NewRouter()

	// TODO: Handling cors should be only in development since react is with dev server
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("all good"))
		})

		app.newsRouter(r)
		app.adminRouter(r)
		app.productRouter(r)

		// static folder serves only images and other non front static files
		r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
			fs := http.StripPrefix("/api/static", http.FileServer(http.Dir(app.staticFolderPath)))
			fs.ServeHTTP(w, r)
		})
	})

	return r
}

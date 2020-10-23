package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/toncek345/webshop/internal/api/v1/model"
	"github.com/toncek345/webshop/internal/pkg/storage"
)

type App struct {
	storage storage.Storage
	models  model.Models
}

func New(
	models model.Models,
	storage storage.Storage) *App {

	return &App{
		models:  models,
		storage: storage,
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

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("all good"))
		})

		app.adminRouter(r)
		app.productRouter(r)
	})

	return r
}

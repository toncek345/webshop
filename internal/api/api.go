package api

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	v1 "github.com/toncek345/webshop/internal/api/v1"
)

func New(db *sqlx.DB, storagePath string) chi.Router {
	r := chi.NewMux()

	r.Mount("/api", v1.New(db, storagePath))

	return r
}

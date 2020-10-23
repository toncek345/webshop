package api

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	v1 "github.com/toncek345/webshop/internal/api/v1"
	"github.com/toncek345/webshop/internal/pkg/storage"
)

func New(db *sqlx.DB, storage storage.Storage) chi.Router {
	r := chi.NewMux()

	r.Mount("/api", v1.New(db, storage))

	return r
}

package v1

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	apiV1 "github.com/toncek345/webshop/internal/api/v1/handler"
	"github.com/toncek345/webshop/internal/api/v1/model"
	"github.com/toncek345/webshop/internal/pkg/storage"
)

func New(db *sqlx.DB, storage storage.Storage) chi.Router {
	r := chi.NewRouter()

	model, err := model.New(db)
	if err != nil {
		panic(err)
	}

	r.Mount("/v1", apiV1.New(model, storage).Router())
	return r
}

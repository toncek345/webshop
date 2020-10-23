package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Models struct {
	Products productsRepo
	Users    usersRepo
	Auth     authRepo
}

func New(sqlDB *sqlx.DB) (Models, error) {
	if sqlDB == nil {
		return Models{}, fmt.Errorf("models: models init failed sqldb is nil")
	}

	return Models{
		Products: newProductsRepo(sqlDB),
		Users:    newUserRepo(sqlDB),
		Auth:     newAuthRepo(sqlDB),
	}, nil
}

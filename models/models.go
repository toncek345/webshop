package models

import (
	"database/sql"
	"fmt"
)

type Models struct {
	News     newsRepo
	Products productsRepo
	Users    usersRepo
	Auth     authRepo
}

func New(sqlDB *sql.DB) (Models, error) {
	if sqlDB == nil {
		return Models{}, fmt.Errorf("models: models init failed sqldb is nil")
	}

	return Models{
		News:     newNewsRepo(sqlDB),
		Products: newProductsRepo(sqlDB),
		Users:    newUserRepo(sqlDB),
		Auth:     newAuthRepo(sqlDB),
	}, nil
}

// 	// creating admin
// 	var hash []byte
// 	hash, err = bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
// 	if err != nil {
// 		clog.Errorf("%s", err)
// 		return
// 	}

// 	admin := User{
// 		Username: "admin",
// 		Password: string(hash),
// 	}

// 	err = CreateUser(admin)
// 	if err != nil {
// 		clog.Errorf("%s", err)
// 		return
// 	}

// 	return err
// }

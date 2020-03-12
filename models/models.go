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

// func InitDb() error {
// if err := initNews(); err != nil {
// 		sql := `CREATE TABLE public.news (
// 	id serial NOT NULL PRIMARY KEY,
// 	header varchar(250),
// 	text text,
// 	imagepath varchar(250)
// )`

// _, err = sqlDB.Exec(sql)
// 	return err
// }

// if err := initProduct(); err != nil {
// 		sqlProduct := `CREATE TABLE public.product (
// id serial NOT NULL PRIMARY KEY,
//   price integer,
//   name varchar(250),
//   description text
// )`

// _, err = sqlDB.Query(sqlProduct)
// if err != nil {
// 	return
// }

// sqlProductImages := `CREATE TABLE public.images (
// id serial NOT NULL PRIMARY KEY,
// product_id integer NOT NULL REFERENCES public.product(id) ON DELETE CASCADE,
// name varchar(250)
// )`
// _, err = sqlDB.Query(sqlProductImages)

// 	return err
// }

// if err := initUser(); err != nil {
// TODO: admin creaton belongs to fixtures
// 			sql := `CREATE TABLE public.user (
// 	id serial NOT NULL PRIMARY KEY,
//     username varchar(250),
//     password varchar(60)
// )`

// 	_, err = sqlDB.Query(sql)
// 	if err != nil {
// 		clog.Errorf("%s", err)
// 		return
// 	}

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

// if err := initAuth(); err != nil {
// 		sql := `CREATE TABLE public.authentification (
// id serial NOT NULL PRIMARY KEY,
//   user_id int NOT NULL REFERENCES public.user (id) ON DELETE CASCADE,
//   valid_until timestamp,
//   token uuid
// )`
// _, err = sqlDB.Query(sql)
// return err
// }

// 	return nil
// }

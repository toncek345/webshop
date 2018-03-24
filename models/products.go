package models

import (
	"database/sql"
	"errors"
)

type Product struct {
	Id          int
	Price       int
	Name        string
	Description string
	ImagePath   string
}

func initProduct() (err error) {
	sql := `CREATE TABLE public.product (
	id serial NOT NULL PRIMARY KEY,
    price integer,
    name varchar(250),
    description text,
    imagepath varchar(250)
	)`

	_, err = sqlDB.Query(sql)
	return
}

// sql-s
var (
	// select
	getAllProducts = "SELECT * FROM public.product"
	getProductById = "SELECT * FROM public.product p WHERE p.id = $1"

	// update
	updateProduct = "UPDATE public.product " +
		"SET price=$1, name=$2, description=$3, imagepath=$4 " +
		"WHERE id = $5"

	// delete
	deleteProduct = "DELETE FROM public.product WHERE id=$1"

	// create
	createProduct = "INSERT INTO public.product (price, name, description, imagepath) " +
		"VALUES ($1, $2, $3, $4)"
)

var (
	ProductNotCreatedError = errors.New("Product not created")
	NoSuchIdProductError   = errors.New("Product not found by given ID")
)

func GetAllProducts() (p []Product, err error) {
	var res *sql.Rows
	res, err = sqlDB.Query(getAllProducts)
	if err != nil {
		return
	}

	for res.Next() {
		temp := Product{}
		err = res.Scan(&temp.Id, &temp.Price, &temp.Name, &temp.Description, &temp.ImagePath)

		p = append(p, temp)
	}

	return
}

func GetProductById(id int) (p Product, err error) {
	var res *sql.Row
	res = sqlDB.QueryRow(getProductById, id)
	err = res.Scan(&p.Id, &p.Price, &p.Name, &p.Description, &p.ImagePath)
	if err != nil {
		return
	}

	return
}

func DeleteProductById(id int) error {
	res, err := sqlDB.Exec(deleteProduct, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return NoSuchIdProductError
	}

	return nil
}

func UpdateProductById(id int, p Product) error {
	res, err := sqlDB.Exec(updateProduct, p.Price, p.Name, p.Description, p.ImagePath, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return NoSuchIdProductError
	}

	return nil
}

func CreateProduct(p Product) error {
	res, err := sqlDB.Exec(createProduct, p.Price, p.Name, p.Description, p.ImagePath)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ProductNotCreatedError
	}

	return nil
}

package models

import (
	"database/sql"
	"errors"
)

type Images struct {
	Id        int
	ProductId int
	Name      string
}

type Product struct {
	Id          int
	Price       int
	Name        string
	Description string
	Images      []Images
}

func initProduct() (err error) {
	sqlProduct := `CREATE TABLE public.product (
	id serial NOT NULL PRIMARY KEY,
    price integer,
    name varchar(250),
    description text
	)`

	_, err = sqlDB.Query(sqlProduct)
	if err != nil {
		return
	}

	sqlProductImages := `CREATE TABLE public.images (
	id serial NOT NULL PRIMARY KEY,
	product_id integer NOT NULL REFERENCES public.product(id) ON DELETE CASCADE,
	name varchar(250)
	)`
	_, err = sqlDB.Query(sqlProductImages)

	return
}

// sql-s
var (
	// select
	getAllProducts      = "SELECT * FROM public.product"
	getProductById      = "SELECT * FROM public.product p WHERE p.id = $1"
	getAllProductImages = "SELECT * FROM public.images i WHERE i.product_id = $1"

	// update
	updateProduct = "UPDATE public.product " +
		"SET price=$1, name=$2, description=$3, imagepath=$4 " +
		"WHERE id = $5"

	// delete
	deleteProduct = "DELETE FROM public.product WHERE id=$1 RETURNING id"
	deleteImage   = "DELETE FROM public.image WHERE id=$1 RETURNING name"

	// create
	createProduct = "INSERT INTO public.product (price, name, description) " +
		"VALUES ($1, $2, $3) RETURNING id"
	createImage = "INSERT INTO public.images (product_id, name) " +
		"VALUES ($1, $2)"
)

var (
	ProductNotCreatedError = errors.New("Product not created")
	NoSuchIdProductError   = errors.New("Product not found by given ID")
	ImageNotCreatedError   = errors.New("Image not created")
)

func GetAllProducts() (p []Product, err error) {
	var productRes *sql.Rows
	productRes, err = sqlDB.Query(getAllProducts)
	if err != nil {
		return
	}

	for productRes.Next() {
		tempProduct := Product{
			Images: []Images{},
		}

		err = productRes.Scan(&tempProduct.Id, &tempProduct.Price, &tempProduct.Name, &tempProduct.Description)
		if err != nil {
			return
		}

		// find all images for product
		var imageRes *sql.Rows
		imageRes, err = sqlDB.Query(getAllProductImages, tempProduct.Id)
		if err != nil {
			return
		}

		for imageRes.Next() {
			tempImg := Images{}
			err = imageRes.Scan(&tempImg.Id, &tempImg.ProductId, &tempImg.Name)
			if err != nil {
				return
			}

			tempProduct.Images = append(tempProduct.Images, tempImg)
		}

		// append all to the product
		p = append(p, tempProduct)
	}

	return
}

func GetProductById(id int) (p Product, err error) {
	var res *sql.Row
	res = sqlDB.QueryRow(getProductById, id)
	err = res.Scan(&p.Id, &p.Price, &p.Name, &p.Description)
	if err != nil {
		return
	}

	return
}

func DeleteProductById(id int) (productId int, err error) {
	res := sqlDB.QueryRow(deleteProduct, id)
	err = res.Scan(&productId)
	if err != nil {
		return
	}

	return
}

func UpdateProductById(id int, p Product) error {
	res, err := sqlDB.Exec(updateProduct, p.Price, p.Name, p.Description, id)
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

func CreateProduct(p Product) (productId int, err error) {
	row := sqlDB.QueryRow(createProduct, p.Price, p.Name, p.Description)

	err = row.Scan(&productId)

	return
}

func InsertImages(productId int, imageNames []string) error {
	for _, name := range imageNames {
		res, err := sqlDB.Exec(createImage, productId, name)
		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rows == 0 {
			return ImageNotCreatedError
		}
	}

	return nil
}

func DeleteImage(imageId int) (imageName string, err error) {
	row := sqlDB.QueryRow(deleteImage, imageId)

	err = row.Scan(&imageName)
	if err != nil {
		return
	}

	return
}

package models

import (
	"database/sql"
	"errors"
)

type productsRepo struct {
	db *sql.DB
}

func newProductsRepo(sqlDB *sql.DB) newsRepo {
	return productsRepo{db: sqlDB}
}

type Images struct {
	Id        int    `db:"id"`
	ProductId int    `db:"product_id"`
	Name      string `db:"name"`
}

type Product struct {
	Id          int    `db:"id"`
	Price       int    `db:"price"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Images      []Images
}

var (
	ProductNotCreatedError = errors.New("Product not created")
	NoSuchIdProductError   = errors.New("Product not found by given ID")
	ImageNotCreatedError   = errors.New("Image not created")
)

func (p *productsRepo) Get() (p []Product, err error) {
	var productRes *sql.Rows
	productRes, err = p.db.Query("SELECT * FROM public.product")
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
		imageRes, err = p.db.Query("SELECT * FROM public.images i WHERE i.product_id = $1",
			tempProduct.Id)
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

func (p *productsRepo) GetById(id int) (p Product, err error) {
	var res *sql.Row
	res = p.db.QueryRow("SELECT * FROM public.product p WHERE p.id = $1", id)
	err = res.Scan(&p.Id, &p.Price, &p.Name, &p.Description)
	if err != nil {
		return
	}

	return
}

func (p *productsRepo) DeleteById(id int) (productId int, err error) {
	res := p.db.QueryRow("DELETE FROM public.product WHERE id=$1 RETURNING id", id)
	err = res.Scan(&productId)
	if err != nil {
		return
	}

	return
}

func (p *productsRepo) UpdateById(id int, p Product) error {
	res, err := p.db.Exec(
		`
		UPDATE public.product
		SET price=$1, name=$2, description=$3, imagepath=$4
		WHERE id = $5
		`,
		p.Price,
		p.Name,
		p.Description,
		id)
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

func (p *productsRepo) Create(p Product) (productId int, err error) {
	row := p.db.QueryRow(
		`
		INSERT INTO public.product (price, name, description)
		VALUES ($1, $2, $3) RETURNING id
		`,
		p.Price,
		p.Name,
		p.Description)

	err = row.Scan(&productId)

	return
}

func (p *productsRepo) InsertImages(productId int, imageNames []string) error {
	for _, name := range imageNames {
		res, err := p.db.Exec(
			`
			INSERT INTO public.images (product_id, name)
			VALUES ($1, $2)
			`,
			productId,
			name)
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

func (p *productsRepo) DeleteImage(imageId int) (imageName string, err error) {
	row := p.db.QueryRow("DELETE FROM public.image WHERE id=$1 RETURNING name", imageId)

	err = row.Scan(&imageName)
	if err != nil {
		return
	}

	return
}

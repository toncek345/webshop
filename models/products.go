package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type productsRepo struct {
	db *sqlx.DB
}

func newProductsRepo(sqlDB *sqlx.DB) productsRepo {
	return productsRepo{db: sqlDB}
}

type Product struct {
	ID          int     `db:"id"`
	Price       float64 `db:"price"`
	Name        string  `db:"name"`
	Description string  `db:"description"`

	Images []*Image `db:"-"`
}

func (pr *productsRepo) Get() ([]Product, error) {
	var products []Product
	if err := pr.db.Select(
		&products,
		"SELECT * FROM products"); err != nil {
		return nil, fmt.Errorf("models/products: error getting products: %w", err)
	}

	productIDs := make([]int, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}

	var images []*Image
	if err := pr.db.Select(
		&images,
		`SELECT * FROM images WHERE product_id IN ($1)`,
		productIDs,
	); err != nil {
		return nil, fmt.Errorf("models/products: error getting images for products: %w", err)
	}

	for _, i := range images {
		for _, p := range products {
			if i.ProductID == p.ID {
				p.Images = append(p.Images, i)
			}
		}
	}

	return products, nil
}

func (pr *productsRepo) GetByID(id int) (Product, error) {
	var p Product
	if err := pr.db.Get(&p, "SELECT * FROM products where id = $1", id); err != nil {
		return p, fmt.Errorf("models/products: error getting product by id: %w", err)
	}

	if err := pr.db.Select(
		&p.Images,
		`SELECT * FROM images
		WHERE image.product_id = $1`,
		id); err != nil {
		return p, fmt.Errorf("models/products: error getting images for product: %w", err)
	}

	return p, nil
}

func (pr *productsRepo) DeleteByID(id int) error {
	// TODO: delete images also..
	if _, err := pr.db.Exec("DELETE FROM products WHERE id=$1", id); err != nil {
		return fmt.Errorf("models/products: error deleting product: %w", err)
	}

	return nil
}

func (pr *productsRepo) UpdateByID(id int, p Product) error {
	if _, err := pr.db.Exec(
		`
		UPDATE products
		SET price = $1, name = $2, description = $3
		WHERE id = $4
		`,
		p.Price,
		p.Name,
		p.Description,
		id); err != nil {
		return fmt.Errorf("models/products: error updating by id: %w", err)
	}

	return nil
}

func (pr *productsRepo) Create(p Product) (int, error) {
	lastID := struct {
		LastID int `db:"id"`
	}{}

	if err := pr.db.Get(
		&lastID,
		`
		INSERT INTO products (price, name, description)
		VALUES ($1, $2, $3) RETURNING id
		`,
		p.Price,
		p.Name,
		p.Description,
	); err != nil {
		return 0, fmt.Errorf("models/products: error creating product: %w", err)
	}

	return lastID.LastID, nil
}

func (pr *productsRepo) InsertImages(productID int, imageKeys []string) error {
	for _, key := range imageKeys {
		lastID := struct {
			LastID int `db:"id"`
		}{}

		if err := pr.db.Get(
			&lastID,
			`INSERT INTO images (key, product_id) VALUES ($1, $2) RETURNING id`,
			key, productID); err != nil {
			return fmt.Errorf("models/products: error inserting image: %w", err)
		}
	}

	return nil
}

func (pr *productsRepo) DeleteImage(imageID int) (string, error) {
	imageKey := struct {
		ImageKey string `db:"name"`
	}{}
	if err := pr.db.Get(
		&imageKey,
		"DELETE FROM images WHERE id = $1 RETURNING key",
		imageID,
	); err != nil {
		return "", fmt.Errorf("models/products: error deleting image: %w", err)
	}

	return imageKey.ImageKey, nil
}

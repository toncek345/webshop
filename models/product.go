package models

import (
	"errors"
)

type Product struct {
	Id          int
	Price       float64
	Name        string
	Description string
}

// database mock
var (
	products = []Product{
		Product{
			Id:          1,
			Price:       2.45,
			Name:        "majica",
			Description: "zakon stvar koja te grije gore",
		},
		Product{
			Id:          2,
			Price:       5.45,
			Name:        "lače",
			Description: "zakon stvar koja te grije dole",
		},
	}
	nextProductId = 3
)

var (
	NoSuchIdProductError = errors.New("Product not found by given ID")
)

func GetAllProducts() []Product {
	return products
}

func GetProductById(id int) (Product, error) {
	for _, v := range products {
		if v.Id == id {
			return v, nil
		}
	}
	return Product{}, NoSuchIdProductError
}

func DeleteProductById(id int) error {
	for i, v := range products {
		if v.Id == id {
			products = append(products[:i], products[i+1:]...)
			return nil
		}
	}
	return NoSuchIdProductError
}

func UpdateProductById(id int) error {
	panic("unimplemented")
}

// product object without the id field
func CreateProduct(p Product) bool {
	p.Id = nextProductId
	nextProductId++
	products = append(products, p)
	return true
}

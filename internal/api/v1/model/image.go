package model

type Image struct {
	ID        int    `db:"id" json:"id"`
	Path      string `db:"path" json:"path"`
	ProductID int    `db:"product_id" json:"product_id"`
}

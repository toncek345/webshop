package model

type Image struct {
	ID        int    `db:"id"`
	Key       string `db:"key"`
	ProductID int    `db:"product_id"`
}

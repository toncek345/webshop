package models

import (
	"database/sql"
	"errors"
)

type newsRepo struct {
	db *sql.DB
}

func newNewsRepo(sqlDB *sql.DB) newsRepo {
	return newsRepo{db: sqlDB}
}

type News struct {
	ID        int    `db:"id"`
	Header    string `db:"header"`
	Text      string `db:"text"`
	ImagePath string `db:"image_path"`
}

var (
	ErrNewsNotCreatedError = errors.New("News not created")
	ErrNoSuchIdNewsError   = errors.New("News not found by given ID")
)

func (n *newsRepo) Get() (n []News, err error) {
	var res *sql.Rows
	res, err = n.db.Query("SELECT * FROM public.news")
	if err != nil {
		return
	}

	for res.Next() {
		temp := News{}
		err = res.Scan(&temp.ID, &temp.Header, &temp.Text,
			&temp.ImagePath)
		if err != nil {
			return
		}

		n = append(n, temp)
	}

	return
}

func (n *newsRepo) GetById(id int) (n News, err error) {
	var res *sql.Row
	res = n.db.QueryRow("SELECT * FROM public.news n WHERE n.id = $1", id)

	err = res.Scan(&n.ID, &n.Header, &n.Text, &n.ImagePath)
	if err != nil {
		return
	}

	return
}

func (n *newsRepo) DeleteById(id int) (n News, err error) {
	row := n.db.QueryRow("DELETE FROM public.news WHERE id=$1 RETURNING *", id)

	err = row.Scan(&n.ID, &n.Header, &n.Text, &n.ImagePath)

	return
}

func (n *newsRepo) UpdateById(id int, n News) error {
	res, err := n.db.Exec(
		`
		UPDATE public.news
		SET header=$1, text=$2, imagepath=$3
		WHERE id=$4
		`,
		n.Header,
		n.Text,
		n.ImagePath,
		id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoSuchIdNewsError
	}

	return nil
}

func (n *newsRepo) CreateNews(n News) error {
	res, err := n.db.Exec(
		`
		INSERT INTO public.news (header, text, imagepath)
		VALUES ($1, $2, $3)
		`,
		n.Header,
		n.Text,
		n.ImagePath)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNewsNotCreatedError
	}

	return nil
}

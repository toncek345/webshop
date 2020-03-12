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

func (nr *newsRepo) Get() (n []News, err error) {
	var res *sql.Rows
	res, err = nr.db.Query("SELECT * FROM public.news")
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

func (nr *newsRepo) GetById(id int) (n News, err error) {
	var res *sql.Row
	res = nr.db.QueryRow("SELECT * FROM public.news n WHERE n.id = $1", id)

	err = res.Scan(&n.ID, &n.Header, &n.Text, &n.ImagePath)
	if err != nil {
		return
	}

	return
}

func (nr *newsRepo) DeleteById(id int) (n News, err error) {
	row := nr.db.QueryRow("DELETE FROM public.news WHERE id=$1 RETURNING *", id)

	err = row.Scan(&n.ID, &n.Header, &n.Text, &n.ImagePath)

	return
}

func (nr *newsRepo) UpdateById(id int, n News) error {
	res, err := nr.db.Exec(
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

func (nr *newsRepo) CreateNews(n News) error {
	res, err := nr.db.Exec(
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

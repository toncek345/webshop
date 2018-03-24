package models

import (
	"database/sql"
	"errors"
)

type News struct {
	Id        int
	Header    string
	Text      string
	ImagePath string // path to image
}

func initNews() (err error) {
	sql := `CREATE TABLE public.news (
		id serial NOT NULL PRIMARY KEY,
		header varchar(250),
		text text,
		imagepath varchar(250)
	)`

	_, err = sqlDB.Exec(sql)
	return
}

// sql-s
var (
	// select
	getAllNews  = "SELECT * FROM public.news"
	getNewsById = "SELECT * FROM public.news n WHERE n.id = $1"

	// update
	updateNews = "UPDATE public.news " +
		"SET header=$1, text=$2, imagepath=$3" +
		"WHERE id=$4"

	// delete
	deleteNews = "DELETE FROM public.news WHERE id=$1"

	// create
	createNews = "INSERT INTO public.news (header, text, imagepath)" +
		"VALUES ($1, $2, $3)"
)

// errors
var (
	NewsNotCreatedError = errors.New("News not created")
	NoSuchIdNewsError   = errors.New("News not found by given ID")
)

func GetAllNews() (n []News, err error) {
	var res *sql.Rows
	res, err = sqlDB.Query(getAllNews)
	if err != nil {
		return
	}

	for res.Next() {
		temp := News{}
		err = res.Scan(&temp.Id, &temp.Header, &temp.Text,
			&temp.ImagePath)
		if err != nil {
			return
		}

		n = append(n, temp)
	}

	return
}

func GetNewsById(id int) (n News, err error) {
	var res *sql.Row
	res = sqlDB.QueryRow(getNewsById, id)
	err = res.Scan(&n.Id, &n.Header, &n.Text, &n.ImagePath)
	if err != nil {
		return
	}

	return
}

func DeleteNewsById(id int) error {
	res, err := sqlDB.Exec(deleteNews, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return NoSuchIdNewsError
	}

	return nil
}

func UpdateNewsById(id int, n News) error {
	res, err := sqlDB.Exec(updateNews, n.Header, n.Text, n.ImagePath, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return NoSuchIdNewsError
	}

	return nil
}

func CreateNews(n News) error {
	res, err := sqlDB.Exec(createNews, n.Header, n.Text, n.ImagePath)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return NewsNotCreatedError
	}

	return nil
}

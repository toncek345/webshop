package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type newsRepo struct {
	db *sqlx.DB
}

func newNewsRepo(sqlDB *sqlx.DB) newsRepo {
	return newsRepo{db: sqlDB}
}

type News struct {
	ID        int    `db:"id"`
	Header    string `db:"header"`
	Text      string `db:"text"`
	ImagePath string `db:"image_path"`
}

func (nr *newsRepo) Get() ([]News, error) {
	var news []News
	if err := nr.db.Select(&news, "SELECT * FROM news"); err != nil {
		return nil, fmt.Errorf("models/news: error getting news: %w", err)
	}

	return news, nil
}

func (nr *newsRepo) GetByID(id int) (News, error) {
	var n News
	if err := nr.db.Get(
		&n,
		"SELECT * FROM news WHERE id = $1",
		id,
	); err != nil {
		return n, fmt.Errorf("models/news: error getting news by id: %w", err)
	}

	return n, nil
}

func (nr *newsRepo) DeleteByID(id int) error {
	if _, err := nr.db.Exec(
		"DELETE FROM news WHERE id = $1",
		id); err != nil {
		return fmt.Errorf("models/news: error deleting news: %w", err)
	}

	return nil
}

func (nr *newsRepo) UpdateByID(id int, n News) error {
	if _, err := nr.db.Exec(
		`UPDATE news SET header = $1, text = $2, imagepath = $3`,
		n.Header, n.Text, n.ImagePath, id); err != nil {
		return fmt.Errorf("models/news: error updating news: %w", err)
	}

	return nil
}

func (nr *newsRepo) CreateNews(n News) error {
	if _, err := nr.db.Exec(
		`
		INSERT INTO public.news (header, text, imagepath)
		VALUES ($1, $2, $3)
		`,
		n.Header,
		n.Text,
		n.ImagePath); err != nil {
		return fmt.Errorf("models/news: error creating news: %w", err)
	}

	return nil
}

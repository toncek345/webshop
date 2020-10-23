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
	ID     int    `db:"id"`
	Header string `db:"header"`
	Text   string `db:"text"`

	Images []*Image `db:"-"`
}

func (nr *newsRepo) Get() ([]News, error) {
	var news []News
	if err := nr.db.Select(&news, "SELECT * FROM news"); err != nil {
		return nil, fmt.Errorf("models/news: error getting news: %w", err)
	}

	newsIDs := make([]int, 0, len(news))
	for _, n := range news {
		newsIDs = append(newsIDs, n.ID)
	}

	var images []struct {
		NewsID   int    `db:"news_id"`
		ImageID  int    `db:"image_id"`
		ImageKey string `db:"image_key"`
	}
	if err := nr.db.Select(
		&images,
		`SELECT *
		FROM images
		JOIN news_images ON news_images.image_id = images.id
		WHERE news_images.id IN ($1)`,
		newsIDs); err != nil {
		return nil, fmt.Errorf("models/news: loading images: %w", err)
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

	if err := nr.db.Select(
		&n.Images,
		`SELECT images.id AS 'id', images.key AS 'key' FROM images
		JOIN news_images ON news_images.news_id = images.id
		WHERE news_images.news_id = $1`,
		id); err != nil {
		return n, fmt.Errorf("models/news: getting image for news: %w", err)
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
		`UPDATE news SET header = $1, text = $2 WHERE id = $3`,
		n.Header, n.Text, id); err != nil {
		return fmt.Errorf("models/news: error updating news: %w", err)
	}

	return nil
}

func (nr *newsRepo) CreateNews(n News, image Image) error {
	if _, err := nr.db.Exec(
		`
			DO $$
			DECLARE
				newsId bigint;
				imageId bigint;
			BEGIN
				INSERT INTO news (header, text) VALUES ($1, $2) RETURNING id INTO newsId;
				INSERT INTO images (key) VALUES ($3) RETURNING id INTO imageId;
				INSERT INTO news_images (news_id, image_id) VALUES (newsId, imageId);
			END $$;
			`,
		n.Header, n.Text, image.Key); err != nil {
		return fmt.Errorf("models/news: creating news: %w", err)
	}

	return nil
}

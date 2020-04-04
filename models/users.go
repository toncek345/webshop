package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type usersRepo struct {
	db *sqlx.DB
}

func newUserRepo(sqlDB *sqlx.DB) usersRepo {
	return usersRepo{db: sqlDB}
}

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func (ur *usersRepo) GetByUsername(username string) (User, error) {
	var u User
	if err := ur.db.Get(
		&u,
		"SELECT * FROM user WHERE username = $1",
		username); err != nil {
		return u, fmt.Errorf("models/users: error getting user by username: %w", err)
	}

	return u, nil
}

func (ur *usersRepo) UpdateUser(id int, u User) error {
	if _, err := ur.db.Exec(
		`
		UPDATE public.user
		SET username=$1, public.user.password=$2
		WHERE id=$3
		`,
		u.Username,
		u.Password,
		id); err != nil {
		return fmt.Errorf("models/users: error updating user: %w", err)
	}

	return nil
}

func (ur *usersRepo) CreateUser(u User) error {
	if _, err := ur.db.Exec(
		`INSERT INTO public.user (username, password)
		VALUES ($1, $2)`,
		u.Username,
		u.Password); err != nil {
		return fmt.Errorf("models/users: error creating user: %w", err)
	}

	return nil
}

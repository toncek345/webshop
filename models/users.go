package models

import (
	"database/sql"
	"errors"
)

type usersRepo struct {
	db *sql.DB
}

func newUserRepo(sqlDB *sql.DB) usersRepo {
	return usersRepo{db: sqlDB}
}

// errors
var (
	UserNoMatch           = errors.New("Username or password does not match")
	AuthNotAuthenticates  = errors.New("User not authenticated")
	UserNotFoundByIdError = errors.New("User not found by id")
	UserNotCreatedError   = errors.New("User not created")
)

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func (ur *usersRepo) GetByUsername(username string) (u User, err error) {
	var res *sql.Row
	res = ur.db.QueryRow("SELECT * FROM public.user WHERE username = $1", username)
	err = res.Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			err = UserNoMatch
			return
		}
		return
	}

	return
}

func (ur *usersRepo) UpdateUser(id int, u User) error {
	res, err := ur.db.Exec(
		`
		UPDATE public.user
		SET username=$1, public.user.password=$2
		WHERE id=$3
		`,
		u.Username,
		u.Password,
		id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return UserNotFoundByIdError
	}

	return nil
}

func (ur *usersRepo) CreateUser(u User) error {
	res, err := ur.db.Exec(
		`INSERT INTO public.user (username, password)
		VALUES ($1, $2)`,
		u.Username,
		u.Password)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return UserNotCreatedError
	}

	return nil
}

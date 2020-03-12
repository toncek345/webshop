package models

import (
	"database/sql"
	"errors"
)

type usersRepo struct {
	db *sql.DB
}

func newUserRepo(sqlDB *sql.DB) newsRepo {
	return authRepo{db: sqlDB}
}

// errors
var (
	UserNoMatch           = errors.New("Username or password does not match")
	AuthNotAuthenticates  = errors.New("User not authenticated")
	UserNotFoundByIdError = errors.New("User not found by id")
	UserNotCreatedError   = errors.New("User not created")
)

// TODO: remove after fixing handlers
// getUserById                  = "SELECT * FROM public.user WHERE id=$1"

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

// TODO: remove after fixing handlers
// func findUser(username string) (u User, err error) {
// 	var res *sql.Row
// 	res = sqlDB.QueryRow(getUserByPasswordAndUsername, username)
// 	err = res.Scan(&u.Id, &u.Username, &u.Password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			err = UserNoMatch
// 			return
// 		}
// 		return
// 	}

// 	return
// }

func UpdateUser(id int, u User) error {
	res, err := sqlDB.Exec(
		`
		UPDATE public.user
		SET username=$1, public.user.password=$2
		WHERE id=$3
		`,
		updateUser,
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

func CreateUser(u User) error {
	res, err := sqlDB.Exec(
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

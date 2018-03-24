package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var sqlDB *sql.DB

func DbConnect(dbString string) (err error) {
	sqlDB, err = sql.Open("postgres", dbString)
	if err != nil {
		return
	}
	return nil
}

func InitDb() error {
	if err := initNews(); err != nil {
		return err
	}

	if err := initProduct(); err != nil {
		return err
	}

	if err := initUser(); err != nil {
		return err
	}

	if err := initAuth(); err != nil {
		return err
	}

	return nil
}

// getting global db object
func GetDb() *sql.DB {
	return sqlDB
}

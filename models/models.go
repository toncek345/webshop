package models

import (
	"database/sql"
)

var sqlDB *sql.DB

// TODO: db connection hardcoded
func init() {
	// connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full" // TODO

	// var err error
	// sqlDB, err = sql.Open("postgres", connStr)
	// if err != nil {
	// 	panic(err)
	// }
}

// getting global db object
func GetDb() *sql.DB {
	return sqlDB
}

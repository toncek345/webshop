package main

import (
	"flag"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/toncek345/webshop/config"
)

var (
	environment = flag.String(
		"env",
		"development",
		"development, production running mode",
	)
)

func main() {
	flag.Parse()

	var env config.Environment
	switch *environment {
	case "development":
		env = config.DevEnv
	case "production":
		env = config.ProdEnv
	}

	config, err := config.New(env)
	if err != nil {
		panic(err)
	}

	m, err := migrate.New(
		"file://migrations",
		config.PGSQLConnString())
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		panic(err)
	}
}

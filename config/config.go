package config

import (
	"fmt"
	"io/ioutil"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Env Environment

	// TODO: static and storage path are currently the same thing.
	StaticPath      string `yaml:"static_path"`
	DiskStoragePath string `yaml:"disk_storage_path"`

	WebPort int64 `yaml:"web_port"`

	Postgres struct {
		Database string `yaml:"database"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"postgres"`
}

type Environment int

const (
	DevEnv = iota
	ProdEnv
)

func New(env Environment) (Config, error) {
	var cfg Config
	var fileName string
	if env == DevEnv {
		fileName = "config/development.yml"
	} else {
		fileName = "config/production.yml"
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cfg, fmt.Errorf("config: error reading file: %w", err)
	}

	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return cfg, fmt.Errorf("config: error unmarshalling yaml: %w", err)
	}

	cfg.Env = env

	return cfg, nil
}

func (cfg *Config) PGSQLConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}

func (cfg *Config) SqlxConn() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", cfg.PGSQLConnString())
}

func (cfg *Config) MigrationObj() (*migrate.Migrate, error) {
	return migrate.New(
		"file://db/migrations",
		cfg.PGSQLConnString())
}

func (cfg *Config) Fixture() (*testfixtures.Loader, error) {
	db, err := cfg.SqlxConn()
	if err != nil {
		return nil, err
	}

	return testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(cfg.fixturePath()),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.SkipResetSequences(),
	)
}

func (cfg *Config) fixturePath() string {
	if cfg.Env == ProdEnv {
		panic("no fixtures on prod env")
	}

	return "db/fixtures/dev"
}

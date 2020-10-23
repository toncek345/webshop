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
	"github.com/toncek345/webshop/internal/pkg/storage"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Env Environment `yaml:"-"`

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

type Environment string

const (
	DevEnv  Environment = "development"
	ProdEnv Environment = "production"
)

func New(env Environment) (*Config, error) {
	cfg := &Config{Env: env}
	var fileName string

	switch env {
	case DevEnv:
		fileName = "config/development.yml"
	case ProdEnv:
		fileName = "config/production.yml"
	default:
		return nil, fmt.Errorf("%s env doesn't exist", env)
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("config: error reading file: %w", err)
	}

	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, fmt.Errorf("config: error unmarshalling yaml: %w", err)
	}

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

func (cfg *Config) Storage() (storage.Storage, error) {
	return &storage.Disk{DirPath: cfg.DiskStoragePath}, nil
}

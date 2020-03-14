package config

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	StaticPath string `yaml:"static_path"`
	WebPort    int64  `yaml:"web_port"`

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

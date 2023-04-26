package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Cfg struct {
		Port string `env:"PORT"`
		DB   DB
	}

	DB struct {
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		Name     string `env:"DB_NAME"`
		Username string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
	}

	option func(*Cfg)
)

func New() *Cfg {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "development" {
		return load(
			readFrom(
				"./config/development.env",
			),
		)
	}

	if appEnv == "production" {
		return load(
			readFrom(
				"./config/production.env",
			),
		)
	}

	return nil
}

func load(opts ...option) *Cfg {
	var cfg Cfg
	for _, option := range opts {
		option(&cfg)
	}

	return &cfg
}

func readFrom(paths ...string) option {
	return func(c *Cfg) {
		for _, path := range paths {
			err := cleanenv.ReadConfig(path, c)
			if err != nil {
				panic(err)
			}
		}
	}
}

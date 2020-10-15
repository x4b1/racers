package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/xabi93/racers/internal/storage/postgres"
)

type Env string

const (
	Local Env = "local"
	Test  Env = "test"
	Prod  Env = "production"
)

type Conf struct {
	Env      Env    `env:"ENVIRONMENT" envDefault:"local"`
	Port     string `env:"PORT" envDefault:"8080"`
	Postgres postgres.Config
}

func Load() (Conf, error) {
	c := Conf{}

	if err := env.Parse(&c); err != nil {
		return Conf{}, err
	}

	return c, nil
}

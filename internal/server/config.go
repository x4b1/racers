package server

import (
	"github.com/xabi93/racers/internal/storage/postgres"

	"github.com/caarlos0/env/v6"
)

type Conf struct {
	Port     string `env:"PORT" envDefault:"8080"`
	Postgres postgres.Config
}

func LoadConf() (Conf, error) {
	var c Conf

	return c, env.Parse(&c)
}

package graph

import (
	"github.com/xabi93/racers/internal/service"
)

//go:generate go run github.com/99designs/gqlgen

func New(races service.Races) Config {
	return Config{Resolvers: &Resolver{races}}
}

type Resolver struct {
	racers service.Races
}

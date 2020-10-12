package graph

import (
	"github.com/xabi93/racers/internal/service"
	"github.com/xabi93/racers/internal/storage/postgres/ent"
	// "github.com/xabi93/racers/internal/storage/postgres/ent"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func New(s *service.Service, cli *ent.Client) Config {
	return Config{Resolvers: &Resolver{s, cli}}
}

type Resolver struct {
	*service.Service
	query *ent.Client
}

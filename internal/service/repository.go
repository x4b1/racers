package service

import (
	"context"

	racers "github.com/xabi93/racers/internal"
)

//go:generate moq -stub -pkg service_test -out mock_repository_test.go . RacesRepository TeamsRepository UsersGetter

type RacesRepository interface {
	RacesGetter
	Exists(ctx context.Context, race racers.Race) (bool, error)
	Save(ctx context.Context, race racers.Race) error
}

type RacesGetter interface {
	All(ctx context.Context) ([]racers.Race, error)
	Get(ctx context.Context, id racers.RaceID) (racers.Race, error)
}

type TeamsRepository interface {
	TeamsGetter
	Save(ctx context.Context, team racers.Team) error
}

type TeamsGetter interface {
	ByMember(ctx context.Context, id racers.UserID) (*racers.Team, error)
	Get(ctx context.Context, id racers.TeamID) (racers.Team, error)
}

type UsersGetter interface {
	Get(ctx context.Context, id racers.UserID) (racers.User, error)
	Current(ctx context.Context) racers.User
}

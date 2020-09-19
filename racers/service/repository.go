package service

import (
	"context"

	"github.com/xabi93/go-clean/racers"
)

//go:generate moq -pkg service_test -out mock_repository_test.go . RacesRepository UsersRepository TeamsRepository

type RacesRepository interface {
	RacesGetter
	Save(ctx context.Context, race racers.Race) error
}

type RacesGetter interface {
	ByID(ctx context.Context, id racers.RaceID) (*racers.Race, error)
}

type UsersRepository interface {
	UsersGetter
	Save(ctx context.Context, user racers.User) error
}

type UsersGetter interface {
	ByID(ctx context.Context, id racers.UserID) (*racers.User, error)
}

type TeamsRepository interface {
	TeamsGetter
	Save(ctx context.Context, team racers.Team) error
}

type TeamsGetter interface {
	ByID(ctx context.Context, id racers.TeamID) (*racers.Team, error)
	ByMember(ctx context.Context, memberID racers.UserID) (*racers.Team, error)
}

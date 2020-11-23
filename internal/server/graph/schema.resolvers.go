package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	racers "github.com/xabi93/racers/internal"
	errorsx "github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/service"
)

func (r *mutationResolver) CreateRace(ctx context.Context, race RaceInput) (CreateRaceResult, error) {
	result, err := r.racers.Create(ctx, service.CreateRace{
		ID:   race.ID,
		Name: race.Name,
		Date: race.Date,
	})

	var (
		invalidID   racers.InvalidRaceIDError
		invalidName racers.InvalidRaceNameError
		invalidDate racers.InvalidRaceDateError
	)
	if err != nil {
		switch {
		case errorsx.As(err, &invalidID):
			return InvalidIDError{Message: invalidID.Error()}, nil
		case errorsx.As(err, &invalidName):
			return InvalidRaceNameError{Message: invalidName.Error()}, nil
		case errorsx.As(err, &invalidDate):
			return InvalidRaceDateError{Message: invalidDate.Error()}, nil
		case errorsx.Is(err, service.ErrRaceAlreadyExists):
			return RaceAlreadyExists{Message: err.Error()}, nil
		}
		return nil, NewInternalError()
	}

	return NewRace(result), err
}

func (r *queryResolver) Race(ctx context.Context, id string) (RaceResult, error) {
	result, err := r.racers.Get(ctx, service.GetRace{ID: id})

	var invalidID racers.InvalidRaceIDError
	if err != nil {
		switch {
		case errorsx.As(err, &invalidID):
			return InvalidIDError{Message: invalidID.Error()}, nil
		case errorsx.Is(err, service.ErrRaceNotFound):
			return RaceNotFound{Message: err.Error()}, nil
		}

		return nil, NewInternalError()
	}

	return NewRace(result), err
}

func (r *queryResolver) Races(ctx context.Context) (*Races, error) {
	races, err := r.racers.List(ctx)
	result := make([]*Race, len(races))
	for i, r := range races {
		result[i] = NewRace(r)
	}
	return &Races{Races: result}, err
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

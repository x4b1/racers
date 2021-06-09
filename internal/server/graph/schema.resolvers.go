package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	racers "github.com/xabi93/racers/internal"
	errorsx "github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/server/graph/models"
	"github.com/xabi93/racers/internal/service"
)

func (r *mutationResolver) CreateRace(ctx context.Context, race models.RaceInput) (models.CreateRaceResult, error) {
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
			return models.InvalidIDError{Message: invalidID.Error()}, nil
		case errorsx.As(err, &invalidName):
			return models.InvalidRaceNameError{Message: invalidName.Error()}, nil
		case errorsx.As(err, &invalidDate):
			return models.InvalidRaceDateError{Message: invalidDate.Error()}, nil
		case errorsx.Is(err, service.ErrRaceAlreadyExists):
			return models.RaceAlreadyExists{Message: err.Error()}, nil
		}
		return nil, models.NewInternalError()
	}

	return models.NewRace(result), err
}

func (r *queryResolver) Race(ctx context.Context, id string) (models.RaceResult, error) {
	result, err := r.racers.Get(ctx, service.GetRace{ID: id})

	var invalidID racers.InvalidRaceIDError
	if err != nil {
		switch {
		case errorsx.As(err, &invalidID):
			return models.InvalidIDError{Message: invalidID.Error()}, nil
		case errorsx.Is(err, service.ErrRaceNotFound):
			return models.RaceNotFound{Message: err.Error()}, nil
		}

		return nil, models.NewInternalError()
	}

	return models.NewRace(result), err
}

func (r *queryResolver) Races(ctx context.Context) (*models.Races, error) {
	races, err := r.racers.List(ctx)
	if err != nil {
		return nil, models.NewInternalError()
	}

	return models.NewRaces(races), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

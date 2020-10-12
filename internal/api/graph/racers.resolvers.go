package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/kr/pretty"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
	"github.com/xabi93/racers/internal/storage/postgres/ent"
	"github.com/xabi93/racers/internal/storage/postgres/ent/race"
)

func (r *mutationResolver) CreateRace(ctx context.Context, race RaceInput) (CreateRaceResult, error) {
	raceID := id.GenerateID().String()
	err := r.Races.Create(ctx, service.CreateRace{
		ID:   raceID,
		Name: race.Name,
		Date: race.Date,
	})

	var invalidName racers.InvalidRaceNameError
	var invalidDate racers.InvalidRaceDateError
	if err != nil {
		switch {
		case errors.As(err, &invalidName):
			pretty.Println(invalidName.Error())
			return InvalidNameError{Message: invalidName.Error()}, nil
		case errors.As(err, &invalidDate):
			pretty.Println(invalidDate.Error())
			return InvalidDateError{Message: invalidDate.Error()}, nil
		}
		return nil, err
	}

	result, err := r.query.Race.Get(ctx, raceID)
	if result != nil {
		return &Race{Race: *result}, err
	}

	return nil, err
}

func (r *queryResolver) Race(ctx context.Context, id string) (RaceResult, error) {
	race, err := r.query.Race.Query().Where(race.ID(id)).CollectFields(ctx, "Race").Only(ctx)
	if race != nil {
		return &Race{Race: *race}, nil
	}

	if ent.IsNotFound(err) {
		pretty.Println(id)
		return RaceByIDNotFound{
			ID: id,
		}, nil
	}

	return nil, err
}

func (r *userResolver) Races(ctx context.Context, obj *User) ([]*Race, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	userResolver     struct{ *Resolver }
)

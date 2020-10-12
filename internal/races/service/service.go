package service

import (
	"context"
	"time"

	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/races"
)

type Service interface {
	Create(ctx context.Context, r Create) error
	Join(ctx context.Context, r Join) error
}

type Races interface {
	Save(ctx context.Context, r races.Race) error
	ByID(ctx context.Context, id races.RaceID) (*races.Race, error)
}

type Competitors interface {
	ByID(ctx context.Context, id races.CompetitorID) (*races.Competitor, error)
}

type service struct {
	races       Races
	competitors Competitors
}

type Create struct {
	ID   string
	Name string
	Date time.Time
}

func (s service) Create(ctx context.Context, r Create) error {
	id, err := races.NewRaceID(r.ID)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}
	name, err := races.NewRaceName(r.Name)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}
	date, err := races.NewRaceDate(r.Date)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}

	return errors.WrapInternalError(s.races.Save(ctx, races.CreateRace(id, name, date)))
}

type Join struct {
	RaceID string
	UserID string
}

func (s service) Join(ctx context.Context, r Join) error {
	raceID, err := races.NewRaceID(r.RaceID)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}
	userID, err := races.NewUserID(r.UserID)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}

	race, err := s.races.ByID(ctx, raceID)
	if err != nil {
		return errors.WrapInternalError(err)
	}
	if race == nil {
		return errors.WrapNotFoundError(RaceByIDNotFoundError{raceID})
	}

	user, err := s.competitors.ByID(ctx, userID)
	if err != nil {
		return errors.WrapInternalError(err)
	}
	if user == nil {
		return errors.WrapNotFoundError(UserByIDNotFoundError{userID})
	}

	if err := race.Join(*user); err != nil {
		return err
	}

	if err := s.races.Save(ctx, *race); err != nil {
		return errors.WrapInternalError(err)
	}

	return nil
}

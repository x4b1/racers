package service

import (
	"context"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
)

const (
	CreateRaceName ActionName = "create race"
	GetRaceName    ActionName = "get race"
	JoinRaceName   ActionName = "join race"
)

type RacesService interface {
	Create(ctx context.Context, r CreateRace) error
	Get(ctx context.Context, r GetRace) (*racers.Race, error)
	Join(ctx context.Context, r JoinRace) error
}

func NewRacesService(races RacesRepository, users UsersGetter) RacesService {
	return racesService{races, users}
}

type racesService struct {
	races RacesRepository
	users UsersGetter
}

type CreateRace struct {
	ID   string
	Name string
	Date time.Time
}

func (rs racesService) Create(ctx context.Context, r CreateRace) error {
	id, err := racers.NewRaceID(r.ID)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}
	name, err := racers.NewRaceName(r.Name)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}
	date, err := racers.NewRaceDate(r.Date)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}

	race := racers.CreateRace(id, name, date)
	if err := rs.races.Save(ctx, race); err != nil {
		return errors.WrapInternalError(err)
	}

	return nil
}

type GetRace struct {
	RaceID string `json:"race_id,omitempty"`
}

func (s racesService) Get(ctx context.Context, r GetRace) (*racers.Race, error) {
	raceID, err := racers.NewRaceID(r.RaceID)
	if err != nil {
		return nil, errors.WrapWrongInputError(err)
	}

	race, err := s.races.ByID(ctx, raceID)
	if err != nil {
		return nil, errors.WrapInternalError(err)
	}
	if race == nil {
		return nil, errors.WrapNotFoundError(RaceByIDNotFoundError{raceID})
	}

	return race, nil
}

type JoinRace struct {
	RaceID string
	UserID string
}

func (s racesService) Join(ctx context.Context, r JoinRace) error {
	raceID, err := racers.NewRaceID(r.RaceID)
	if err != nil {
		return errors.WrapWrongInputError(err)
	}
	competitorID, err := racers.NewUserID(r.UserID)
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

	user, err := s.users.ByID(ctx, competitorID)
	if err != nil {
		return errors.WrapInternalError(err)
	}
	if user == nil {
		return errors.WrapNotFoundError(UserByIDNotFoundError{competitorID})
	}

	if err := race.Join(*user); err != nil {
		return err
	}

	if err := s.races.Save(ctx, *race); err != nil {
		return errors.WrapInternalError(err)
	}

	return nil
}

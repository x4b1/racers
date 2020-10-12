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
		return errors.WrapWrongInputError(err, CreateRaceName.String())
	}
	name, err := racers.NewRaceName(r.Name)
	if err != nil {
		return errors.WrapWrongInputError(err, CreateRaceName.String())
	}
	date, err := racers.NewRaceDate(r.Date)
	if err != nil {
		return errors.WrapWrongInputError(err, CreateRaceName.String())
	}

	race := racers.CreateRace(id, name, date)
	if err := rs.races.Save(ctx, race); err != nil {
		return errors.WrapInternalError(err, CreateRaceName.String())
	}

	return nil
}

type JoinRace struct {
	RaceID string
	UserID string
}

func (rs racesService) Join(ctx context.Context, r JoinRace) error {
	var ve errors.ValidationError

	raceID, err := racers.NewRaceID(r.RaceID)
	if err != nil {
		ve.Add(err)
	}
	userID, err := racers.NewUserID(r.UserID)
	if err != nil {
		ve.Add(err)
	}
	if err := ve.Valid(); err != nil {
		return err
	}

	race, err := rs.races.ByID(ctx, raceID)
	if err != nil {
		return errors.WrapInternalError(err, "getting race to join")
	}
	if race == nil {
		return errors.WrapNotFoundError(RaceByIDNotFoundError{raceID}, "getting race to join")
	}

	user, err := rs.users.ByID(ctx, userID)
	if err != nil {
		return errors.WrapInternalError(err, "getting user to join")
	}
	if user == nil {
		return errors.WrapNotFoundError(UserByIDNotFoundError{userID}, "getting user to join")
	}

	if err := race.Join(*user); err != nil {
		return err
	}

	if err := rs.races.Save(ctx, *race); err != nil {
		return errors.WrapInternalError(err, "saving race")
	}

	return nil
}

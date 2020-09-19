package service

import (
	"context"
	"time"

	"github.com/xabi93/go-clean/internal/errors"
	"github.com/xabi93/go-clean/racers"
)

func NewRacesService(races RacesRepository, users UsersGetter) RacesService {
	return RacesService{races, users}
}

type RacesService struct {
	races RacesRepository
	users UsersGetter
}

type CreateRace struct {
	ID   string
	Name string
	Date time.Time
}

func (rs RacesService) Create(ctx context.Context, r CreateRace) error {
	var ve errors.ValidationError

	id, err := racers.NewRaceID(r.ID)
	if err != nil {
		ve.Add(err)
	}
	name, err := racers.NewRaceName(r.Name)
	if err != nil {
		ve.Add(err)
	}
	date, err := racers.NewRaceDate(r.Date)
	if err != nil {
		ve.Add(err)
	}
	if err := ve.Valid(); err != nil {
		return err
	}

	if err := rs.races.Save(ctx, racers.CreateRace(id, name, date)); err != nil {
		return errors.WrapInternalError(err, "saving race")
	}

	return nil
}

type JoinRace struct {
	RaceID string
	UserID string
}

func (rs RacesService) Join(ctx context.Context, r JoinRace) error {
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

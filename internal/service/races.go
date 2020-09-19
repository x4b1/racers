package service

import (
	"context"
	"time"

	racers "github.com/xabi93/racers/internal"
)

func NewRaces(races RacesRepository, users UsersGetter, uow UnitOfWork, eb EventBus) Races {
	return Races{races, users, eb, uow}
}

type Races struct {
	races RacesRepository
	users UsersGetter
	eb    EventBus
	uow   UnitOfWork
}

type CreateRace struct {
	ID   string    `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
	Date time.Time `json:"date,omitempty"`
}

type RaceCreated struct {
	Race racers.Race `json:"race,omitempty"`
}

func (s Races) Create(ctx context.Context, r CreateRace) (race racers.Race, err error) {
	id, err := racers.NewRaceID(r.ID)
	if err != nil {
		return racers.Race{}, err
	}
	name, err := racers.NewRaceName(r.Name)
	if err != nil {
		return racers.Race{}, err
	}
	date, err := racers.NewRaceDate(r.Date)
	if err != nil {
		return racers.Race{}, err
	}

	race = racers.NewRace(id, name, date, s.users.Current(ctx).ID)

	exists, err := s.races.Exists(ctx, race)
	if err != nil {
		return racers.Race{}, err
	}
	if exists {
		return racers.Race{}, ErrRaceAlreadyExists
	}

	err = s.uow(ctx, func(ctx context.Context) error {
		if err := s.races.Save(ctx, race); err != nil {
			return err
		}

		return s.eb.Publish(ctx, newEvent(RaceCreated{Race: race}, s.users.Current(ctx).ID))
	})
	if err != nil {
		return racers.Race{}, err
	}

	return race, nil
}

type GetRace struct {
	ID string `json:"id,omitempty"`
}

func (s Races) Get(ctx context.Context, r GetRace) (racers.Race, error) {
	raceID, err := racers.NewRaceID(r.ID)
	if err != nil {
		return racers.Race{}, err
	}

	race, err := s.races.Get(ctx, raceID)
	if err != nil {
		return racers.Race{}, err
	}

	return race, nil
}

type JoinRace struct {
	RaceID string
	UserID string
}

type UserJoinedRace struct {
	User racers.User
	Race racers.Race
}

func (s Races) Join(ctx context.Context, r JoinRace) error {
	raceID, err := racers.NewRaceID(r.RaceID)
	if err != nil {
		return err
	}

	competitorID, err := racers.NewUserID(r.UserID)
	if err != nil {
		return err
	}

	race, err := s.races.Get(ctx, raceID)
	if err != nil {
		return err
	}

	user, err := s.users.Get(ctx, competitorID)
	if err != nil {
		return err
	}

	if err := race.Join(user); err != nil {
		return err
	}

	return s.uow(ctx, func(ctx context.Context) error {
		if err := s.races.Save(ctx, race); err != nil {
			return err
		}

		return s.eb.Publish(ctx, newEvent(UserJoinedRace{Race: race, User: user}, s.users.Current(ctx).ID))
	})
}

func (s Races) List(ctx context.Context) ([]racers.Race, error) {
	return s.races.All(ctx)
}

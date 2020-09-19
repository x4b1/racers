package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/matryer/is"
	"github.com/xabi93/go-clean/internal/errors"
	"github.com/xabi93/go-clean/internal/types"
	"github.com/xabi93/go-clean/racers"
	"github.com/xabi93/go-clean/racers/service"
)

type testRacesService struct {
	service service.RacesService
	races   *RacesRepositoryMock
	users   *UsersRepositoryMock
}

func newRacesServiceSuite() testRacesService {
	s := testRacesService{
		races: &RacesRepositoryMock{},
		users: &UsersRepositoryMock{},
	}

	s.service = service.NewRacesService(s.races, s.users)

	return s
}

func TestCreateRace(t *testing.T) {
	raceID := racers.RaceID(types.GenerateID())
	raceName := racers.RaceName("marathon")
	raceDate := racers.RaceDate(time.Now().AddDate(0, 0, 1))

	req := service.CreateRace{
		ID:   types.ID(raceID).String(),
		Name: string(raceName),
		Date: time.Time(raceDate),
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		err := s.service.Create(context.Background(), service.CreateRace{})

		var ve errors.ValidationError
		is.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidRaceIDError{},
			racers.InvalidRaceNameError{},
			racers.InvalidUserIDError{},
		}
		is.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			is.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, but fails on saving race", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		s.races.SaveFunc = func(context.Context, racers.Race) error {
			return errors.New("")
		}

		err := s.service.Create(context.Background(), req)
		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, and league exits, and saves race", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		s.races.SaveFunc = func(context.Context, racers.Race) error {
			return nil
		}

		err := s.service.Create(context.Background(), req)
		is.NoErr(err)

		is.Equal(len(s.races.SaveCalls()), 1)
		savedRace := s.races.SaveCalls()[0].Race

		events := savedRace.ConsumeEvents()
		is.Equal(len(events), 1)
		_, isCreated := events[0].(racers.RaceCreated)
		is.True(isCreated)

		is.Equal(racers.NewRace(raceID, raceName, raceDate), savedRace)
	})
}

func TestJoinRace(t *testing.T) {
	raceID := racers.RaceID(types.GenerateID())
	raceName := racers.RaceName("marathon")
	raceDate := racers.RaceDate(time.Now())

	userID := racers.UserID(types.GenerateID())
	userName := racers.UserName("Usain")

	req := service.JoinRace{
		RaceID: types.ID(raceID).String(),
		UserID: types.ID(userID).String(),
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		err := s.service.Join(context.Background(), service.JoinRace{})

		var ve errors.ValidationError
		is.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidRaceIDError{},
			racers.InvalidUserIDError{},
		}
		is.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			is.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, fails to get the race returns internal error", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return nil, errors.New("")
		}

		err := s.service.Join(context.Background(), req)
		pretty.Println(err)
		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, the race does not exists returns not found error", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return nil, nil
		}

		err := s.service.Join(context.Background(), req)

		is.True(errors.IsNotFoundError(err))
		var notFoundErr service.RaceByIDNotFoundError
		is.True(errors.As(err, &notFoundErr))
		is.Equal(notFoundErr.ID, raceID)
	})

	t.Run("When valid request, fails to get the user returns not found error", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		race := racers.NewRace(raceID, raceName, raceDate)
		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return &race, nil
		}

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, errors.New("")
		}

		err := s.service.Join(context.Background(), req)

		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, and the user does not exists returns not found error", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		race := racers.NewRace(raceID, raceName, raceDate)
		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return &race, nil
		}

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, nil
		}

		err := s.service.Join(context.Background(), req)

		is.True(errors.IsNotFoundError(err))
		var notFound service.UserByIDNotFoundError
		is.True(errors.As(err, &notFound))
		is.Equal(notFound.ID, userID)
	})

	t.Run("When user already joined race, returns CompetitorInRaceError error", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		race := racers.NewRace(raceID, raceName, raceDate, racers.RaceCompetitorsOpt(racers.NewRaceCompetitors(userID)))
		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return &race, nil
		}

		user := racers.NewUser(userID, userName)
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &user, nil
		}

		err := s.service.Join(context.Background(), req)

		var expectedErr racers.CompetitorInRaceError
		is.True(errors.As(err, &expectedErr))
	})

	t.Run("When user not in race but fails saving race, returns InternalError", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		race := racers.NewRace(raceID, raceName, raceDate)
		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return &race, nil
		}
		s.races.SaveFunc = func(context.Context, racers.Race) error {
			return errors.New("")
		}

		user := racers.NewUser(userID, userName)
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &user, nil
		}

		err := s.service.Join(context.Background(), req)
		is.True(errors.IsInternalError(err))
	})

	t.Run("When user not in race and saves the race, returns no error", func(t *testing.T) {
		is := is.New(t)
		s := newRacesServiceSuite()

		race := racers.NewRace(raceID, raceName, raceDate)
		s.races.ByIDFunc = func(context.Context, racers.RaceID) (*racers.Race, error) {
			return &race, nil
		}
		s.races.SaveFunc = func(context.Context, racers.Race) error {
			return nil
		}

		user := racers.NewUser(userID, userName)
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &user, nil
		}

		err := s.service.Join(context.Background(), req)
		is.NoErr(err)

		is.Equal(len(s.races.SaveCalls()), 1)
		savedRace := s.races.SaveCalls()[0].Race

		events := savedRace.ConsumeEvents()
		is.Equal(len(events), 1)
		_, isJoined := events[0].(racers.RaceCompetitorJoined)
		is.True(isJoined)
		is.Equal(racers.NewRace(raceID, raceName, raceDate, racers.RaceCompetitorsOpt(racers.NewRaceCompetitors(userID))), savedRace)
	})
}

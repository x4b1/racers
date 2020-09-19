package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
)

func TestRacesService(t *testing.T) {
	suite.Run(t, new(createRaceSuite))
	suite.Run(t, new(getRaceSuite))
	suite.Run(t, new(joinRaceSuite))
	suite.Run(t, new(listRacesSuite))
}

type createRaceSuite struct {
	suite.Suite

	service service.Races

	req service.CreateRace

	races    *RacesRepositoryMock
	users    *UsersGetterMock
	eventBus *EventBusMock
}

func (s *createRaceSuite) SetupTest() {
	s.races = &RacesRepositoryMock{}
	s.users = &UsersGetterMock{}
	s.eventBus = &EventBusMock{}

	s.req = service.CreateRace{
		ID:   id.Generate().String(),
		Name: "Black Mamba Race",
		Date: time.Now().AddDate(0, 1, 0),
	}

	s.service = service.NewRaces(s.races, s.users, service.NoopUnitOfWork, s.eventBus)
}

func (s createRaceSuite) TestCreateRace_InvalidRequest() {
	for field, c := range map[string]service.CreateRace{
		"id":   {Name: s.req.Name, Date: s.req.Date},
		"name": {ID: s.req.ID, Date: s.req.Date},
		"date": {ID: s.req.ID, Name: s.req.Name},
	} {
		s.Run(field, func() {
			_, err := s.service.Create(context.Background(), c)
			s.Error(err)
		})
	}
}

func (s createRaceSuite) TestCreateRace_CheckExistsFails() {
	s.races.ExistsFunc = func(context.Context, racers.Race) (bool, error) {
		return false, errors.New("")
	}

	_, err := s.service.Create(context.Background(), s.req)

	s.Error(err)
}

func (s createRaceSuite) TestCreateRace_AlreadyExists() {
	s.races.ExistsFunc = func(context.Context, racers.Race) (bool, error) {
		return true, nil
	}

	_, err := s.service.Create(context.Background(), s.req)

	s.Equal(service.ErrRaceAlreadyExists, err)
}

func (s createRaceSuite) TestCreateRace_SaveFails() {
	s.races.SaveFunc = func(context.Context, racers.Race) error {
		return errors.New("")
	}

	_, err := s.service.Create(context.Background(), s.req)
	s.Error(err)
}

func (s createRaceSuite) TestCreateRace_PublishEventsFails() {
	s.eventBus.PublishFunc = func(context.Context, ...service.Event) error {
		return errors.New("")
	}

	_, err := s.service.Create(context.Background(), s.req)
	s.Error(err)
}

func (s createRaceSuite) TestCreateRace_Success() {
	s.races.ExistsFunc = func(context.Context, racers.Race) (bool, error) {
		return false, nil
	}

	s.races.SaveFunc = func(context.Context, racers.Race) error {
		return nil
	}

	s.eventBus.PublishFunc = func(context.Context, ...service.Event) error {
		return nil
	}

	result, err := s.service.Create(context.Background(), s.req)
	s.NoError(err)

	expected := racers.NewRace(
		racers.RaceID(id.MustParse(s.req.ID)),
		racers.RaceName(s.req.Name),
		racers.RaceDate(s.req.Date),
		racers.UserID{},
	)
	s.Equal(expected, result)

	s.Len(s.races.ExistsCalls(), 1)
	s.Len(s.races.SaveCalls(), 1)

	s.Len(s.eventBus.PublishCalls(), 1)
	s.Len(s.eventBus.PublishCalls()[0].Events, 1)
	s.Equal(service.RaceCreated{Race: result}, s.eventBus.PublishCalls()[0].Events[0].Payload)
}

type getRaceSuite struct {
	suite.Suite

	service service.Races

	req service.GetRace

	dummyRace racers.Race

	races *RacesRepositoryMock
	users *UsersGetterMock
}

func (s *getRaceSuite) SetupTest() {
	s.races = &RacesRepositoryMock{}
	s.users = &UsersGetterMock{}

	s.req = service.GetRace{
		ID: id.Generate().String(),
	}

	s.service = service.NewRaces(s.races, s.users, service.NoopUnitOfWork, nil)
}

func (s getRaceSuite) TestGetRace_InvalidRequest() {
	_, err := s.service.Get(context.Background(), service.GetRace{ID: ""})
	s.Error(err)
}

func (s getRaceSuite) TestGetRace_Fails() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return racers.Race{}, errors.New("")
	}
	_, err := s.service.Get(context.Background(), s.req)
	s.Error(err)
}

func (s getRaceSuite) TestGetRace_Success() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return s.dummyRace, nil
	}

	result, err := s.service.Get(context.Background(), s.req)

	s.NoError(err)
	s.Equal(s.dummyRace, result)
}

type joinRaceSuite struct {
	suite.Suite

	service service.Races

	req service.JoinRace

	dummyRace racers.Race
	dummyUser racers.User

	races    *RacesRepositoryMock
	users    *UsersGetterMock
	eventBus *EventBusMock
}

func (s *joinRaceSuite) SetupTest() {
	s.races = &RacesRepositoryMock{}
	s.users = &UsersGetterMock{}
	s.eventBus = &EventBusMock{}

	s.dummyRace = racers.NewRace(
		racers.RaceID(id.Generate()),
		racers.RaceName("Black Mamba Race"),
		racers.RaceDate(time.Now().AddDate(0, 1, 0)),
		racers.UserID(id.Generate()),
	)

	s.dummyUser = racers.User{ID: racers.UserID(id.Generate())}

	s.req = service.JoinRace{
		RaceID: id.ID(s.dummyRace.ID).String(),
		UserID: id.ID(s.dummyUser.ID).String(),
	}

	s.service = service.NewRaces(s.races, s.users, service.NoopUnitOfWork, s.eventBus)
}

func (s joinRaceSuite) TestJoinRace_InvalidRequest() {
	for field, r := range map[string]service.JoinRace{
		"race_id": {UserID: s.req.UserID},
		"user_id": {RaceID: s.req.RaceID},
	} {
		s.Run(field, func() {
			err := s.service.Join(context.Background(), r)
			s.Error(err)
		})
	}
}

func (s joinRaceSuite) TestJoinRace_FailsGettingRace() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return racers.Race{}, errors.New("")
	}

	err := s.service.Join(context.Background(), s.req)
	s.Error(err)
}

func (s joinRaceSuite) TestJoinRace_FailsGettingUser() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return s.dummyRace, nil
	}

	s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
		return racers.User{}, errors.New("")
	}

	err := s.service.Join(context.Background(), s.req)
	s.Error(err)
}

func (s joinRaceSuite) TestJoinRace_FailsJoiningRace() {
	s.dummyRace.Join(s.dummyUser)

	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return s.dummyRace, nil
	}

	s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
		return s.dummyUser, nil
	}

	err := s.service.Join(context.Background(), s.req)
	s.Error(err)
}

func (s joinRaceSuite) TestJoinRace_FailsSaving() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return s.dummyRace, nil
	}
	s.races.SaveFunc = func(context.Context, racers.Race) error {
		return errors.New("")
	}

	s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
		return s.dummyUser, nil
	}

	err := s.service.Join(context.Background(), s.req)
	s.Error(err)
}

func (s joinRaceSuite) TestJoinRace_PublishEventsFails() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return s.dummyRace, nil
	}
	s.races.SaveFunc = func(context.Context, racers.Race) error {
		return nil
	}

	s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
		return s.dummyUser, nil
	}

	s.eventBus.PublishFunc = func(context.Context, ...service.Event) error {
		return errors.New("")
	}

	s.Error(s.service.Join(context.Background(), s.req))
}

func (s joinRaceSuite) TestJoinRace_Success() {
	s.races.GetFunc = func(context.Context, racers.RaceID) (racers.Race, error) {
		return s.dummyRace, nil
	}
	s.races.SaveFunc = func(context.Context, racers.Race) error {
		return nil
	}

	s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
		return s.dummyUser, nil
	}

	s.eventBus.PublishFunc = func(context.Context, ...service.Event) error {
		return nil
	}

	s.NoError(s.service.Join(context.Background(), s.req))

	s.Len(s.races.SaveCalls(), 1)

	s.Len(s.eventBus.PublishCalls(), 1)
	s.Len(s.eventBus.PublishCalls()[0].Events, 1)

	s.dummyRace.Competitors = racers.NewRaceCompetitors(s.dummyUser.ID)

	s.Equal(
		service.UserJoinedRace{
			User: s.dummyUser,
			Race: s.dummyRace,
		},
		s.eventBus.PublishCalls()[0].Events[0].Payload,
	)
}

type listRacesSuite struct {
	suite.Suite

	service service.Races

	races *RacesRepositoryMock
}

func (s *listRacesSuite) SetupTest() {
	s.races = &RacesRepositoryMock{}

	s.service = service.NewRaces(s.races, nil, service.NoopUnitOfWork, nil)
}

func (s listRacesSuite) TestListRaces_Success() {
	owner := racers.UserID(id.Generate())
	races := make([]racers.Race, 5)
	for i := range races {
		races[i] = racers.NewRace(racers.RaceID(id.Generate()), racers.RaceName("race-name"), racers.RaceDate{}, owner)
	}

	s.races.AllFunc = func(ctx context.Context) ([]racers.Race, error) {
		return races, nil
	}

	result, err := s.service.List(context.Background())

	s.Equal(races, result)
	s.NoError(err)
}

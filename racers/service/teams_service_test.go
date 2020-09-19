package service_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/xabi93/go-clean/internal/errors"
	"github.com/xabi93/go-clean/internal/types"
	"github.com/xabi93/go-clean/racers"
	"github.com/xabi93/go-clean/racers/service"
)

var (
	teamID      = racers.TeamID(types.GenerateID())
	teamName    = racers.TeamName("black panters")
	teamAdminID = racers.UserID(types.GenerateID())
	teamAdmin   = racers.NewUser(teamAdminID, racers.UserName("Usain"))
)

type testTeamsService struct {
	service service.TeamsService
	teams   *TeamsRepositoryMock
	users   *UsersRepositoryMock
}

func newTestTeamsService() testTeamsService {
	s := testTeamsService{
		teams: &TeamsRepositoryMock{},
		users: &UsersRepositoryMock{},
	}

	s.service = service.NewTeamsService(s.teams, s.users)

	return s
}

func TestCreateTeam(t *testing.T) {
	req := service.CreateTeam{
		ID:      types.ID(teamID).String(),
		Name:    string(teamName),
		AdminID: types.ID(teamAdminID).String(),
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		err := s.service.Create(context.Background(), service.CreateTeam{})

		var ve errors.ValidationError
		is.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidTeamIDError{},
			racers.InvalidTeamNameError{},
			racers.InvalidUserIDError{},
		}
		is.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			is.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, but fails on getting the admin", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, errors.New("")
		}

		err := s.service.Create(context.Background(), req)

		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, but admin does not exists", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, nil
		}

		err := s.service.Create(context.Background(), req)

		is.True(errors.IsNotFoundError(err))
		var notFound service.UserByIDNotFoundError
		is.True(errors.As(err, &notFound))
		is.Equal(notFound.ID, teamAdminID)
	})

	t.Run("When valid request, but fails on saving team", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &teamAdmin, nil
		}

		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return errors.New("")
		}

		err := s.service.Create(context.Background(), req)

		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, and saves team", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &teamAdmin, nil
		}

		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return nil
		}

		err := s.service.Create(context.Background(), req)
		is.NoErr(err)

		is.Equal(len(s.teams.SaveCalls()), 1)
		savedTeam := s.teams.SaveCalls()[0].Team

		events := savedTeam.ConsumeEvents()
		is.Equal(len(events), 1)
		_, isCreated := events[0].(racers.TeamCreated)
		is.True(isCreated)

		is.Equal(racers.NewTeam(teamID, teamName, teamAdminID, racers.TeamMembersOpt(racers.NewTeamMembers(teamAdminID))), savedTeam)
	})
}

func TestJoinTeam(t *testing.T) {
	newMemberID := racers.UserID(types.GenerateID())
	newMember := racers.NewUser(newMemberID, racers.UserName("Martin"))

	req := service.JoinTeam{
		TeamID: types.ID(teamID).String(),
		UserID: string(newMemberID),
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		err := s.service.Join(context.Background(), service.JoinTeam{})

		var ve errors.ValidationError
		is.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidTeamIDError{},
			racers.InvalidUserIDError{},
		}
		is.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			is.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, but fails on getting the team", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return nil, errors.New("")
		}

		err := s.service.Join(context.Background(), req)

		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, but team does not exists", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return nil, nil
		}

		err := s.service.Join(context.Background(), req)

		is.True(errors.IsNotFoundError(err))
		var notFound service.TeamByIDNotFoundError
		is.True(errors.As(err, &notFound))
		is.Equal(notFound.ID, teamID)
	})

	t.Run("When valid request, but fails on getting the user", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &racers.Team{}, nil
		}

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, errors.New("")
		}

		is.True(errors.IsInternalError(s.service.Join(context.Background(), req)))
	})

	t.Run("When valid request, but user does not exists", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &racers.Team{}, nil
		}
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, nil
		}

		err := s.service.Join(context.Background(), req)

		is.True(errors.IsNotFoundError(err))
		var notFound service.UserByIDNotFoundError
		is.True(errors.As(err, &notFound))
		is.Equal(notFound.ID, newMemberID)
	})

	t.Run("When valid request, but fails when tries to get the user team", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &racers.Team{}, nil
		}
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return nil, errors.New("")
		}

		is.True(errors.IsInternalError(s.service.Join(context.Background(), req)))
	})

	t.Run("When valid request, and the user is already in other team", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		team := racers.NewTeam(teamID, teamName, teamAdminID)
		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &team, nil
		}
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return &team, nil
		}

		err := s.service.Join(context.Background(), req)

		var expectedErr racers.UserAlreadyInTeam
		is.True(errors.As(err, &expectedErr))
		is.Equal(expectedErr.TeamID, teamID)
		is.Equal(expectedErr.UserID, newMemberID)
	})

	t.Run("When valid request, and fails saving", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		team := racers.NewTeam(teamID, teamName, teamAdminID)
		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &team, nil
		}
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return nil, nil
		}
		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return errors.New("")
		}

		is.True(errors.IsInternalError(s.service.Join(context.Background(), req)))
	})

	t.Run("When valid request, and saves team", func(t *testing.T) {
		is := is.New(t)
		s := newTestTeamsService()

		team := racers.NewTeam(teamID, teamName, teamAdminID)
		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &team, nil
		}
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return nil, nil
		}
		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return nil
		}

		err := s.service.Join(context.Background(), req)
		is.NoErr(err)

		is.Equal(len(s.teams.ByIDCalls()), 1)
		is.Equal(s.teams.ByIDCalls()[0].ID, teamID)

		is.Equal(len(s.users.ByIDCalls()), 1)
		is.Equal(s.users.ByIDCalls()[0].ID, newMemberID)

		is.Equal(len(s.teams.ByMemberCalls()), 1)
		is.Equal(s.teams.ByMemberCalls()[0].MemberID, newMemberID)

		is.Equal(len(s.teams.SaveCalls()), 1)
		savedTeam := s.teams.SaveCalls()[0].Team

		events := savedTeam.ConsumeEvents()
		is.Equal(len(events), 1)
		_, isJoined := events[0].(racers.UserJoinedTeam)
		is.True(isJoined)
	})
}

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
)

var (
	teamID      = racers.TeamID(id.GenerateID())
	teamName    = racers.TeamName("black panters")
	teamAdminID = racers.UserID(id.GenerateID())
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
	require := require.New(t)

	req := service.CreateTeam{
		ID:      id.ID(teamID).String(),
		Name:    string(teamName),
		AdminID: id.ID(teamAdminID).String(),
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		s := newTestTeamsService()

		err := s.service.Create(context.Background(), service.CreateTeam{})

		var ve errors.ValidationError
		require.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidTeamIDError{},
			racers.InvalidTeamNameError{},
			racers.InvalidUserIDError{},
		}
		require.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			require.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, but fails on getting the admin", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, errors.New("")
		}

		err := s.service.Create(context.Background(), req)

		require.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, but admin does not exists", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, nil
		}

		err := s.service.Create(context.Background(), req)

		require.True(errors.IsNotFoundError(err))
		var notFound service.UserByIDNotFoundError
		require.True(errors.As(err, &notFound))
		require.Equal(notFound.ID, teamAdminID)
	})

	t.Run("When valid request, but fails on saving team", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &teamAdmin, nil
		}

		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return errors.New("")
		}

		err := s.service.Create(context.Background(), req)

		require.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, and saves team", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return &teamAdmin, nil
		}

		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return nil
		}

		err := s.service.Create(context.Background(), req)
		require.NoError(err)

		require.Equal(len(s.teams.SaveCalls()), 1)
		savedTeam := s.teams.SaveCalls()[0].Team

		events := savedTeam.ConsumeEvents()
		require.Equal(len(events), 1)
		_, isCreated := events[0].(racers.TeamCreated)
		require.True(isCreated)

		require.Equal(racers.NewTeam(teamID, teamName, teamAdminID, racers.TeamMembersOpt(racers.NewTeamMembers(teamAdminID))), savedTeam)
	})
}

func TestJoinTeam(t *testing.T) {
	require := require.New(t)

	newMemberID := racers.UserID(id.GenerateID())
	newMember := racers.NewUser(newMemberID, racers.UserName("Martin"))

	req := service.JoinTeam{
		TeamID: id.ID(teamID).String(),
		UserID: string(newMemberID),
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		s := newTestTeamsService()

		err := s.service.Join(context.Background(), service.JoinTeam{})

		var ve errors.ValidationError
		require.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidTeamIDError{},
			racers.InvalidUserIDError{},
		}
		require.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			require.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, but fails on getting the team", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return nil, errors.New("")
		}

		err := s.service.Join(context.Background(), req)

		require.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, but team does not exists", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return nil, nil
		}

		err := s.service.Join(context.Background(), req)

		require.True(errors.IsNotFoundError(err))
		var notFound service.TeamByIDNotFoundError
		require.True(errors.As(err, &notFound))
		require.Equal(notFound.ID, teamID)
	})

	t.Run("When valid request, but fails on getting the user", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &racers.Team{}, nil
		}

		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, errors.New("")
		}

		require.True(errors.IsInternalError(s.service.Join(context.Background(), req)))
	})

	t.Run("When valid request, but user does not exists", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.ByIDFunc = func(context.Context, racers.TeamID) (*racers.Team, error) {
			return &racers.Team{}, nil
		}
		s.users.ByIDFunc = func(context.Context, racers.UserID) (*racers.User, error) {
			return nil, nil
		}

		err := s.service.Join(context.Background(), req)

		require.True(errors.IsNotFoundError(err))
		var notFound service.UserByIDNotFoundError
		require.True(errors.As(err, &notFound))
		require.Equal(notFound.ID, newMemberID)
	})

	t.Run("When valid request, but fails when tries to get the user team", func(t *testing.T) {
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

		require.True(errors.IsInternalError(s.service.Join(context.Background(), req)))
	})

	t.Run("When valid request, and the user is already in other team", func(t *testing.T) {
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
		require.True(errors.As(err, &expectedErr))
		require.Equal(expectedErr.TeamID, teamID)
		require.Equal(expectedErr.UserID, newMemberID)
	})

	t.Run("When valid request, and fails saving", func(t *testing.T) {
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

		require.True(errors.IsInternalError(s.service.Join(context.Background(), req)))
	})

	t.Run("When valid request, and saves team", func(t *testing.T) {
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
		require.NoError(err)

		require.Equal(len(s.teams.ByIDCalls()), 1)
		require.Equal(s.teams.ByIDCalls()[0].ID, teamID)

		require.Equal(len(s.users.ByIDCalls()), 1)
		require.Equal(s.users.ByIDCalls()[0].ID, newMemberID)

		require.Equal(len(s.teams.ByMemberCalls()), 1)
		require.Equal(s.teams.ByMemberCalls()[0].MemberID, newMemberID)

		require.Equal(len(s.teams.SaveCalls()), 1)
		savedTeam := s.teams.SaveCalls()[0].Team

		events := savedTeam.ConsumeEvents()
		require.Equal(len(events), 1)
		_, isJoined := events[0].(racers.UserJoinedTeam)
		require.True(isJoined)
	})
}
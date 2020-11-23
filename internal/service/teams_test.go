package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
)

var (
	teamID      = racers.TeamID(id.Generate())
	teamName    = racers.TeamName("black panters")
	teamAdminID = racers.UserID(id.Generate())
	teamAdmin   = racers.User{ID: teamAdminID}
)

type testTeamsService struct {
	service service.Teams
	teams   *TeamsRepositoryMock
	users   *UsersGetterMock
}

func newTestTeamsService() testTeamsService {
	s := testTeamsService{
		teams: &TeamsRepositoryMock{},
		users: &UsersGetterMock{},
	}

	s.service = service.NewTeams(s.teams, s.users)

	return s
}

func TestCreateTeam(t *testing.T) {
	require := require.New(t)

	req := service.CreateTeam{
		ID:      id.ID(teamID).String(),
		Name:    string(teamName),
		AdminID: id.ID(teamAdminID).String(),
	}

	t.Run("Scenario: invalid request", func(t *testing.T) {
		type testCase struct {
			req service.CreateTeam
		}
		for field, c := range map[string]testCase{
			"id":   {req: service.CreateTeam{Name: req.Name, AdminID: req.AdminID}},
			"name": {req: service.CreateTeam{ID: req.ID, AdminID: req.AdminID}},
			"date": {req: service.CreateTeam{ID: req.ID, Name: req.Name}},
		} {
			t.Run(fmt.Sprintf("when invalid %s", field), func(t *testing.T) {
				s := newTestTeamsService()
				err := s.service.Create(context.Background(), c.req)
				require.Error(err)
			})
		}
	})

	t.Run("When valid request, but fails on getting the admin", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return racers.User{}, service.ErrUserNotFound
		}

		require.Error(s.service.Create(context.Background(), req))
	})

	t.Run("When valid request, but fails on saving team", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return teamAdmin, nil
		}

		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return errors.New("")
		}

		require.Error(s.service.Create(context.Background(), req))
	})

	t.Run("When valid request, and saves team", func(t *testing.T) {
		s := newTestTeamsService()

		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return teamAdmin, nil
		}

		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return nil
		}

		err := s.service.Create(context.Background(), req)
		require.NoError(err)

		require.Equal(len(s.teams.SaveCalls()), 1)
		savedTeam := s.teams.SaveCalls()[0].Team

		require.Equal(racers.NewTeam(teamID, teamName, teamAdminID, racers.TeamMembersOpt(racers.NewTeamMembers(teamAdminID))), savedTeam)
	})
}

func TestJoinTeam(t *testing.T) {
	require := require.New(t)

	team := racers.NewTeam(teamID, teamName, teamAdminID)

	newMemberID := racers.UserID(id.Generate())
	newMember := racers.User{ID: newMemberID}

	req := service.JoinTeam{
		TeamID: id.ID(teamID).String(),
		UserID: id.ID(newMemberID).String(),
	}

	t.Run("Scenario: invalid request", func(t *testing.T) {
		type testCase struct {
			req service.JoinTeam
		}
		for field, c := range map[string]testCase{
			"team id": {req: service.JoinTeam{UserID: req.UserID}},
			"user id": {req: service.JoinTeam{TeamID: req.TeamID}},
		} {
			t.Run(fmt.Sprintf("when invalid %s", field), func(t *testing.T) {
				s := newTestTeamsService()
				err := s.service.Join(context.Background(), c.req)
				require.Error(err)
			})
		}
	})

	t.Run("When valid request, but fails on getting the team", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.GetFunc = func(context.Context, racers.TeamID) (racers.Team, error) {
			return racers.Team{}, service.ErrTeamNotFound
		}

		require.Error(s.service.Join(context.Background(), req))
	})

	t.Run("When valid request, but fails on getting the user", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.GetFunc = func(context.Context, racers.TeamID) (racers.Team, error) {
			return racers.Team{}, nil
		}

		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return racers.User{}, service.ErrUserNotFound
		}

		require.Error(s.service.Join(context.Background(), req))
	})

	t.Run("When valid request, but fails when tries to get the user team", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.GetFunc = func(context.Context, racers.TeamID) (racers.Team, error) {
			return team, nil
		}
		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return nil, errors.New("")
		}

		require.Error(s.service.Join(context.Background(), req))
	})

	t.Run("When valid request, and the user is already in other team", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.GetFunc = func(context.Context, racers.TeamID) (racers.Team, error) {
			return team, nil
		}
		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return &team, nil
		}

		err := s.service.Join(context.Background(), req)

		var expectedErr racers.UserAlreadyInTeamError
		require.True(errors.As(err, &expectedErr))
		require.Equal(expectedErr.TeamID, teamID)
		require.Equal(expectedErr.UserID, newMemberID)
	})

	t.Run("When valid request, and fails saving", func(t *testing.T) {
		s := newTestTeamsService()

		s.teams.GetFunc = func(context.Context, racers.TeamID) (racers.Team, error) {
			return team, nil
		}
		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return nil, nil
		}
		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return errors.New("")
		}

		require.Error(s.service.Join(context.Background(), req))
	})

	t.Run("When valid request, and saves team", func(t *testing.T) {
		s := newTestTeamsService()

		team := racers.NewTeam(teamID, teamName, teamAdminID)
		s.teams.GetFunc = func(context.Context, racers.TeamID) (racers.Team, error) {
			return team, nil
		}
		s.users.GetFunc = func(context.Context, racers.UserID) (racers.User, error) {
			return newMember, nil
		}
		s.teams.ByMemberFunc = func(context.Context, racers.UserID) (*racers.Team, error) {
			return nil, nil
		}
		s.teams.SaveFunc = func(context.Context, racers.Team) error {
			return nil
		}

		err := s.service.Join(context.Background(), req)
		require.NoError(err)

		require.Equal(len(s.teams.GetCalls()), 1)
		require.Equal(s.teams.GetCalls()[0].ID, teamID)

		require.Equal(len(s.users.GetCalls()), 1)
		require.Equal(s.users.GetCalls()[0].ID, newMemberID)

		require.Equal(len(s.teams.ByMemberCalls()), 1)
		require.Equal(s.teams.ByMemberCalls()[0].ID, newMemberID)

		require.Equal(len(s.teams.SaveCalls()), 1)
	})
}

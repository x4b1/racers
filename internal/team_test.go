package racers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
)

var (
	teamID      = racers.TeamID(id.Generate())
	teamName    = racers.TeamName("athletic")
	teamAdminID = racers.UserID(id.Generate())
)

func TestTeamID(t *testing.T) {
	require := require.New(t)
	t.Run("when invalid id returns InvalidTeamIDError error", func(t *testing.T) {
		_, err := racers.NewTeamID("")
		require.True(errors.As(err, &racers.InvalidTeamIDError{}))
	})

	t.Run("when valid id returns TeamID and no error", func(t *testing.T) {
		id := id.Generate()
		teamID, err := racers.NewTeamID(id.String())

		require.Equal(racers.TeamID(id), teamID)
		require.NoError(err)
	})
}

func TestTeamName(t *testing.T) {
	require := require.New(t)
	t.Run("when New with empty name returns InvalidTeamNameError error", func(t *testing.T) {
		_, err := racers.NewTeamName("")
		require.True(errors.As(err, &racers.InvalidTeamNameError{}))
	})

	t.Run("when New with valid name returns TeamName and no error", func(t *testing.T) {
		name := "athletic"
		teamName, err := racers.NewTeamName(name)

		require.Equal(racers.TeamName(name), teamName)
		require.NoError(err)
	})
}

func TestNewTeam(t *testing.T) {
	require := require.New(t)

	team := racers.NewTeam(teamID, teamName, teamAdminID)

	require.Equal(team.ID, teamID)
	require.Equal(team.Name, teamName)
	require.Equal(team.Admin, teamAdminID)
}

func TestCreateTeam(t *testing.T) {
	require := require.New(t)

	r := racers.CreateTeam(teamID, teamName, racers.User{teamAdminID})

	require.Equal(r, racers.NewTeam(teamID, teamName, teamAdminID, racers.TeamMembersOpt(racers.NewTeamMembers(teamAdminID))))
}

func TestJoinTeam(t *testing.T) {
	require := require.New(t)
	t.Run(`Given a team with only with admin,
	When joins other user with team,
	Then returns UserAlreadyInTeam error`, func(t *testing.T) {
		team := racers.NewTeam(teamID, teamName, teamAdminID)

		_, err := racers.JoinTeam(team, racers.User{userID}, &team)

		require.True(errors.Is(err, racers.UserAlreadyInTeamError{UserID: userID, TeamID: teamID}))
	})

	t.Run(`Given a team with only with admin,
	When joins other user without team,
	Then returns no error and generates event`, func(t *testing.T) {
		team := racers.NewTeam(teamID, teamName, teamAdminID)

		team, err := racers.JoinTeam(team, racers.User{userID}, nil)
		require.NoError(err)
	})
}

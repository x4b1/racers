package racers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
)

var (
	teamID      = racers.TeamID(id.GenerateID())
	teamName    = racers.TeamName("athletic")
	teamAdminID = racers.UserID(id.GenerateID())
)

func TestTeamID(t *testing.T) {
	require := require.New(t)
	t.Run("when invalid id returns InvalidTeamIDError error", func(t *testing.T) {
		_, err := racers.NewTeamID("")
		require.True(errors.As(err, &racers.InvalidTeamIDError{}))
	})

	t.Run("when valid id returns TeamID and no error", func(t *testing.T) {
		id := id.GenerateID()
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

	require.Equal(team.ID(), teamID)
	require.Equal(team.Name(), teamName)
	require.Equal(team.Admin(), teamAdminID)

	events := team.ConsumeEvents()
	require.Empty(events)
}

func TestCreateTeam(t *testing.T) {
	require := require.New(t)

	r := racers.CreateTeam(teamID, teamName, racers.NewUser(teamAdminID, ""))

	events := r.ConsumeEvents()
	require.Len(events, 1)
	createdEvent := events[0].(racers.TeamCreated)
	require.Equal(createdEvent.TeamID, teamID)
	require.Equal(createdEvent.TeamName, teamName)
	require.Equal(createdEvent.TeamAdmin, teamAdminID)
	require.Equal(createdEvent.TeamMembers, racers.NewTeamMembers(teamAdminID))

	require.Equal(r, racers.NewTeam(teamID, teamName, teamAdminID, racers.TeamMembersOpt(racers.NewTeamMembers(teamAdminID))))
}

func TestTeam(t *testing.T) {
	require := require.New(t)
	t.Run(`Given a team with only with admin,
	When joins other user with team,
	Then returns UserAlreadyInTeam error`, func(t *testing.T) {
		team := racers.NewTeam(teamID, teamName, teamAdminID)

		err := team.Join(racers.NewTeamMember(racers.NewUser(teamAdminID, ""), &team))

		require.True(errors.Is(err, racers.UserAlreadyInTeam{UserID: teamAdminID, TeamID: teamID}))
		require.Equal(len(team.ConsumeEvents()), 0)
	})

	t.Run(`Given a team with only with admin,
	When joins other user without team,
	Then returns no error and generates event`, func(t *testing.T) {
		team := racers.NewTeam(teamID, teamName, teamAdminID)

		newUserID := racers.UserID(id.GenerateID())
		require.NoError(team.Join(racers.NewTeamMember(racers.NewUser(newUserID, ""), nil)))

		events := team.ConsumeEvents()
		require.Equal(len(events), 1)
		joinedEvent := events[0].(racers.UserJoinedTeam)
		require.Equal(joinedEvent.TeamID, teamID)
		require.Equal(joinedEvent.UserID, newUserID)
	})
}

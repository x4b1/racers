package racers_test

import (
	"errors"
	"testing"

	"github.com/matryer/is"
	"github.com/xabi93/go-clean/internal/types"
	"github.com/xabi93/go-clean/racers"
)

var (
	teamID      = racers.TeamID(types.GenerateID())
	teamName    = racers.TeamName("athletic")
	teamAdminID = racers.UserID(types.GenerateID())
)

func TestTeamID(t *testing.T) {
	is := is.New(t)
	t.Run("when invalid id returns InvalidTeamIDError error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewTeamID("")
		is.True(errors.As(err, &racers.InvalidTeamIDError{}))
	})

	t.Run("when valid id returns TeamID and no error", func(t *testing.T) {
		is := is.New(t)

		id := types.GenerateID()
		teamID, err := racers.NewTeamID(id.String())

		is.Equal(racers.TeamID(id), teamID)
		is.NoErr(err)
	})
}

func TestTeamName(t *testing.T) {
	is := is.New(t)
	t.Run("when New with empty name returns InvalidTeamNameError error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewTeamName("")
		is.True(errors.As(err, &racers.InvalidTeamNameError{}))
	})

	t.Run("when New with valid name returns TeamName and no error", func(t *testing.T) {
		is := is.New(t)

		name := "athletic"
		teamName, err := racers.NewTeamName(name)

		is.Equal(racers.TeamName(name), teamName)
		is.NoErr(err)
	})
}

func TestNewTeam(t *testing.T) {
	is := is.New(t)

	team := racers.NewTeam(teamID, teamName, teamAdminID)

	is.Equal(team.ID(), teamID)
	is.Equal(team.Name(), teamName)
	is.Equal(team.Admin(), teamAdminID)

	events := team.ConsumeEvents()
	is.True(len(events) == 0)
}

func TestCreateTeam(t *testing.T) {
	is := is.New(t)

	r := racers.CreateTeam(teamID, teamName, racers.NewUser(teamAdminID, ""))

	events := r.ConsumeEvents()
	is.True(len(events) == 1)
	createdEvent := events[0].(racers.TeamCreated)
	is.Equal(createdEvent.TeamID, teamID)
	is.Equal(createdEvent.TeamName, teamName)
	is.Equal(createdEvent.TeamAdmin, teamAdminID)
	is.Equal(createdEvent.TeamMembers, racers.NewTeamMembers(teamAdminID))

	is.Equal(r, racers.NewTeam(teamID, teamName, teamAdminID, racers.TeamMembersOpt(racers.NewTeamMembers(teamAdminID))))
}

func TestTeam(t *testing.T) {
	t.Run(`Given a team with only with admin,
	When joins other user with team,
	Then returns UserAlreadyInTeam error`, func(t *testing.T) {
		is := is.New(t)

		team := racers.NewTeam(teamID, teamName, teamAdminID)

		err := team.Join(racers.NewTeamMember(racers.NewUser(teamAdminID, ""), &team))

		is.True(errors.Is(err, racers.UserAlreadyInTeam{UserID: teamAdminID, TeamID: teamID}))
		is.Equal(len(team.ConsumeEvents()), 0)
	})

	t.Run(`Given a team with only with admin,
	When joins other user without team,
	Then returns no error and generates event`, func(t *testing.T) {
		is := is.New(t)

		team := racers.NewTeam(teamID, teamName, teamAdminID)

		newUserID := racers.UserID(types.GenerateID())
		is.NoErr(team.Join(racers.NewTeamMember(racers.NewUser(newUserID, ""), nil)))

		events := team.ConsumeEvents()
		is.Equal(len(events), 1)
		joinedEvent := events[0].(racers.UserJoinedTeam)
		is.Equal(joinedEvent.TeamID, teamID)
		is.Equal(joinedEvent.UserID, newUserID)
	})
}

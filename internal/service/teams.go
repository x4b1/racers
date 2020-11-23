package service

import (
	"context"

	racers "github.com/xabi93/racers/internal"
)

func NewTeams(teams TeamsRepository, users UsersGetter) Teams {
	return Teams{teams, users}
}

type Teams struct {
	teams TeamsRepository
	users UsersGetter
}

type CreateTeam struct {
	ID      string
	Name    string
	AdminID string
}

func (t Teams) Create(ctx context.Context, r CreateTeam) error {
	id, err := racers.NewTeamID(r.ID)
	if err != nil {
		return err
	}

	name, err := racers.NewTeamName(r.Name)
	if err != nil {
		return err
	}

	adminID, err := racers.NewUserID(r.AdminID)
	if err != nil {
		return err
	}

	admin, err := t.users.Get(ctx, adminID)
	if err != nil {
		return err
	}

	return t.teams.Save(ctx, racers.CreateTeam(id, name, admin))
}

type JoinTeam struct {
	TeamID string
	UserID string
}

func (t Teams) Join(ctx context.Context, r JoinTeam) error {
	teamID, err := racers.NewTeamID(r.TeamID)
	if err != nil {
		return err
	}
	userID, err := racers.NewUserID(r.UserID)
	if err != nil {
		return err
	}

	team, err := t.teams.Get(ctx, teamID)
	if err != nil {
		return err
	}

	user, err := t.users.Get(ctx, userID)
	if err != nil {
		return err
	}

	userTeam, err := t.teams.ByMember(ctx, userID)
	if err != nil {
		return err
	}

	team, err = racers.JoinTeam(team, user, userTeam)
	if err != nil {
		return err
	}

	if err := t.teams.Save(ctx, team); err != nil {
		return err
	}

	return nil
}

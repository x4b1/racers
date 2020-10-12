package service

import (
	"context"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
)

type TeamsService interface {
	Create(ctx context.Context, r CreateTeam) error
	Join(ctx context.Context, r JoinTeam) error
}

func NewTeamsService(teams TeamsRepository, users UsersGetter) TeamsService {
	return teamsService{teams, users}
}

type teamsService struct {
	teams TeamsRepository
	users UsersGetter
}

type CreateTeam struct {
	ID      string
	Name    string
	AdminID string
}

func (rs teamsService) Create(ctx context.Context, r CreateTeam) error {
	var ve errors.ValidationError

	id, err := racers.NewTeamID(r.ID)
	if err != nil {
		ve.Add(err)
	}
	name, err := racers.NewTeamName(r.Name)
	if err != nil {
		ve.Add(err)
	}
	adminID, err := racers.NewUserID(r.AdminID)
	if err != nil {
		ve.Add(err)
	}
	if err := ve.Valid(); err != nil {
		return err
	}

	admin, err := rs.users.ByID(ctx, adminID)
	if err != nil {
		return errors.WrapInternalError(err, "fetching admin")
	}
	if admin == nil {
		return errors.WrapNotFoundError(UserByIDNotFoundError{ID: adminID}, "fetching admin")
	}

	if err := rs.teams.Save(ctx, racers.CreateTeam(id, name, *admin)); err != nil {
		return errors.WrapInternalError(err, "saving team")
	}

	return nil
}

type JoinTeam struct {
	TeamID string
	UserID string
}

func (rs teamsService) Join(ctx context.Context, r JoinTeam) error {
	var ve errors.ValidationError

	teamID, err := racers.NewTeamID(r.TeamID)
	if err != nil {
		ve.Add(err)
	}
	userID, err := racers.NewUserID(r.UserID)
	if err != nil {
		ve.Add(err)
	}
	if err := ve.Valid(); err != nil {
		return err
	}

	team, err := rs.teams.ByID(ctx, teamID)
	if err != nil {
		return errors.WrapInternalError(err, "fetching team")
	}
	if team == nil {
		return errors.WrapNotFoundError(TeamByIDNotFoundError{ID: teamID}, "fetching team")
	}

	user, err := rs.users.ByID(ctx, userID)
	if err != nil {
		return errors.WrapInternalError(err, "fetching user")
	}
	if user == nil {
		return errors.WrapNotFoundError(UserByIDNotFoundError{ID: userID}, "fetching user")
	}

	userTeam, err := rs.teams.ByMember(ctx, userID)
	if err != nil {
		return errors.WrapInternalError(err, "fetching user team")
	}

	if err := team.Join(racers.NewTeamMember(*user, userTeam)); err != nil {
		return err
	}

	if err := rs.teams.Save(ctx, *team); err != nil {
		return errors.WrapInternalError(err, "saving team")
	}

	return nil
}

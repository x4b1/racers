package racers

import (
	"errors"
	"fmt"

	"github.com/xabi93/racers/internal/id"
)

type (
	// TeamID defines a unique identifier for a team
	TeamID id.ID
	// InvalidTeamIDError means the given id is not valid
	InvalidTeamIDError struct{ error }
)

func (err InvalidTeamIDError) Error() string {
	return fmt.Sprintf("invalid team id: %s", err.error)
}

// NewTeamID validates the id and returns a TeamID instance
func NewTeamID(s string) (TeamID, error) {
	id, err := id.NewID(s)
	if err != nil {
		return TeamID{}, InvalidTeamIDError{err}
	}

	return TeamID(id), nil
}

type (
	// TeamName defines the name of a team
	TeamName             string
	InvalidTeamNameError struct{ error }
)

func (err InvalidTeamNameError) Error() string {
	return fmt.Sprintf("invalid team name: %s", err.error)
}

// NewTeamName validates the name and returns a TeamName instance
func NewTeamName(s string) (TeamName, error) {
	if s == "" {
		return "", InvalidTeamNameError{errors.New("empty name")}
	}

	return TeamName(s), nil
}

// TeamMembers is a list of users that are in a team
type TeamMembers struct{ userList }

type teamOption func(*Team)

// TeamMembersOpt is an optional parameter for NewTeam to initialize a team with members
func TeamMembersOpt(c TeamMembers) teamOption {
	return func(r *Team) {
		r.Members = c
	}
}

// NewTeam is a constructor for Team
func NewTeam(id TeamID, name TeamName, admin UserID, opts ...teamOption) Team {
	r := Team{
		ID:    id,
		Name:  name,
		Admin: admin,
	}
	for _, opt := range opts {
		opt(&r)
	}

	return r
}

// CreateTeam creates a new team and attach TeamCreated event
func CreateTeam(teamID TeamID, name TeamName, admin User) Team {
	return NewTeam(teamID, name, admin.ID, TeamMembersOpt(NewTeamMembers(admin.ID)))
}

// NewTeamMembers is a constructor that given a users ids it builds TeamMembers list
func NewTeamMembers(users ...UserID) TeamMembers {
	ul := make(userList, len(users))
	for _, u := range users {
		ul.add(u)
	}

	return TeamMembers{ul}
}

// ErrUnknownTeam means the team does not exists in the system
var ErrUnknownTeam = errors.New("unknown team")

// Team represents a team in the service
type Team struct {
	ID      TeamID      `json:"id,omitempty"`
	Name    TeamName    `json:"name,omitempty"`
	Admin   UserID      `json:"admin,omitempty"`
	Members TeamMembers `json:"members,omitempty"`
}

// UserAlreadyInTeamError a user cannot join a team because it's already in a team
type UserAlreadyInTeamError struct {
	UserID UserID
	TeamID TeamID
}

func (err UserAlreadyInTeamError) Error() string {
	return fmt.Sprintf("user %s cant join is member of %s team", err.UserID, err.TeamID)
}

// Join add a user as a team member
func (t *Team) join(u User) {
	t.Members.add(u.ID)
}

// JoinTeam is a domain service that checks if the given user is in a team and in that case returns UserAlreadyInTeamError
// if not adds the user to the team
func JoinTeam(joinTeam Team, u User, t *Team) (Team, error) {
	if t != nil {
		return Team{}, UserAlreadyInTeamError{UserID: u.ID, TeamID: t.ID}
	}

	joinTeam.join(u)

	return joinTeam, nil
}

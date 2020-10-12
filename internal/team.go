package racers

import (
	"errors"
	"fmt"

	"github.com/xabi93/racers/internal/id"
	baseid "github.com/xabi93/racers/internal/id"
)

// Events
type (
	TeamCreated struct {
		BaseEvent   `json:"base_event,omitempty"`
		TeamID      `json:"id,omitempty"`
		TeamName    `json:"name,omitempty"`
		TeamAdmin   UserID      `json:"admin_id,omitempty"`
		TeamMembers TeamMembers `json:"members,omitempty"`
	}
	UserJoinedTeam struct {
		BaseEvent `json:"base_event,omitempty"`
		TeamID    `json:"team_id,omitempty"`
		UserID    `json:"user_id,omitempty"`
	}
)

type (
	TeamID             id.ID
	InvalidTeamIDError struct{ error }
)

func (err InvalidTeamIDError) Error() string {
	return fmt.Sprintf("invalid team id: %s", err.error)
}

func NewTeamID(s string) (TeamID, error) {
	id, err := id.NewID(s)
	if err != nil {
		return "", InvalidTeamIDError{err}
	}

	return TeamID(id), nil
}

type (
	TeamName             string
	InvalidTeamNameError struct{ error }
)

func (err InvalidTeamNameError) Error() string {
	return fmt.Sprintf("invalid team name: %s", err.error)
}

func NewTeamName(s string) (TeamName, error) {
	if s == "" {
		return "", InvalidTeamNameError{errors.New("empty name")}
	}

	return TeamName(s), nil
}

type TeamMembers struct{ userList }

type TeamOption func(*Team)

func TeamMembersOpt(c TeamMembers) TeamOption {
	return func(r *Team) {
		r.members = c
	}
}

func NewTeam(id TeamID, name TeamName, admin UserID, opts ...TeamOption) Team {
	r := Team{
		aggregate: newAggregate(),
		id:        id,
		name:      name,
		admin:     admin,
	}
	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func CreateTeam(id TeamID, name TeamName, admin User) Team {
	t := NewTeam(id, name, admin.ID(), TeamMembersOpt(NewTeamMembers(admin.id)))

	t.aggregate.record(TeamCreated{
		NewBaseEvent(baseid.ID(t.id)),
		t.id,
		t.name,
		t.admin,
		t.members,
	})

	return t
}

func NewTeamMembers(users ...UserID) TeamMembers {
	ul := make(userList, len(users))
	for _, u := range users {
		ul.add(u)
	}

	return TeamMembers{ul}
}

type Team struct {
	aggregate

	id      TeamID
	name    TeamName
	admin   UserID
	members TeamMembers
}

func (t Team) ID() TeamID {
	return t.id
}

func (t Team) Name() TeamName {
	return t.name
}

func (t Team) Admin() UserID {
	return t.admin
}

func (t Team) Members() TeamMembers {
	return t.members
}

type UserAlreadyInTeam struct {
	UserID UserID
	TeamID TeamID
}

func (err UserAlreadyInTeam) Error() string {
	return fmt.Sprintf("user %s cant join is member of %s team", err.UserID, err.TeamID)
}

func (t *Team) Join(tm TeamMember) error {
	if uteam := tm.t; uteam != nil {
		return UserAlreadyInTeam{tm.u.id, uteam.id}
	}

	t.members.add(tm.u.id)

	t.record(UserJoinedTeam{
		BaseEvent: NewBaseEvent(id.ID(t.id)),
		UserID:    tm.u.id,
		TeamID:    t.id,
	})

	return nil
}

func (t *Team) ConsumeEvents() []Event {
	return t.aggregate.ConsumeEvents()
}

func NewTeamMember(u User, t *Team) TeamMember {
	return TeamMember{u, t}
}

type TeamMember struct {
	u User
	t *Team
}

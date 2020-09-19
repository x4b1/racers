package racers

import (
	"errors"
	"fmt"

	"github.com/xabi93/go-clean/internal/types"
)

type (
	UserCreated struct {
		BaseEvent
		UserID   `json:"user_id,omitempty"`
		UserName `json:"user_name,omitempty"`
	}
)

type (
	UserID             types.ID
	InvalidUserIDError struct{ error }
)

func (err InvalidUserIDError) Error() string {
	return fmt.Sprintf("invalid user id: %s", err.error)
}

func NewUserID(s string) (UserID, error) {
	id, err := types.NewID(s)
	if err != nil {
		return "", InvalidUserIDError{err}
	}

	return UserID(id), nil
}

type (
	UserName             string
	InvalidUserNameError struct{ error }
)

func (err InvalidUserNameError) Error() string {
	return fmt.Sprintf("invalid user name: %s", err.error)
}

func NewUserName(s string) (UserName, error) {
	if s == "" {
		return "", InvalidUserNameError{errors.New("empty name")}
	}

	return UserName(s), nil
}

type userList map[UserID]struct{}

func (ul userList) is(id UserID) bool {
	_, ok := ul[id]

	return ok
}

func (ul *userList) add(id UserID) {
	if *ul == nil {
		*ul = make(map[UserID]struct{})
	}

	(*ul)[id] = struct{}{}
}

func (ul userList) List() []UserID {
	l := make([]UserID, 0, len(ul))
	for u := range ul {
		l = append(l, u)
	}

	return l
}

func NewUser(id UserID, name UserName) User {
	return User{newAggregate(), id, name}
}

func CreateUser(id UserID, name UserName) User {
	u := NewUser(id, name)

	u.aggregate.record(UserCreated{
		NewBaseEvent(types.ID(u.id)),
		u.id,
		u.name,
	})

	return u
}

type User struct {
	aggregate

	id   UserID
	name UserName
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Name() UserName {
	return u.name
}

func (u *User) ConsumeEvents() []Event {
	return u.aggregate.ConsumeEvents()
}

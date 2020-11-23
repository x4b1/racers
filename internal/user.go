package racers

import (
	"errors"
	"fmt"

	"github.com/xabi93/racers/internal/id"
)

type (
	// UserID defines a unique identifier for a user
	UserID id.ID
	// InvalidUserIDError means the given user id is not correct
	InvalidUserIDError struct{ error }
)

func (err InvalidUserIDError) Error() string {
	return fmt.Sprintf("invalid user id: %s", err.error)
}

// NewUserID parses a id and returns a UserID if it's ok
func NewUserID(s string) (UserID, error) {
	id, err := id.NewID(s)
	if err != nil {
		return UserID{}, InvalidUserIDError{err}
	}

	return UserID(id), nil
}

// userList is a list of users
type userList map[UserID]struct{}

// is returns if a user is in the list
func (ul userList) is(id UserID) bool {
	_, ok := ul[id]

	return ok
}

// add adds a user to the list
func (ul *userList) add(id UserID) {
	if *ul == nil {
		*ul = make(map[UserID]struct{})
	}

	(*ul)[id] = struct{}{}
}

// List returns a list of users
func (ul userList) List() []UserID {
	l := make([]UserID, 0, len(ul))
	for u := range ul {
		l = append(l, u)
	}

	return l
}

// ErrUnknownUser means a user does not exists in the service
var ErrUnknownUser = errors.New("unknown user")

// User represents a user in the service
type User struct {
	ID UserID
}

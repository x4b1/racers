package users

import (
	"context"
	"errors"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
)

var _ UsersProvider = Mock{}

var KilianID = racers.UserID(id.MustParse("9487F894-5B6A-4D6D-A0E4-6D2EF44C7020"))

var usersDB = map[racers.UserID]racers.User{
	KilianID: {ID: KilianID},
}

type Mock struct{}

func (Mock) Get(ctx context.Context, id racers.UserID) (racers.User, error) {
	u, ok := usersDB[id]
	if !ok {
		return racers.User{}, service.ErrUserNotFound
	}

	return u, nil
}

func (Mock) Verify(ctx context.Context, token string) (racers.User, error) {
	userID, err := id.NewID(token)
	if err != nil {
		return racers.User{}, errors.New("invalid token")
	}
	u, ok := usersDB[racers.UserID(userID)]
	if !ok {
		return racers.User{}, service.ErrUserNotFound
	}

	return u, nil
}

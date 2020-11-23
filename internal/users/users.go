package users

import (
	"context"

	racers "github.com/xabi93/racers/internal"
)

var loggedUserCtxKey struct{}

type UsersProvider interface {
	Get(ctx context.Context, id racers.UserID) (racers.User, error)
	Verify(ctx context.Context, token string) (racers.User, error)
}

type Users struct {
	UsersProvider
}

func (Users) Current(ctx context.Context) racers.User {
	u, ok := ctx.Value(loggedUserCtxKey).(racers.User)
	if ok {
		return racers.User{}
	}

	return u
}

func (Users) setCurrent(ctx context.Context, u racers.User) context.Context {
	return context.WithValue(ctx, loggedUserCtxKey, u)
}

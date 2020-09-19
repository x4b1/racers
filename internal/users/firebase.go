package users

import (
	"context"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"

	"firebase.google.com/go/auth"
)

func NewFirebase(c *auth.Client) Firebase {
	return Firebase{c}
}

type Firebase struct {
	cli *auth.Client
}

func (f Firebase) Get(ctx context.Context, userID racers.UserID) (racers.User, error) {
	_, err := f.cli.GetUser(ctx, id.ID(userID).String())
	if auth.IsUserNotFound(err) {
		return racers.User{}, service.ErrUserNotFound
	}
	if err != nil {
		return racers.User{}, err
	}

	return racers.User{ID: userID}, nil
}

func (f Firebase) Verify(ctx context.Context, token string) (racers.User, error) {
	t, err := f.cli.VerifyIDToken(ctx, token)
	if err != nil {
		return racers.User{}, err
	}

	id, err := id.NewID(t.UID)
	if err != nil {
		return racers.User{}, err
	}

	return racers.User{ID: racers.UserID(id)}, nil
}

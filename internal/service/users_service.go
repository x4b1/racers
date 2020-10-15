package service

import (
	"context"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
)

func NewUsersService(users UsersRepository) UsersService {
	return UsersService{users}
}

type UsersService struct {
	users UsersRepository
}

type CreateUser struct {
	ID   string
	Name string
}

func (ls UsersService) Create(ctx context.Context, r CreateUser) error {
	var ve errors.ValidationError

	id, err := racers.NewUserID(r.ID)
	if err != nil {
		ve.Add(err)
	}
	name, err := racers.NewUserName(r.Name)
	if err != nil {
		ve.Add(err)
	}
	if err := ve.Valid(); err != nil {
		return ve
	}

	if err := ls.users.Save(ctx, racers.CreateUser(id, name)); err != nil {
		return errors.WrapInternalError(err)
	}

	return nil
}

package service_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/xabi93/go-clean/internal/errors"
	"github.com/xabi93/go-clean/internal/types"
	"github.com/xabi93/go-clean/racers"
	"github.com/xabi93/go-clean/racers/service"
)

func TestUserService(t *testing.T) {
	type suite struct {
		ctx     context.Context
		service service.UsersService
		users   *UsersRepositoryMock
	}

	newSuite := func(t *testing.T) (suite, *is.I) {
		var users UsersRepositoryMock

		return suite{
			context.Background(),
			service.NewUsersService(&users),
			&users,
		}, is.New(t)
	}

	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		s, is := newSuite(t)

		err := s.service.Create(s.ctx, service.CreateUser{
			ID:   "",
			Name: "",
		})

		var ve errors.ValidationError
		is.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidUserIDError{},
			racers.InvalidUserNameError{},
		}
		is.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			is.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, and fails to save returns error", func(t *testing.T) {
		s, is := newSuite(t)

		s.users.SaveFunc = func(context.Context, racers.User) error {
			return errors.New("")
		}

		err := s.service.Create(s.ctx, service.CreateUser{
			ID:   types.GenerateID().String(),
			Name: "world champion",
		})

		is.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, and success saving returns no error", func(t *testing.T) {
		s, is := newSuite(t)

		s.users.SaveFunc = func(context.Context, racers.User) error {
			return nil
		}

		id := types.GenerateID()
		r := service.CreateUser{
			ID:   id.String(),
			Name: "world champion",
		}
		err := s.service.Create(s.ctx, r)

		is.NoErr(err)

		is.Equal(len(s.users.SaveCalls()), 1)
		savedUser := s.users.SaveCalls()[0].User

		events := savedUser.ConsumeEvents()
		is.Equal(len(events), 1)
		_, isCreated := events[0].(racers.UserCreated)
		is.True(isCreated)

		is.Equal(racers.NewUser(racers.UserID(id), racers.UserName(r.Name)), savedUser)
	})
}

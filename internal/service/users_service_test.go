package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
)

type testUsersService struct {
	service service.UsersService
	users   *UsersRepositoryMock
}

func newTestUsersService() testUsersService {
	s := testUsersService{
		users: &UsersRepositoryMock{},
	}

	s.service = service.NewUsersService(s.users)

	return s
}

func TestUserService(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	t.Run("When invalid request, returns validation errors", func(t *testing.T) {
		s := newTestUsersService()

		err := s.service.Create(ctx, service.CreateUser{
			ID:   "",
			Name: "",
		})

		var ve errors.ValidationError
		require.True(errors.As(err, &ve))

		expectedErrs := []error{
			racers.InvalidUserIDError{},
			racers.InvalidUserNameError{},
		}
		require.Equal(len(ve.Errors), len(expectedErrs))
		for i, err := range ve.Errors {
			require.True(errors.As(err, &expectedErrs[i]))
		}
	})

	t.Run("When valid request, and fails to save returns error", func(t *testing.T) {
		s := newTestUsersService()

		s.users.SaveFunc = func(context.Context, racers.User) error {
			return errors.New("")
		}

		err := s.service.Create(ctx, service.CreateUser{
			ID:   id.GenerateID().String(),
			Name: "world champion",
		})

		require.True(errors.IsInternalError(err))
	})

	t.Run("When valid request, and success saving returns no error", func(t *testing.T) {
		s := newTestUsersService()

		s.users.SaveFunc = func(context.Context, racers.User) error {
			return nil
		}

		id := id.GenerateID()
		r := service.CreateUser{
			ID:   id.String(),
			Name: "world champion",
		}
		err := s.service.Create(ctx, r)

		require.NoError(err)

		require.Equal(len(s.users.SaveCalls()), 1)
		savedUser := s.users.SaveCalls()[0].User

		events := savedUser.ConsumeEvents()
		require.Equal(len(events), 1)
		_, isCreated := events[0].(racers.UserCreated)
		require.True(isCreated)

		require.Equal(racers.NewUser(racers.UserID(id), racers.UserName(r.Name)), savedUser)
	})
}

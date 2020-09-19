package racers_test

import (
	"errors"
	"testing"

	"github.com/matryer/is"
	"github.com/xabi93/go-clean/internal/types"
	"github.com/xabi93/go-clean/racers"
)

var (
	userID   = racers.UserID(types.GenerateID())
	userName = racers.UserName("kilian Jornet")
)

func TestUserID(t *testing.T) {
	is := is.New(t)
	t.Run("when invalid id returns InvalidUserIDError error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewUserID("")
		is.True(errors.As(err, &racers.InvalidUserIDError{}))
	})

	t.Run("when valid id returns UserID and no error", func(t *testing.T) {
		is := is.New(t)

		id := types.GenerateID()
		userID, err := racers.NewUserID(id.String())

		is.Equal(racers.UserID(id), userID)
		is.NoErr(err)
	})
}

func TestUserName(t *testing.T) {
	is := is.New(t)
	t.Run("when New with empty name returns InvalidUserNameError error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewUserName("")
		is.True(errors.As(err, &racers.InvalidUserNameError{}))
	})

	t.Run("when New with valid name returns UserName and no error", func(t *testing.T) {
		is := is.New(t)

		name := "athletic"
		userName, err := racers.NewUserName(name)

		is.Equal(racers.UserName(name), userName)
		is.NoErr(err)
	})
}

func TestNewUser(t *testing.T) {
	is := is.New(t)

	user := racers.NewUser(userID, userName)

	is.Equal(user.ID(), userID)
	is.Equal(user.Name(), userName)

	events := user.ConsumeEvents()
	is.True(len(events) == 0)
}

func TestCreateUser(t *testing.T) {
	is := is.New(t)

	r := racers.CreateUser(userID, userName)

	events := r.ConsumeEvents()
	is.True(len(events) == 1)
	createdEvent := events[0].(racers.UserCreated)
	is.Equal(createdEvent.UserID, userID)
	is.Equal(createdEvent.UserName, userName)

	is.Equal(r, racers.NewUser(userID, userName))
}

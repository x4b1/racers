package racers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
)

var (
	userID   = racers.UserID(id.GenerateID())
	userName = racers.UserName("kilian Jornet")
)

func TestUserID(t *testing.T) {
	require := require.New(t)
	t.Run("when invalid id returns InvalidUserIDError error", func(t *testing.T) {
		_, err := racers.NewUserID("")
		require.True(errors.As(err, &racers.InvalidUserIDError{}))
	})

	t.Run("when valid id returns UserID and no error", func(t *testing.T) {
		id := id.GenerateID()
		userID, err := racers.NewUserID(id.String())

		require.Equal(racers.UserID(id), userID)
		require.NoError(err)
	})
}

func TestUserName(t *testing.T) {
	require := require.New(t)
	t.Run("when New with empty name returns InvalidUserNameError error", func(t *testing.T) {
		_, err := racers.NewUserName("")
		require.True(errors.As(err, &racers.InvalidUserNameError{}))
	})

	t.Run("when New with valid name returns UserName and no error", func(t *testing.T) {
		name := "athletic"
		userName, err := racers.NewUserName(name)

		require.Equal(racers.UserName(name), userName)
		require.NoError(err)
	})
}

func TestNewUser(t *testing.T) {
	require := require.New(t)

	user := racers.NewUser(userID, userName)

	require.Equal(user.ID(), userID)
	require.Equal(user.Name(), userName)

	require.Empty(user.ConsumeEvents())
}

func TestCreateUser(t *testing.T) {
	require := require.New(t)

	r := racers.CreateUser(userID, userName)

	events := r.ConsumeEvents()
	require.Len(events, 1)
	createdEvent := events[0].(racers.UserCreated)
	require.Equal(createdEvent.UserID, userID)
	require.Equal(createdEvent.UserName, userName)

	require.Equal(r, racers.NewUser(userID, userName))
}

package racers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
)

var userID = racers.UserID(id.Generate())

func TestUserID(t *testing.T) {
	require := require.New(t)
	t.Run("when invalid id returns InvalidUserIDError error", func(t *testing.T) {
		_, err := racers.NewUserID("")
		require.True(errors.As(err, &racers.InvalidUserIDError{}))
	})

	t.Run("when valid id returns UserID and no error", func(t *testing.T) {
		id := id.Generate()
		userID, err := racers.NewUserID(id.String())

		require.Equal(racers.UserID(id), userID)
		require.NoError(err)
	})
}

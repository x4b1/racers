package racers_test

import (
	"errors"
	"testing"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"

	"github.com/stretchr/testify/require"
)

var (
	raceID         = racers.RaceID(id.Generate())
	raceName       = racers.RaceName("New York Marathon")
	raceDate       = racers.RaceDate(time.Now())
	raceCompetitor = racers.User{racers.UserID(id.Generate())}
	ownerID        = racers.UserID(id.Generate())
)

func TestRaceID(t *testing.T) {
	require := require.New(t)
	t.Run("when invalid id returns InvalidRaceIDError error", func(t *testing.T) {
		_, err := racers.NewRaceID("")
		require.True(errors.As(err, &racers.InvalidRaceIDError{}))
	})

	t.Run("when valid id returns RaceID and no error", func(t *testing.T) {
		id := id.Generate()
		raceID, err := racers.NewRaceID(id.String())

		require.Equal(racers.RaceID(id), raceID)
		require.NoError(err)
	})
}

func TestRaceName(t *testing.T) {
	require := require.New(t)
	t.Run("when New with empty name returns InvalidRaceNameError error", func(t *testing.T) {
		_, err := racers.NewRaceName("")
		require.True(errors.As(err, &racers.InvalidRaceNameError{}))
	})

	t.Run("when New with valid name returns RaceName and no error", func(t *testing.T) {
		name := "New York Marathon"
		raceName, err := racers.NewRaceName(name)

		require.Equal(racers.RaceName(name), raceName)
		require.NoError(err)
	})
}

func TestRaceDate(t *testing.T) {
	require := require.New(t)
	t.Run("when New with past date, returns error", func(t *testing.T) {
		_, err := racers.NewRaceDate(time.Now().AddDate(0, 0, -1))
		require.True(errors.As(err, &racers.InvalidRaceDateError{}))
	})

	t.Run("when New with valid date returns RaceName and no error", func(t *testing.T) {
		tomorrow := time.Now().AddDate(0, 0, 1)
		date, err := racers.NewRaceDate(tomorrow)

		require.Equal(racers.RaceDate(tomorrow), date)
		require.NoError(err)
	})
}

func TestRaceCompetitors(t *testing.T) {
	require := require.New(t)
	usersIDs := []racers.UserID{
		racers.UserID(id.Generate()),
		racers.UserID(id.Generate()),
	}

	rc := racers.NewRaceCompetitors(usersIDs...)

	t.Run(`Given a list of competitors
	when List returns a list of users`, func(t *testing.T) {
		require.Equal(len(usersIDs), len(rc.List()))
	})
}

func TestRace(t *testing.T) {
	require := require.New(t)

	t.Run(`Given a race with no competitors,
	When joins one,
	Then returns no error and generates RaceCompetitorJoined event`, func(t *testing.T) {
		r := racers.Race{
			ID:    raceID,
			Name:  raceName,
			Date:  raceDate,
			Owner: ownerID,
		}
		require.NoError(r.Join(raceCompetitor))
	})

	t.Run(`Given a race with one competitor,
	When tries to join the same competitor,
	Then returns CompetitorInRaceError error`, func(t *testing.T) {
		r := racers.Race{
			ID:          raceID,
			Name:        raceName,
			Date:        raceDate,
			Owner:       ownerID,
			Competitors: racers.NewRaceCompetitors(raceCompetitor.ID),
		}

		err := r.Join(raceCompetitor)

		var competirorInRaceErr racers.CompetitorInRaceError
		require.True(errors.As(err, &competirorInRaceErr))
		require.Equal(competirorInRaceErr.RaceID, r.ID)
		require.Equal(competirorInRaceErr.CompetitorID, raceCompetitor.ID)
	})
}

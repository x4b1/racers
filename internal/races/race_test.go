package races_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/common/id"
)

var (
	raceID         = racers.RaceID(id.GenerateID())
	raceName       = racers.RaceName("New York Marathon")
	raceDate       = racers.RaceDate(time.Now())
	raceCompetitor = racers.NewUser(racers.UserID(id.GenerateID()), racers.UserName("Usaint"))
)

func TestRaceID(t *testing.T) {
	require := require.New(t)
	t.Run("when invalid id returns InvalidRaceIDError error", func(t *testing.T) {
		_, err := racers.NewRaceID("")
		require.True(errors.As(err, &racers.InvalidRaceIDError{}))
	})

	t.Run("when valid id returns RaceID and no error", func(t *testing.T) {
		id := id.GenerateID()
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
		racers.UserID(id.GenerateID()),
		racers.UserID(id.GenerateID()),
	}

	rc := racers.NewRaceCompetitors(usersIDs...)

	t.Run(`Given a list of competitors
	when List returns a list of users`, func(t *testing.T) {
		require.Equal(len(usersIDs), len(rc.List()))
	})
}

func TestNewRace(t *testing.T) {
	require := require.New(t)

	competitors := racers.NewRaceCompetitors(
		racers.UserID(id.GenerateID()),
		racers.UserID(id.GenerateID()),
	)

	r := racers.NewRace(raceID, raceName, raceDate, racers.RaceCompetitorsOpt(competitors))

	require.Equal(r.ID(), raceID)
	require.Equal(r.Name(), raceName)
	require.Equal(r.Competitors(), competitors)

	require.Empty(r.ConsumeEvents())
}

func TestCreateRace(t *testing.T) {
	require := require.New(t)

	r := racers.CreateRace(raceID, raceName, raceDate)

	events := r.ConsumeEvents()
	require.Len(events, 1)
	createdEvent := events[0].(racers.RaceCreated)
	require.Equal(createdEvent.RaceID, raceID)
	require.Equal(createdEvent.RaceName, raceName)
	require.Equal(createdEvent.RaceDate, raceDate)

	require.Equal(r, racers.NewRace(raceID, raceName, raceDate))
}

func TestRace(t *testing.T) {
	require := require.New(t)

	t.Run(`Given a race with no competitors,
	When joins one,
	Then returns no error and generates RaceCompetitorJoined event`, func(t *testing.T) {
		r := racers.NewRace(
			raceID,
			raceName,
			raceDate,
		)
		require.NoError(r.Join(raceCompetitor))

		events := r.ConsumeEvents()
		require.Len(events, 1)

		joinedEvent := events[0].(racers.RaceCompetitorJoined)
		require.Equal(joinedEvent.RaceID, r.ID())
		require.Equal(joinedEvent.CompetitorID, raceCompetitor.ID())
	})

	t.Run(`Given a race with one competitor,
	When tries to join the same competitor,
	Then returns CompetitorInRaceError error`, func(t *testing.T) {
		r := racers.NewRace(
			raceID,
			raceName,
			raceDate,
			racers.RaceCompetitorsOpt(racers.NewRaceCompetitors(raceCompetitor.ID())),
		)

		err := r.Join(raceCompetitor)

		var competirorInRaceErr racers.CompetitorInRaceError
		require.True(errors.As(err, &competirorInRaceErr))
		require.Equal(competirorInRaceErr.RaceID, r.ID())
		require.Equal(competirorInRaceErr.CompetitorID, raceCompetitor.ID())

		require.Empty(r.ConsumeEvents())
	})
}

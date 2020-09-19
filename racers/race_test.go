package racers_test

import (
	"errors"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/xabi93/go-clean/internal/types"
	"github.com/xabi93/go-clean/racers"
)

var (
	raceID         = racers.RaceID(types.GenerateID())
	raceName       = racers.RaceName("New York Marathon")
	raceDate       = racers.RaceDate(time.Now())
	raceCompetitor = racers.NewUser(racers.UserID(types.GenerateID()), racers.UserName("Usaint"))
)

func TestRaceID(t *testing.T) {
	is := is.New(t)
	t.Run("when invalid id returns InvalidRaceIDError error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewRaceID("")
		is.True(errors.As(err, &racers.InvalidRaceIDError{}))
	})

	t.Run("when valid id returns RaceID and no error", func(t *testing.T) {
		is := is.New(t)

		id := types.GenerateID()
		raceID, err := racers.NewRaceID(id.String())

		is.Equal(racers.RaceID(id), raceID)
		is.NoErr(err)
	})
}

func TestRaceName(t *testing.T) {
	is := is.New(t)
	t.Run("when New with empty name returns InvalidRaceNameError error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewRaceName("")
		is.True(errors.As(err, &racers.InvalidRaceNameError{}))
	})

	t.Run("when New with valid name returns RaceName and no error", func(t *testing.T) {
		is := is.New(t)

		name := "New York Marathon"
		raceName, err := racers.NewRaceName(name)

		is.Equal(racers.RaceName(name), raceName)
		is.NoErr(err)
	})
}

func TestRaceDate(t *testing.T) {
	is := is.New(t)
	t.Run("when New with past date, returns error", func(t *testing.T) {
		is := is.New(t)
		_, err := racers.NewRaceDate(time.Now().AddDate(0, 0, -1))
		is.True(errors.As(err, &racers.InvalidRaceDateError{}))
	})

	t.Run("when New with valid date returns RaceName and no error", func(t *testing.T) {
		is := is.New(t)

		tomorrow := time.Now().AddDate(0, 0, 1)
		date, err := racers.NewRaceDate(tomorrow)

		is.Equal(racers.RaceDate(tomorrow), date)
		is.NoErr(err)
	})
}

func TestRaceCompetitors(t *testing.T) {
	is := is.New(t)
	usersIDs := []racers.UserID{
		racers.UserID(types.GenerateID()),
		racers.UserID(types.GenerateID()),
	}

	rc := racers.NewRaceCompetitors(usersIDs...)

	t.Run(`Given a list of competitors
	when List returns a list of users`, func(t *testing.T) {
		is := is.New(t)
		is.Equal(len(usersIDs), len(rc.List()))
	})
}

func TestNewRace(t *testing.T) {
	is := is.New(t)

	competitors := racers.NewRaceCompetitors(
		racers.UserID(types.GenerateID()),
		racers.UserID(types.GenerateID()),
	)

	r := racers.NewRace(raceID, raceName, raceDate, racers.RaceCompetitorsOpt(competitors))

	is.Equal(r.ID(), raceID)
	is.Equal(r.Name(), raceName)
	is.Equal(r.Competitors(), competitors)

	events := r.ConsumeEvents()
	is.True(len(events) == 0)
}

func TestCreateRace(t *testing.T) {
	is := is.New(t)

	r := racers.CreateRace(raceID, raceName, raceDate)

	events := r.ConsumeEvents()
	is.True(len(events) == 1)
	createdEvent := events[0].(racers.RaceCreated)
	is.Equal(createdEvent.RaceID, raceID)
	is.Equal(createdEvent.RaceName, raceName)
	is.Equal(createdEvent.RaceDate, raceDate)

	is.Equal(r, racers.NewRace(raceID, raceName, raceDate))
}

func TestRace(t *testing.T) {
	is := is.New(t)

	t.Run(`Given a race with no competitors,
	When joins one,
	Then returns no error and generates RaceCompetitorJoined event`, func(t *testing.T) {
		is := is.New(t)

		r := racers.NewRace(
			raceID,
			raceName,
			raceDate,
		)
		is.NoErr(r.Join(raceCompetitor))

		events := r.ConsumeEvents()
		is.True(len(events) == 1)

		joinedEvent := events[0].(racers.RaceCompetitorJoined)
		is.Equal(joinedEvent.RaceID, r.ID())
		is.Equal(joinedEvent.CompetitorID, raceCompetitor.ID())
	})

	t.Run(`Given a race with one competitor,
	When tries to join the same competitor,
	Then returns CompetitorInRaceError error`, func(t *testing.T) {
		is := is.New(t)

		r := racers.NewRace(
			raceID,
			raceName,
			raceDate,
			racers.RaceCompetitorsOpt(racers.NewRaceCompetitors(raceCompetitor.ID())),
		)

		err := r.Join(raceCompetitor)

		var competirorInRaceErr racers.CompetitorInRaceError
		is.True(errors.As(err, &competirorInRaceErr))
		is.Equal(competirorInRaceErr.RaceID, r.ID())
		is.Equal(competirorInRaceErr.CompetitorID, raceCompetitor.ID())

		is.True(len(r.ConsumeEvents()) == 0)
	})
}

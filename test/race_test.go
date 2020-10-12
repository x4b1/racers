package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"
	"github.com/xabi93/racers/internal/api/graph"
	"github.com/xabi93/racers/internal/id"
)

// fixtures
var (
	raceName = "super-race"
	raceDate = time.Now().AddDate(0, 1, 2).Truncate(1 * time.Second).UTC()
)

func TestCreateRace(t *testing.T) {
	require := require.New(t)
	s := newSuite(t)

	t.Run("invalid payload", func(t *testing.T) {
		type testCase struct {
			name      string
			date      time.Time
			errorType interface{}
		}
		testCases := map[string]testCase{
			"empty name":   {date: raceDate, errorType: graph.InvalidNameError{}},
			"invalid date": {name: raceName, errorType: graph.InvalidDateError{}},
		}
		for name, c := range testCases {
			t.Run(fmt.Sprintf("when %s", name), func(t *testing.T) {
				var race graph.Race
				race.Name = c.name
				race.Date = c.date
				resp := createRace(s.graphql, race)

				require.Equal(reflect.TypeOf(c.errorType).Name(), resp.CreateRace.Typename)
				require.NotEmpty(resp.CreateRace.Message)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		var race graph.Race
		race.Name = raceName
		race.Date = raceDate

		resp := createRace(s.graphql, race)

		require.Equal(raceName, resp.CreateRace.Name)

		respDate, err := graphql.UnmarshalTime(resp.CreateRace.Date)
		require.NoError(err)
		require.Equal(raceDate, respDate)
	})
}

func TestGetRace(t *testing.T) {
	require := require.New(t)

	t.Run("not exists", func(t *testing.T) {
		s := newSuite(t)
		raceID := id.GenerateID()
		resp := getRace(s.graphql, raceID)

		require.Equal(reflect.TypeOf(graph.RaceByIDNotFound{}).Name(), resp.Race.Typename)
		require.Equal(raceID.String(), resp.Race.ID)
	})

	t.Run("exists", func(t *testing.T) {
		s := newSuite(t)
		var race graph.Race
		race.Name = raceName
		race.Date = raceDate

		create := createRace(s.graphql, race)
		raceID, err := id.NewID(create.CreateRace.ID)
		require.NoError(err)
		resp := getRace(s.graphql, raceID)

		require.Equal(reflect.TypeOf(graph.Race{}).Name(), resp.Race.Typename)
		require.Equal(raceID.String(), resp.Race.ID)
		require.Equal(raceName, resp.Race.Name)
		respDate, err := graphql.UnmarshalTime(resp.Race.Date)
		require.NoError(err)
		require.Equal(raceDate, respDate)
	})
}

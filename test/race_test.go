package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/server/graph"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"
)

// fixtures
var (
	blackMambaRace = graph.Race{
		ID:   id.Generate().String(),
		Name: "black mamba race",
		Date: time.Now().AddDate(0, 1, 2).Truncate(1 * time.Second).UTC(),
	}
)

func TestCreateRace(t *testing.T) {
	require := require.New(t)
	s := newSuite(t)
	defer s.db.Close()

	t.Run("invalid payload", func(t *testing.T) {
		type testCase struct {
			reqFactory func(graph.Race) graph.Race
			errorType  interface{}
		}

		testCases := map[string]testCase{
			"empty name":   {errorType: graph.InvalidRaceNameError{}, reqFactory: func(r graph.Race) graph.Race { r.Name = ""; return r }},
			"invalid date": {errorType: graph.InvalidRaceDateError{}, reqFactory: func(r graph.Race) graph.Race { r.Date = time.Time{}; return r }},
		}
		for name, c := range testCases {
			t.Run(fmt.Sprintf("when %s", name), func(t *testing.T) {
				resp := createRace(s.graphql, c.reqFactory(blackMambaRace))

				require.Equal(reflect.TypeOf(c.errorType).Name(), resp.CreateRace.Typename)
				require.NotEmpty(resp.CreateRace.Message)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		resp := createRace(s.graphql, blackMambaRace)

		require.Equal(blackMambaRace.Name, resp.CreateRace.Name)

		respDate, err := graphql.UnmarshalTime(resp.CreateRace.Date)
		require.NoError(err)
		require.Equal(blackMambaRace.Date, respDate)
	})

	t.Run("duplicated", func(t *testing.T) {
		resp := createRace(s.graphql, blackMambaRace)

		require.Equal(reflect.TypeOf(graph.RaceAlreadyExists{}).Name(), resp.CreateRace.Typename)
	})
}

func TestGetRace(t *testing.T) {
	require := require.New(t)

	s := newSuite(t)
	defer s.db.Close()

	t.Run("not exists", func(t *testing.T) {
		resp := getRace(s.graphql, id.MustParse(blackMambaRace.ID))

		require.Equal(reflect.TypeOf(graph.RaceNotFound{}).Name(), resp.Race.Typename)
	})

	t.Run("exists", func(t *testing.T) {
		createRace(s.graphql, blackMambaRace)

		resp := getRace(s.graphql, id.MustParse(blackMambaRace.ID))

		require.Equal(reflect.TypeOf(graph.Race{}).Name(), resp.Race.Typename)
		require.Equal(blackMambaRace.Name, resp.Race.Name)
		respDate, err := graphql.UnmarshalTime(resp.Race.Date)
		require.NoError(err)
		require.Equal(blackMambaRace.Date, respDate)
	})
}

func TestAllRaces(t *testing.T) {
	require := require.New(t)
	s := newSuite(t)
	defer s.db.Close()

	t.Run("exists", func(t *testing.T) {
		createRace(s.graphql, blackMambaRace)

		resp := getRace(s.graphql, id.MustParse(blackMambaRace.ID))

		require.Equal(reflect.TypeOf(graph.Race{}).Name(), resp.Race.Typename)
		require.Equal(blackMambaRace.Name, resp.Race.Name)
		respDate, err := graphql.UnmarshalTime(resp.Race.Date)
		require.NoError(err)
		require.Equal(blackMambaRace.Date, respDate)
	})
}

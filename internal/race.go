package racers

import (
	"fmt"
	"time"

	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
)

type (
	RaceID             id.ID
	InvalidRaceIDError struct{ error }
)

func (err InvalidRaceIDError) Error() string {
	return fmt.Sprintf("invalid race id: %s", err.error)
}

func NewRaceID(s string) (RaceID, error) {
	id, err := id.NewID(s)
	if err != nil {
		return RaceID{}, InvalidRaceIDError{err}
	}

	return RaceID(id), nil
}

type (
	RaceName             string
	InvalidRaceNameError struct{ error }
)

func (err InvalidRaceNameError) Error() string {
	return fmt.Sprintf("invalid race name: %s", err.error)
}

func NewRaceName(s string) (RaceName, error) {
	if s == "" {
		return "", InvalidRaceNameError{errors.New("empty name")}
	}

	return RaceName(s), nil
}

type (
	RaceDate             time.Time
	InvalidRaceDateError struct{ time.Time }
)

func (err InvalidRaceDateError) Error() string {
	return fmt.Sprintf("race date cannot be past: %s", err.Time)
}

func NewRaceDate(t time.Time) (RaceDate, error) {
	if t.Before(time.Now()) {
		return RaceDate{}, InvalidRaceDateError{t}
	}

	return RaceDate(t), nil
}

func NewRaceCompetitors(users ...UserID) RaceCompetitors {
	ul := make(userList, len(users))
	for _, u := range users {
		ul.add(u)
	}

	return RaceCompetitors{ul}
}

type RaceCompetitors struct{ userList }

type Race struct {
	ID          RaceID
	Name        RaceName
	Date        RaceDate
	Owner       UserID
	Competitors RaceCompetitors
}

type CompetitorInRaceError struct {
	RaceID       RaceID
	CompetitorID UserID
}

func (err CompetitorInRaceError) Error() string {
	return fmt.Sprintf("competitor %s already joined race %s", err.CompetitorID, err.RaceID)
}

func (r *Race) Join(u User) error {
	if r.Competitors.is(u.ID) {
		return CompetitorInRaceError{r.ID, u.ID}
	}

	r.Competitors.add(u.ID)

	return nil
}

package racers

import (
	"errors"
	"fmt"
	"time"

	"github.com/xabi93/racers/internal/id"
	baseid "github.com/xabi93/racers/internal/id"
)

// Events
type (
	RaceCreated struct {
		BaseEvent
		RaceID   RaceID   `json:"id,omitempty"`
		RaceName RaceName `json:"name,omitempty"`
		RaceDate RaceDate `json:"date,omitempty"`
	}
	RaceCompetitorJoined struct {
		BaseEvent
		RaceID       RaceID `json:"race_id,omitempty"`
		CompetitorID UserID `json:"competitor_id,omitempty"`
	}
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
		return "", InvalidRaceIDError{err}
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

type RaceOption func(*Race)

func RaceCompetitorsOpt(c RaceCompetitors) RaceOption {
	return func(r *Race) {
		r.competitors = c
	}
}

func NewRace(id RaceID, name RaceName, date RaceDate, opts ...RaceOption) Race {
	r := Race{
		aggregate: newAggregate(),
		id:        id,
		name:      name,
		date:      date,
	}
	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func CreateRace(id RaceID, name RaceName, date RaceDate) Race {
	r := NewRace(id, name, date)

	r.aggregate.record(RaceCreated{
		NewBaseEvent(baseid.ID(r.id)),
		r.id,
		r.name,
		r.date,
	})

	return r
}

type Race struct {
	aggregate

	id          RaceID
	name        RaceName
	date        RaceDate
	competitors RaceCompetitors
}

func (r Race) ID() RaceID {
	return r.id
}

func (r Race) Name() RaceName {
	return r.name
}

func (r Race) Date() RaceDate {
	return r.date
}

func (r Race) Competitors() RaceCompetitors {
	return r.competitors
}

type CompetitorInRaceError struct {
	RaceID       RaceID
	CompetitorID UserID
}

func (err CompetitorInRaceError) Error() string {
	return fmt.Sprintf("competitor %s already joined race %s", err.CompetitorID, err.RaceID)
}

func (r *Race) Join(u User) error {
	if r.competitors.is(u.id) {
		return CompetitorInRaceError{r.id, u.id}
	}

	r.competitors.add(u.id)

	r.record(RaceCompetitorJoined{
		NewBaseEvent(id.ID(r.id)),
		r.id,
		u.id,
	})

	return nil
}

func (r *Race) ConsumeEvents() []Event {
	return r.aggregate.ConsumeEvents()
}

package races

import (
	"errors"
	"fmt"
	"time"

	"github.com/xabi93/racers/internal/common/aggregate"
	"github.com/xabi93/racers/internal/common/id"
)

// Events
type (
	RaceCreated struct {
		aggregate.BaseEvent
		RaceID   RaceID   `json:"id,omitempty"`
		RaceName RaceName `json:"name,omitempty"`
		RaceDate RaceDate `json:"date,omitempty"`
	}
	RaceCompetitorJoined struct {
		aggregate.BaseEvent
		RaceID       RaceID       `json:"race_id,omitempty"`
		CompetitorID CompetitorID `json:"competitor_id,omitempty"`
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

type (
	CompetitorID id.ID
	Competitor   struct {
		id CompetitorID
	}
)

func NewRaceCompetitors(competitors ...CompetitorID) RaceCompetitors {
	rCompetitors := make(RaceCompetitors, len(competitors))
	for _, c := range competitors {
		rCompetitors.add(c)
	}

	return rCompetitors
}

type RaceCompetitors map[CompetitorID]struct{}

func (rc RaceCompetitors) is(c CompetitorID) bool {
	_, is := rc[c]
	return is
}

func (rc RaceCompetitors) add(c CompetitorID) {
	rc[c] = struct{}{}
}

type RaceOption func(*Race)

func RaceCompetitorsOpt(c RaceCompetitors) RaceOption {
	return func(r *Race) {
		r.competitors = c
	}
}

func NewRace(id RaceID, name RaceName, date RaceDate, opts ...RaceOption) Race {
	r := Race{
		aggregate: aggregate.NewAggregate(),
		id:        id,
		name:      name,
		date:      date,
	}
	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func CreateRace(raceID RaceID, name RaceName, date RaceDate) Race {
	r := NewRace(raceID, name, date)

	r.aggregate.Record(RaceCreated{
		aggregate.NewBaseEvent(id.ID(r.id)),
		r.id,
		r.name,
		r.date,
	})

	return r
}

type Race struct {
	aggregate aggregate.Aggregate

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
	CompetitorID CompetitorID
}

func (err CompetitorInRaceError) Error() string {
	return fmt.Sprintf("competitor %s already joined race %s", err.CompetitorID, err.RaceID)
}

func (r *Race) Join(c Competitor) error {
	if r.competitors.is(c.id) {
		return CompetitorInRaceError{r.id, c.id}
	}

	r.competitors.add(c.id)

	r.aggregate.Record(RaceCompetitorJoined{
		aggregate.NewBaseEvent(id.ID(r.id)),
		r.id,
		c.id,
	})

	return nil
}

func (r *Race) ConsumeEvents() []aggregate.Event {
	return r.aggregate.ConsumeEvents()
}

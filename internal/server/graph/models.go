package graph

import (
	"errors"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
)

type Race struct {
	ID             string
	Name           string
	Date           time.Time
	Competitors    []*User
	competitorsIDs []racers.UserID
}

func (Race) IsCreateRaceResult() {}
func (Race) IsRaceResult()       {}

func NewRace(race racers.Race) *Race {
	return &Race{
		ID:             id.ID(race.ID).String(),
		Name:           string(race.Name),
		Date:           time.Time(race.Date),
		competitorsIDs: race.Competitors.List(),
	}
}

func NewInternalError() error {
	return errors.New("internal error")
}

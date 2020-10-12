package service

import (
	"fmt"

	"github.com/xabi93/racers/internal/races"
)

type RaceByIDNotFoundError struct {
	ID races.RaceID
}

func (err RaceByIDNotFoundError) Error() string {
	return fmt.Sprintf("race %s not found", err.ID)
}

type CompetitorByIDNotFoundError struct {
	ID races.CompetitorID
}

func (err CompetitorByIDNotFoundError) Error() string {
	return fmt.Sprintf("competitor %s not found", err.ID)
}

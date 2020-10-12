package service

import (
	"fmt"

	racers "github.com/xabi93/racers/internal"
)

// Races errors
type RaceByIDNotFoundError struct {
	ID racers.RaceID
}

func (err RaceByIDNotFoundError) Error() string {
	return fmt.Sprintf("race %s not found", err.ID)
}

type RaceNameAlreadyExistsError struct {
	Name racers.RaceName
}

func (err RaceNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("race name %s already exits", err.Name)
}

type UserByIDNotFoundError struct {
	ID racers.UserID
}

func (err UserByIDNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", err.ID)
}

type TeamByIDNotFoundError struct {
	ID racers.TeamID
}

func (err TeamByIDNotFoundError) Error() string {
	return fmt.Sprintf("team %s not found", err.ID)
}

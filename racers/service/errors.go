package service

import (
	"fmt"

	"github.com/xabi93/go-clean/racers"
)

type RaceByIDNotFoundError struct {
	ID racers.RaceID
}

func (err RaceByIDNotFoundError) Error() string {
	return fmt.Sprintf("race %s not found", err.ID)
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

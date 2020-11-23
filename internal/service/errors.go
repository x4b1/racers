package service

import "github.com/xabi93/racers/internal/errors"

// Races errors
var (
	ErrRaceNotFound      = errors.New("race not found")
	ErrRaceAlreadyExists = errors.New("race already exists")
)

// Users errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// Teams errors
var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamAlreadyExists = errors.New("team already exists")
)

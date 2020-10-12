package service

import (
	"context"
)

var _ RacesService = RacesMetrics{}

func NewRacesMetrics(rs RacesService) RacesService {
	return RacesMetrics{rs}
}

type RacesMetrics struct {
	next RacesService
}

func (rm RacesMetrics) Create(ctx context.Context, r CreateRace) error {
	return rm.next.Create(ctx, r)
}

func (rm RacesMetrics) Join(ctx context.Context, r JoinRace) error {
	return rm.next.Join(ctx, r)
}

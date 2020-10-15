package service

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xabi93/racers/internal/errors"
)

type ActionName string

func (an ActionName) String() string {
	return string(an)
}

func NewService(races RacesService, teams TeamsService) *Service {
	s := Service{
		races,
		teams,
	}

	return &s
}

type Service struct {
	Races RacesService
	Teams TeamsService
}

func NewInstrumenting(l log.Logger) Instrumenting {
	return Instrumenting{l, nil, nil}
}

type Instrumenting struct {
	logger    log.Logger
	count     prometheus.Counter
	histogram prometheus.Histogram
}

func (i Instrumenting) Log(ctx context.Context, action ActionName, req interface{}, h func() error) {
	var err error

	defer func(begin time.Time) {
		elapsed := time.Since(begin)

		i.count.With("action", action.String()).Add(1)
		i.histogram.With("action", action.String()).Observe(elapsed)

		fields := []interface{}{
			"action", action,
			"request", req,
			"took", elapsed,
		}
		if errors.IsInternalError(err) {
			i.logger.Error(ctx, err, fields...)
			return
		}
		fields = append(fields, "error", err)
		i.logger.Info(ctx, fields)
	}(time.Now())

	err = h()
}

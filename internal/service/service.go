package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/log"
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
	count     metrics.Counter
	histogram metrics.Histogram
}

func (i Instrumenting) Log(ctx context.Context, action ActionName, req interface{}, h func() error) {
	var err error

	defer func(begin time.Time) {
		elapsed := time.Since(begin)

		i.count.With("action", action.String()).Add(1)
		i.histogram.With("action", action.String()).Observe(float64(elapsed))

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

package service

import "context"

type Work func(ctx context.Context) error

type UnitOfWork func(ctx context.Context, work Work) error

func NoopUnitOfWork(ctx context.Context, work Work) error {
	return work(ctx)
}

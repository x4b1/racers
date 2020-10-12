package ent

import (
	"context"
	"database/sql"

	"github.com/xabi93/racers/internal/log"
	"github.com/xabi93/racers/internal/storage/postgres"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
)

func New(conf postgres.Config, logger log.Logger, db *sql.DB) *Client {
	d := dialect.DebugWithContext(entsql.OpenDB(dialect.Postgres, db), func(ctx context.Context, args ...interface{}) {
		logger.Debug(ctx, "ent sql", log.Fields{"payload": args})
	})

	return NewClient(Driver(d))
}

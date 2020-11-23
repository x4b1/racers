package postgres

import (
	"context"
	"database/sql"

	"github.com/xabi93/racers/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var transactionContextKey struct{}

type Repository struct{ db *gorm.DB }

func (r Repository) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionContextKey).(*gorm.DB)
	if !ok {
		return r.db
	}

	return tx
}

func New(conn *sql.DB) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func TransactionFactory(db *gorm.DB) service.UnitOfWork {
	return func(ctx context.Context, w service.Work) error {
		return db.Transaction(func(tx *gorm.DB) error {
			return w(context.WithValue(ctx, transactionContextKey, tx))
		})
	}
}

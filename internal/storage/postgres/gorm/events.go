package gorm

import (
	"context"
	"encoding/json"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
	"gorm.io/gorm"
)

type event struct {
	ID          id.ID     `db:"id"`
	Payload     string    `db:"payload"`
	AggregateID id.ID     `db:"aggregate_id"`
	UserID      id.ID     `db:"user_id"`
	OccurredOn  time.Time `db:"date"`
}

func NewEvents(db *gorm.DB) Events {
	return Events{db, "events"}
}

type Events struct {
	db        *gorm.DB
	tableName string
}

func (e Events) Save(ctx context.Context, tx *gorm.DB, events ...racers.Event) error {
	if tx == nil {
		tx = e.db
	}
	if len(events) == 0 {
		return nil
	}
	eventsDB := make([]event, len(events))
	for i, e := range events {
		payload, err := json.Marshal(e)
		if err != nil {
			return errors.WrapInternalError(err, "marshaling event")
		}
		eventsDB[i] = event{
			ID:          e.ID(),
			Payload:     string(payload),
			AggregateID: e.AggregateID(),
			UserID:      id.GenerateID(),
			OccurredOn:  e.OccurredOn(),
		}
	}
	// return errors.WrapInternalError(tx.Table(e.tableName).Create(&eventsDB).Error, "saving events")

	return nil
}

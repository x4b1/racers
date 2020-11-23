package postgres

import (
	"context"
	"encoding/json"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/service"
	"gorm.io/gorm"
)

type event struct {
	ID         id.ID         `gorm:"type:uuid"`
	Payload    string        `db:"payload"`
	UserID     racers.UserID `db:"user_id"`
	OccurredAt time.Time     `db:"occurred_at"`
}

func (event) TableName() string {
	return "events"
}

func NewEvents(db *gorm.DB) Events {
	return Events{Repository{db}}
}

type Events struct {
	repo Repository
}

func (e Events) Publish(ctx context.Context, events ...service.Event) error {
	db := e.repo.DB(ctx)
	if len(events) == 0 {
		return nil
	}
	eventsDB := make([]event, len(events))
	for i, e := range events {
		payload, err := json.Marshal(e.Payload)
		if err != nil {
			return errors.Wrap(err, "marshaling event")
		}
		eventsDB[i] = event{
			ID:         e.ID,
			OccurredAt: e.OccurredAt,
			Payload:    string(payload),
			UserID:     e.UserID,
		}
	}
	return errors.Wrap(db.Create(&eventsDB).Error, "saving events")
}

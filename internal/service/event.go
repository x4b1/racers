package service

import (
	"context"
	"reflect"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/id"
)

//go:generate moq -stub -pkg service_test -out mock_event_test.go . EventBus

func newEvent(payload interface{}, userID racers.UserID) Event {
	return Event{id.Generate(), payload, userID, time.Now()}
}

type Event struct {
	ID         id.ID
	Payload    interface{}
	UserID     racers.UserID
	OccurredAt time.Time
}

func (e Event) Name() string {
	return reflect.ValueOf(e.Payload).Type().Name()
}

type EventBus interface {
	Publish(ctx context.Context, events ...Event) error
}

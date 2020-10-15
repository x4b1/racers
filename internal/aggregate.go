package racers

import (
	"time"

	"github.com/xabi93/racers/internal/id"
)

type Event interface {
	ID() id.ID
	AggregateID() id.ID
	OccurredOn() time.Time
}

var _ Event = BaseEvent{}

func NewBaseEvent(aggID id.ID) BaseEvent {
	return BaseEvent{
		id:          id.GenerateID(),
		aggregateID: aggID,
		occurredOn:  time.Now(),
	}
}

type BaseEvent struct {
	id          id.ID
	aggregateID id.ID
	occurredOn  time.Time
}

func (be BaseEvent) ID() id.ID {
	return be.id
}

func (be BaseEvent) AggregateID() id.ID {
	return be.aggregateID
}

func (be BaseEvent) OccurredOn() time.Time {
	return be.occurredOn
}

func newAggregate() aggregate {
	return aggregate{make([]Event, 0)}
}

type aggregate struct {
	events []Event
}

func (a *aggregate) ConsumeEvents() []Event {
	events := a.events
	a.ClearEvents()

	return events
}

func (a *aggregate) record(e Event) {
	a.events = append(a.events, e)
}

func (a *aggregate) ClearEvents() {
	a.events = make([]Event, 0)
}

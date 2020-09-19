package racers

import (
	"time"

	"github.com/xabi93/go-clean/internal/types"
)

type Event interface {
	ID() types.ID
	AggregateID() types.ID
	OccurredOn() time.Time
}

var _ Event = BaseEvent{}

func NewBaseEvent(aggID types.ID) BaseEvent {
	return BaseEvent{
		id:          types.GenerateID(),
		aggregateID: aggID,
		occurredOn:  time.Now(),
	}
}

type BaseEvent struct {
	id          types.ID
	aggregateID types.ID
	occurredOn  time.Time
}

func (be BaseEvent) ID() types.ID {
	return be.id
}

func (be BaseEvent) AggregateID() types.ID {
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

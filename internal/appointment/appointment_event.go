package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type EventType int

const (
	CreatedEventType EventType = iota
	CanceledEventType
)

type Event pubsub.Event[EventType]

type CreatedEvent struct {
	AppointmentAggregate
}

func NewCreated(appointment AppointmentAggregate) CreatedEvent {
	return CreatedEvent{
		AppointmentAggregate: appointment,
	}
}

func (e CreatedEvent) Type() EventType {
	return CreatedEventType
}

type CanceledEvent struct {
	AppointmentAggregate
}

func NewAppointmentCanceled(appointment AppointmentAggregate) CanceledEvent {
	return CanceledEvent{
		AppointmentAggregate: appointment,
	}
}

func (e CanceledEvent) Type() EventType {
	return CanceledEventType
}

package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type EventType int

const (
	CreatedEventType EventType = iota
	CanceledEventType
	ChangedEventType
)

type Event pubsub.Event[EventType]
type Publisher pubsub.Publisher[EventType]

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

type ChangeType int

const (
	CreatedChangeType ChangeType = iota
	StatusChangeType
	DateTimeChangeType
	RemovedChangeType
)

type ChangedEvent struct {
	ChangeType  ChangeType
	Appointment AppointmentAggregate
}

func (e ChangedEvent) Type() EventType {
	return ChangedEventType
}

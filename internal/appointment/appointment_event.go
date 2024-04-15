package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type EventType int

const (
	AppointmentCreated EventType = iota
	AppointmentCanceled
)

type Event pubsub.Event[EventType]

type AppointmentCreatedEvent struct {
	AppointmentAggregate
}

func NewAppointmentCreated(appointment AppointmentAggregate) AppointmentCreatedEvent {
	return AppointmentCreatedEvent{
		AppointmentAggregate: appointment,
	}
}

func (e AppointmentCreatedEvent) Type() EventType {
	return AppointmentCreated
}

type AppointmentCanceledEvent struct {
	AppointmentAggregate
}

func NewAppointmentCanceled(appointment AppointmentAggregate) AppointmentCanceledEvent {
	return AppointmentCanceledEvent{
		AppointmentAggregate: appointment,
	}
}

func (e AppointmentCanceledEvent) Type() EventType {
	return AppointmentCanceled
}

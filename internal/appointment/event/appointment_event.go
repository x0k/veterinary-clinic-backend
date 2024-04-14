package appointment_event

import "github.com/x0k/veterinary-clinic-backend/internal/appointment"

type Type int

const (
	AppointmentCreated Type = iota
	AppointmentCanceled
)

type AppointmentCreatedEvent struct {
	appointment.AppointmentAggregate
}

func NewAppointmentCreated(appointment appointment.AppointmentAggregate) AppointmentCreatedEvent {
	return AppointmentCreatedEvent{
		AppointmentAggregate: appointment,
	}
}

func (e AppointmentCreatedEvent) Type() Type {
	return AppointmentCreated
}

type AppointmentCanceledEvent struct {
	appointment.AppointmentAggregate
}

func NewAppointmentCanceled(appointment appointment.AppointmentAggregate) AppointmentCanceledEvent {
	return AppointmentCanceledEvent{
		AppointmentAggregate: appointment,
	}
}

func (e AppointmentCanceledEvent) Type() Type {
	return AppointmentCanceled
}

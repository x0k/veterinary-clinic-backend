package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type AppointmentPeriodChecker interface {
	IsAppointmentPeriodBusy(context.Context, entity.DateTimePeriod) (bool, error)
}

type AppointmentCreator interface {
	CreateAppointment(context.Context, *appointment.AppointmentAggregate) error
}

type CustomerLoader interface {
	Customer(context.Context, appointment.CustomerId) (appointment.CustomerEntity, error)
}

type ServiceLoader interface {
	Service(context.Context, appointment.ServiceId) (appointment.ServiceEntity, error)
}

type ServicesLoader interface {
	Services(context.Context) ([]appointment.ServiceEntity, error)
}

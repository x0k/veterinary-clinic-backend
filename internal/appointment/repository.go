package appointment

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type AppointmentPeriodChecker interface {
	IsAppointmentPeriodBusy(context.Context, entity.DateTimePeriod) (bool, error)
}

type AppointmentCreator interface {
	CreateAppointment(context.Context, *AppointmentAggregate) error
}

type CustomerLoader interface {
	Customer(context.Context, CustomerId) (CustomerEntity, error)
}

type ServiceLoader interface {
	Service(context.Context, ServiceId) (ServiceEntity, error)
}

type ServicesLoader interface {
	Services(context.Context) ([]ServiceEntity, error)
}

package appointment

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type AppointmentRepository interface {
	IsAppointmentPeriodBusy(context.Context, entity.DateTimePeriod) (bool, error)
	CreateAppointment(context.Context, *AppointmentAggregate) error
}

type CustomerRepository interface {
	Customer(context.Context, CustomerId) (CustomerEntity, error)
}

type ServiceRepository interface {
	Service(context.Context, ServiceId) (ServiceEntity, error)
}

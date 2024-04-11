package appointment

import (
	"context"
	"time"

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

type CustomerCreator interface {
	CreateCustomer(context.Context, CustomerEntity) error
}

type ServiceLoader interface {
	Service(context.Context, ServiceId) (ServiceEntity, error)
}

type ServicesLoader interface {
	Services(context.Context) ([]ServiceEntity, error)
}

type ProductionCalendarLoader interface {
	ProductionCalendar(context.Context) (ProductionCalendar, error)
}

type WorkingHoursLoader interface {
	WorkingHours(context.Context) (WorkingHours, error)
}

type BusyPeriodsLoader interface {
	BusyPeriods(context.Context, time.Time) (BusyPeriods, error)
}

type WorkBreaksLoader interface {
	WorkBreaks(context.Context) (WorkBreaks, error)
}

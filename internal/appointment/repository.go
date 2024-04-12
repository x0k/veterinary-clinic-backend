package appointment

import (
	"context"
	"time"
)

type AppointmentCreator interface {
	CreateAppointment(context.Context, *AppointmentAggregate) error
}

type CustomerLoader interface {
	Customer(context.Context, CustomerIdentity) (CustomerEntity, error)
}

type CustomerCreator interface {
	CreateCustomer(context.Context, *CustomerEntity) error
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

type CustomerActiveAppointmentLoader interface {
	CustomerActiveAppointment(context.Context, CustomerEntity) (AppointmentAggregate, error)
}

type AppointmentRemover interface {
	RemoveAppointment(context.Context, RecordId) error
}

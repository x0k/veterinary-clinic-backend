package appointment

import (
	"context"
	"time"
)

type AppointmentCreator func(context.Context, *RecordEntity) error

type CustomerByIdentityLoader func(context.Context, CustomerIdentity) (CustomerEntity, error)

type CustomerByIdLoader func(context.Context, CustomerId) (CustomerEntity, error)

type CustomerCreator interface {
	CreateCustomer(context.Context, *CustomerEntity) error
}

type ServiceLoader interface {
	Service(context.Context, ServiceId) (ServiceEntity, error)
}

type ServicesLoader interface {
	Services(context.Context) ([]ServiceEntity, error)
}

type ProductionCalendarLoader func(context.Context) (ProductionCalendar, error)

type WorkingHoursLoader func(context.Context) (WorkingHours, error)

type BusyPeriodsLoader func(context.Context, time.Time) (BusyPeriods, error)

type WorkBreaksLoader func(context.Context) (WorkBreaks, error)

type CustomerActiveAppointmentLoader func(context.Context, CustomerId) (RecordEntity, error)

type AppointmentRemover func(context.Context, RecordId) error

type RecordsArchiver func(context.Context) error

type ActualAppointmentsLoader func(context.Context, time.Time) ([]RecordEntity, error)

type AppointmentsStateLoader func(context.Context) (AppointmentsState, error)

type AppointmentsStateSaver func(context.Context, AppointmentsState) error

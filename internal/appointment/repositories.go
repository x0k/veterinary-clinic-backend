package appointment

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type AppointmentRepository interface {
	LockPeriod(context.Context, entity.DateTimePeriod) error
	GetAppointmentsForPeriod(context.Context, entity.DateTimePeriod) ([]Appointment, error)
	UnLockPeriod(context.Context, entity.DateTimePeriod) error
	SaveAppointment(context.Context, Appointment) error
}

type ClientRepository interface {
	GetClient(context.Context, ClientId) (Client, error)
}

type ServiceRepository interface {
	GetService(context.Context, ServiceId) (Service, error)
}

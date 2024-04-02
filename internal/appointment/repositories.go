package appointment

import (
	"context"
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrPeriodIsLocked = errors.New("period is locked")
var ErrAppointmentsLoadFailed = errors.New("appointments load failed")
var ErrUnknownRecordStatus = errors.New("unknown record status")
var ErrAppointmentSaveFailed = errors.New("appointment save failed")

type AppointmentRepository interface {
	LockPeriod(context.Context, entity.DateTimePeriod) error
	UnLockPeriod(context.Context, entity.DateTimePeriod) error
	IsAppointmentPeriodBusy(context.Context, entity.DateTimePeriod) (bool, error)
	SaveAppointment(context.Context, *AppointmentAggregate) error
}

type ClientRepository interface {
	Client(context.Context, ClientId) (ClientEntity, error)
}

var ErrServiceLoadFailed = errors.New("service load failed")
var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository interface {
	Service(context.Context, ServiceId) (ServiceEntity, error)
}

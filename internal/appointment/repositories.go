package appointment

import (
	"context"
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrUnknownRecordStatus = errors.New("unknown record status")
var ErrBusyPeriodsLoadFailed = errors.New("busy periods load failed")

type RecordRepository interface {
	IsAppointmentPeriodBusy(context.Context, entity.DateTimePeriod) (bool, error)
	SaveRecord(context.Context, *RecordEntity) error
}

var ErrCustomerNotFound = errors.New("customer not found")
var ErrCustomerLoadFailed = errors.New("customer load failed")

type CustomerRepository interface {
	Customer(context.Context, CustomerId) (CustomerEntity, error)
}

var ErrServiceLoadFailed = errors.New("service load failed")
var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository interface {
	Service(context.Context, ServiceId) (ServiceEntity, error)
}

package appointment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrAppointmentInvalidCustomer = errors.New("invalid customer")
var ErrAppointmentInvalidService = errors.New("invalid service")

type AppointmentAggregate struct {
	// Root Entity
	record   RecordEntity
	service  ServiceEntity
	customer CustomerEntity
}

func NewAppointmentAggregate(record RecordEntity, service ServiceEntity, customer CustomerEntity) (AppointmentAggregate, error) {
	if record.CustomerId != customer.Id {
		return AppointmentAggregate{}, ErrAppointmentInvalidCustomer
	}
	if record.ServiceId != service.Id {
		return AppointmentAggregate{}, ErrAppointmentInvalidService
	}
	return AppointmentAggregate{
		record:   record,
		service:  service,
		customer: customer,
	}, nil
}

func (a *AppointmentAggregate) Id() RecordId {
	return a.record.Id
}

func (a *AppointmentAggregate) SetId(recordId RecordId) error {
	return a.record.SetId(recordId)
}

func (a *AppointmentAggregate) SetCreatedAt(t time.Time) {
	a.record.SetCreatedAt(t)
}

func (a *AppointmentAggregate) Title() (string, error) {
	idType, err := a.customer.IdentityType()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s, %s, %s",
		a.service.Title,
		strings.ToUpper(idType.String()),
		a.record.CreatedAt.Format("02.01.06"),
	), nil
}

func (a *AppointmentAggregate) CreatedAt() time.Time {
	return a.record.CreatedAt
}

func (a *AppointmentAggregate) DateTimePeriod() shared.DateTimePeriod {
	return a.record.DateTimePeriod
}

func (a *AppointmentAggregate) Status() RecordStatus {
	return a.record.Status
}

func (a *AppointmentAggregate) IsArchived() bool {
	return a.record.IsArchived
}

func (a *AppointmentAggregate) ServiceId() ServiceId {
	return a.service.Id
}

func (a *AppointmentAggregate) Service() ServiceEntity {
	return a.service
}

func (a *AppointmentAggregate) CustomerId() CustomerId {
	return a.customer.Id
}

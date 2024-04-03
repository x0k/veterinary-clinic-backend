package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/entity"

type AppointmentAggregate struct {
	// Root Entity
	record   RecordEntity
	service  ServiceEntity
	customer CustomerEntity
}

func NewAppointmentAggregate(record RecordEntity, service ServiceEntity, customer CustomerEntity) *AppointmentAggregate {
	return &AppointmentAggregate{
		record:   record,
		service:  service,
		customer: customer,
	}
}

func (a *AppointmentAggregate) SetId(recordId RecordId) {
	a.record.SetId(recordId)
}

func (a *AppointmentAggregate) Title() string {
	return ""
}

func (a *AppointmentAggregate) DateTimePeriod() entity.DateTimePeriod {
	return a.record.DateTimePeriod
}

func (a *AppointmentAggregate) State() RecordStatus {
	return a.record.Status
}

func (a *AppointmentAggregate) IsArchived() bool {
	return a.record.IsArchived
}

func (a *AppointmentAggregate) ServiceId() ServiceId {
	return a.service.Id
}

func (a *AppointmentAggregate) CustomerId() CustomerId {
	return a.customer.Id
}

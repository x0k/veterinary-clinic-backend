package appointment

import (
	"github.com/google/uuid"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type AppointmentAggregate struct {
	record  RecordEntity
	client  ClientEntity
	service ServiceEntity
}

func NewAppointment(
	client ClientEntity,
	service ServiceEntity,
	dateTimePeriod entity.DateTimePeriod,
) (*AppointmentAggregate, error) {
	recordId := NewRecordId(uuid.New().String())
	record, err := NewRecord(recordId, dateTimePeriod)
	if err != nil {
		return nil, err
	}
	return &AppointmentAggregate{
		record:  record,
		service: service,
		client:  client,
	}, nil
}

func (r *AppointmentAggregate) Id() RecordId {
	return r.record.Id
}

func (r *AppointmentAggregate) SetId(id RecordId) {
	r.record.SetId(id)
}

func (r *AppointmentAggregate) Status() RecordStatus {
	return r.record.Status
}

func (r *AppointmentAggregate) Record() RecordEntity {
	return r.record
}

func (r *AppointmentAggregate) Client() ClientEntity {
	return r.client
}

func (r *AppointmentAggregate) Service() ServiceEntity {
	return r.service
}

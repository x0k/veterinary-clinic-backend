package appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type Appointment struct {
	record  Record
	client  Client
	service Service
}

func NewAppointment(
	client Client,
	service Service,
	dateTimePeriod entity.DateTimePeriod,
) (*Appointment, error) {
	record, err := NewRecord(dateTimePeriod)
	if err != nil {
		return nil, err
	}
	return &Appointment{
		record:  record,
		service: service,
		client:  client,
	}, nil
}

func (r *Appointment) Id() RecordId {
	return r.record.Id
}

func (r *Appointment) SetId(id RecordId) {
	r.record.SetId(id)
}

func (r *Appointment) Status() RecordStatus {
	return r.record.Status
}

func (r *Appointment) Record() Record {
	return r.record
}

func (r *Appointment) Client() Client {
	return r.client
}

func (r *Appointment) Service() Service {
	return r.service
}

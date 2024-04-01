package appointment

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrInvalidStatusForArchivedRecord = errors.New("invalid status for archived record")
var ErrInvalidDateTimePeriod = errors.New("invalid date time period")
var ErrRecordIsArchived = errors.New("record is archived")

type RecordStatus string

const (
	RecordAwaits    RecordStatus = "awaits"
	RecordDone      RecordStatus = "done"
	RecordNotAppear RecordStatus = "failed"
)

type RecordId uuid.UUID

type Record struct {
	Id             RecordId
	Status         RecordStatus
	IsArchived     bool
	DateTimePeriod entity.DateTimePeriod
}

func NewRecord(dateTimePeriod entity.DateTimePeriod) (Record, error) {
	if !entity.DateTimePeriodApi.IsValidPeriod(dateTimePeriod) {
		return Record{}, fmt.Errorf("%w: %s", ErrInvalidDateTimePeriod, dateTimePeriod)
	}
	return Record{
		Id:             RecordId(uuid.New()),
		Status:         RecordAwaits,
		IsArchived:     false,
		DateTimePeriod: dateTimePeriod,
	}, nil
}

func (r *Record) Archive() error {
	if r.IsArchived {
		return nil
	}
	if r.Status == RecordAwaits {
		return fmt.Errorf("%w: %s", ErrInvalidStatusForArchivedRecord, r.Status)
	}
	r.IsArchived = true
	return nil
}

func (r *Record) SetStatus(status RecordStatus) error {
	if r.IsArchived {
		return ErrRecordIsArchived
	}
	r.Status = status
	return nil
}

func (r *Record) SetDateTimePeriod(dateTimePeriod entity.DateTimePeriod) error {
	if !entity.DateTimePeriodApi.IsValidPeriod(dateTimePeriod) {
		return fmt.Errorf("%w: %s", ErrInvalidDateTimePeriod, dateTimePeriod)
	}
	r.DateTimePeriod = dateTimePeriod
	return nil
}

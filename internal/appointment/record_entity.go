package appointment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrInvalidStatusForArchivedRecord = errors.New("invalid status for archived record")
var ErrInvalidDateTimePeriod = errors.New("invalid date time period")
var ErrRecordIsArchived = errors.New("record is archived")
var ErrRecordIdIsNotTemporal = errors.New("id is not temporal")

type RecordStatus string

const (
	RecordAwaits    RecordStatus = "awaits"
	RecordDone      RecordStatus = "done"
	RecordNotAppear RecordStatus = "failed"
)

type RecordId string

const TemporalRecordId RecordId = "tmp_record_id"

func NewRecordId(str string) RecordId {
	return RecordId(str)
}

func (r RecordId) String() string {
	return string(r)
}

type RecordEntity struct {
	Id             RecordId
	Title          string
	Status         RecordStatus
	IsArchived     bool
	DateTimePeriod shared.DateTimePeriod
	CustomerId     CustomerId
	ServiceId      ServiceId
	CreatedAt      time.Time
}

func RecordTitle(
	customer CustomerEntity,
	service ServiceEntity,
	createdAt time.Time,
) (string, error) {
	idType, err := customer.IdentityType()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s %s %s",
		service.Title,
		strings.ToUpper(idType.String()),
		createdAt.Format("02.01.2006"),
	), nil
}

func NewRecord(
	id RecordId,
	title string,
	status RecordStatus,
	isArchived bool,
	dateTimePeriod shared.DateTimePeriod,
	customerId CustomerId,
	serviceId ServiceId,
	createdAt time.Time,
) (RecordEntity, error) {
	if status == RecordAwaits && isArchived {
		return RecordEntity{}, fmt.Errorf("%w: %s", ErrInvalidStatusForArchivedRecord, status)
	}
	if !shared.DateTimePeriodApi.IsValidPeriod(dateTimePeriod) {
		return RecordEntity{}, fmt.Errorf("%w: %s", ErrInvalidDateTimePeriod, dateTimePeriod)
	}
	return RecordEntity{
		Id:             id,
		Title:          title,
		Status:         status,
		IsArchived:     isArchived,
		DateTimePeriod: dateTimePeriod,
		CustomerId:     customerId,
		ServiceId:      serviceId,
		CreatedAt:      createdAt,
	}, nil
}

func (r *RecordEntity) SetId(id RecordId) error {
	if r.Id != TemporalRecordId {
		return fmt.Errorf("%w: %s", ErrRecordIdIsNotTemporal, r.Id)
	}
	r.Id = id
	return nil
}

func (r *RecordEntity) SetCreatedAt(t time.Time) {
	r.CreatedAt = t
}

func (r *RecordEntity) Archive() error {
	if r.IsArchived {
		return nil
	}
	if r.Status == RecordAwaits {
		return fmt.Errorf("%w: %s", ErrInvalidStatusForArchivedRecord, r.Status)
	}
	r.IsArchived = true
	return nil
}

func (r *RecordEntity) SetStatus(status RecordStatus) error {
	if r.IsArchived {
		return ErrRecordIsArchived
	}
	r.Status = status
	return nil
}

func (r *RecordEntity) SetDateTimePeriod(dateTimePeriod shared.DateTimePeriod) error {
	if !shared.DateTimePeriodApi.IsValidPeriod(dateTimePeriod) {
		return fmt.Errorf("%w: %s", ErrInvalidDateTimePeriod, dateTimePeriod)
	}
	r.DateTimePeriod = dateTimePeriod
	return nil
}

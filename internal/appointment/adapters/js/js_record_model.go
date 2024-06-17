package appointment_js_adapters

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	shared_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/shared/adapters/js"
)

type RecordDTO struct {
	Id             string                               `js:"id"`
	Title          string                               `js:"title"`
	Status         string                               `js:"status"`
	IsArchived     bool                                 `js:"isArchived"`
	DateTimePeriod shared_js_adapters.DateTimePeriodDTO `js:"dateTimePeriod"`
	CustomerId     string                               `js:"customerId"`
	ServiceId      string                               `js:"serviceId"`
	CreatedAt      string                               `js:"createdAt"`
}

func RecordToDTO(record appointment.RecordEntity) RecordDTO {
	return RecordDTO{
		Id:             record.Id.String(),
		Title:          record.Title,
		Status:         record.Status.String(),
		IsArchived:     record.IsArchived,
		DateTimePeriod: shared_js_adapters.DateTimePeriodToDTO(record.DateTimePeriod),
		CustomerId:     record.CustomerId.String(),
		ServiceId:      record.ServiceId.String(),
		CreatedAt:      record.CreatedAt.String(),
	}
}

func RecordFromDTO(dto RecordDTO) (appointment.RecordEntity, error) {
	createdAt, err := time.Parse(time.RFC3339, dto.CreatedAt)
	if err != nil {
		return appointment.RecordEntity{}, err
	}
	return appointment.NewRecord(
		appointment.RecordId(dto.Id),
		dto.Title,
		appointment.NewRecordStatus(dto.Status),
		dto.IsArchived,
		shared_js_adapters.DateTimePeriodFromDTO(dto.DateTimePeriod),
		appointment.NewCustomerId(dto.CustomerId),
		appointment.NewServiceId(dto.ServiceId),
		createdAt,
	)
}

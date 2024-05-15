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

type WorkBreakDTO struct {
	Id              string                           `js:"id"`
	Title           string                           `js:"title"`
	MatchExpression string                           `js:"matchExpression"`
	Period          shared_js_adapters.TimePeriodDTO `js:"period"`
}

func WorkBreakFromDTO(dto WorkBreakDTO) appointment.WorkBreak {
	return appointment.NewWorkBreak(
		appointment.WorkBreakId(dto.Id),
		dto.Title,
		dto.MatchExpression,
		shared_js_adapters.TimePeriodFromDTO(dto.Period),
	)
}

type ScheduleEntryDTO struct {
	Type           int                                  `js:"type"`
	Title          string                               `js:"title"`
	DateTimePeriod shared_js_adapters.DateTimePeriodDTO `js:"dateTimePeriod"`
}

func ScheduleEntryToDTO(entry appointment.ScheduleEntry) ScheduleEntryDTO {
	return ScheduleEntryDTO{
		Type:           entry.Type.Int(),
		Title:          entry.Title,
		DateTimePeriod: shared_js_adapters.DateTimePeriodToDTO(entry.DateTimePeriod),
	}
}

type ScheduleDTO struct {
	Date     string             `js:"date"`
	Entries  []ScheduleEntryDTO `js:"entries"`
	NextDate string             `js:"nextDate"`
	PrevDate string             `js:"prevDate"`
}

func ScheduleToDTO(schedule appointment.Schedule) ScheduleDTO {
	entries := make([]ScheduleEntryDTO, len(schedule.Entries))
	for i, entry := range schedule.Entries {
		entries[i] = ScheduleEntryToDTO(entry)
	}
	return ScheduleDTO{
		Date:     schedule.Date.String(),
		Entries:  entries,
		NextDate: schedule.NextDate.String(),
		PrevDate: schedule.PrevDate.String(),
	}
}

type ProductionCalendarDataDTO map[string]int

type ServiceDTO struct {
	Id                string `js:"id"`
	Title             string `js:"title"`
	DurationInMinutes int    `js:"durationInMinutes"`
	Description       string `js:"description"`
	CostDescription   string `js:"costDescription"`
}

func ServiceToDTO(service appointment.ServiceEntity) ServiceDTO {
	return ServiceDTO{
		Id:                service.Id.String(),
		Title:             service.Title,
		DurationInMinutes: service.DurationInMinutes.Int(),
		Description:       service.Description,
		CostDescription:   service.CostDescription,
	}
}

type AppointmentInfoDTO struct {
	Record  RecordDTO  `js:"record"`
	Service ServiceDTO `js:"service"`
}

func AppointmentInfoToDTO(
	record appointment.RecordEntity,
	service appointment.ServiceEntity,
) AppointmentInfoDTO {
	return AppointmentInfoDTO{
		Record:  RecordToDTO(record),
		Service: ServiceToDTO(service),
	}
}

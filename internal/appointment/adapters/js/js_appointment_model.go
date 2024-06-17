package appointment_js_adapters

import (
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	shared_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/shared/adapters/js"
)

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

type AppointmentInfoDTO struct {
	Record  RecordDTO  `js:"record"`
	Service ServiceDTO `js:"service"`
}

func AppointmentInfoToDTO(
	record appointment.RecordEntity,
	service appointment.ServiceEntity,
) (AppointmentInfoDTO, error) {
	serviceDto, err := ServiceToDTO(service)
	if err != nil {
		return AppointmentInfoDTO{}, err
	}
	return AppointmentInfoDTO{
		Record:  RecordToDTO(record),
		Service: serviceDto,
	}, nil
}

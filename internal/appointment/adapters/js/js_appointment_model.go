package appointment_js_adapters

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type TimeDTO struct {
	Minutes int `js:"minutes"`
	Hours   int `js:"hours"`
}

func TimeToDTO(time shared.Time) TimeDTO {
	return TimeDTO{
		Minutes: time.Minutes,
		Hours:   time.Hours,
	}
}

func TimeFromDTO(time TimeDTO) shared.Time {
	return shared.Time{
		Minutes: time.Minutes,
		Hours:   time.Hours,
	}
}

type TimePeriodDTO struct {
	Start TimeDTO `js:"start"`
	End   TimeDTO `js:"end"`
}

func TimePeriodFromDTO(period TimePeriodDTO) shared.TimePeriod {
	return shared.TimePeriod{
		Start: TimeFromDTO(period.Start),
		End:   TimeFromDTO(period.End),
	}
}

type DateDTO struct {
	Day   int `js:"day"`
	Month int `js:"month"`
	Year  int `js:"year"`
}

func DateToDTO(date shared.Date) DateDTO {
	return DateDTO{
		Day:   date.Day,
		Month: date.Month,
		Year:  date.Year,
	}
}

func DateFromDTO(date DateDTO) shared.Date {
	return shared.Date{
		Day:   date.Day,
		Month: date.Month,
		Year:  date.Year,
	}
}

type DateTimeDTO struct {
	Date DateDTO `js:"date"`
	Time TimeDTO `js:"time"`
}

func DateTimeToDTO(dateTime shared.DateTime) DateTimeDTO {
	return DateTimeDTO{
		Date: DateToDTO(dateTime.Date),
		Time: TimeToDTO(dateTime.Time),
	}
}

func DateTimeFromDTO(dateTime DateTimeDTO) shared.DateTime {
	return shared.DateTime{
		Date: DateFromDTO(dateTime.Date),
		Time: TimeFromDTO(dateTime.Time),
	}
}

type DateTimePeriodDTO struct {
	Start DateTimeDTO `js:"start"`
	End   DateTimeDTO `js:"end"`
}

func DateTimePeriodToDTO(period shared.DateTimePeriod) DateTimePeriodDTO {
	return DateTimePeriodDTO{
		Start: DateTimeToDTO(period.Start),
		End:   DateTimeToDTO(period.End),
	}
}

func DateTimePeriodFromDTO(period DateTimePeriodDTO) shared.DateTimePeriod {
	return shared.DateTimePeriod{
		Start: DateTimeFromDTO(period.Start),
		End:   DateTimeFromDTO(period.End),
	}
}

type RecordDTO struct {
	Id             string            `js:"id"`
	Title          string            `js:"title"`
	Status         string            `js:"status"`
	IsArchived     bool              `js:"isArchived"`
	DateTimePeriod DateTimePeriodDTO `js:"dateTimePeriod"`
	CustomerId     string            `js:"customerId"`
	ServiceId      string            `js:"serviceId"`
	CreatedAt      string            `js:"createdAt"`
}

func RecordToDTO(record appointment.RecordEntity) RecordDTO {
	return RecordDTO{
		Id:             record.Id.String(),
		Title:          record.Title,
		Status:         record.Status.String(),
		IsArchived:     record.IsArchived,
		DateTimePeriod: DateTimePeriodToDTO(record.DateTimePeriod),
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
		DateTimePeriodFromDTO(dto.DateTimePeriod),
		appointment.NewCustomerId(dto.CustomerId),
		appointment.NewServiceId(dto.ServiceId),
		createdAt,
	)
}

type WorkBreakDTO struct {
	Id              string        `js:"id"`
	Title           string        `js:"title"`
	MatchExpression string        `js:"matchExpression"`
	Period          TimePeriodDTO `js:"period"`
}

func WorkBreakFromDTO(dto WorkBreakDTO) appointment.WorkBreak {
	return appointment.NewWorkBreak(
		appointment.WorkBreakId(dto.Id),
		dto.Title,
		dto.MatchExpression,
		TimePeriodFromDTO(dto.Period),
	)
}

type ScheduleEntryDTO struct {
	Type           int               `js:"type"`
	Title          string            `js:"title"`
	DateTimePeriod DateTimePeriodDTO `js:"dateTimePeriod"`
}

func ScheduleEntryToDTO(entry appointment.ScheduleEntry) ScheduleEntryDTO {
	return ScheduleEntryDTO{
		Type:           entry.Type.Int(),
		Title:          entry.Title,
		DateTimePeriod: DateTimePeriodToDTO(entry.DateTimePeriod),
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

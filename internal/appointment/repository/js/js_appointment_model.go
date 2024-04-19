package appointment_js_repository

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type TimeDto struct {
	Minutes int `js:"minutes"`
	Hours   int `js:"hours"`
}

func TimeToDto(time shared.Time) TimeDto {
	return TimeDto{
		Minutes: time.Minutes,
		Hours:   time.Hours,
	}
}

func TimeFromDto(time TimeDto) shared.Time {
	return shared.Time{
		Minutes: time.Minutes,
		Hours:   time.Hours,
	}
}

type TimePeriodDto struct {
	Start TimeDto `js:"start"`
	End   TimeDto `js:"end"`
}

func TimePeriodFromDto(period TimePeriodDto) shared.TimePeriod {
	return shared.TimePeriod{
		Start: TimeFromDto(period.Start),
		End:   TimeFromDto(period.End),
	}
}

type DateDto struct {
	Day   int `js:"day"`
	Month int `js:"month"`
	Year  int `js:"year"`
}

func DateToDto(date shared.Date) DateDto {
	return DateDto{
		Day:   date.Day,
		Month: date.Month,
		Year:  date.Year,
	}
}

func DateFromDto(date DateDto) shared.Date {
	return shared.Date{
		Day:   date.Day,
		Month: date.Month,
		Year:  date.Year,
	}
}

type DateTimeDto struct {
	Date DateDto `js:"date"`
	Time TimeDto `js:"time"`
}

func DateTimeToDto(dateTime shared.DateTime) DateTimeDto {
	return DateTimeDto{
		Date: DateToDto(dateTime.Date),
		Time: TimeToDto(dateTime.Time),
	}
}

func DateTimeFromDto(dateTime DateTimeDto) shared.DateTime {
	return shared.DateTime{
		Date: DateFromDto(dateTime.Date),
		Time: TimeFromDto(dateTime.Time),
	}
}

type DateTimePeriodDto struct {
	Start DateTimeDto `js:"start"`
	End   DateTimeDto `js:"end"`
}

func DateTimePeriodToDto(period shared.DateTimePeriod) DateTimePeriodDto {
	return DateTimePeriodDto{
		Start: DateTimeToDto(period.Start),
		End:   DateTimeToDto(period.End),
	}
}

func DateTimePeriodFromDto(period DateTimePeriodDto) shared.DateTimePeriod {
	return shared.DateTimePeriod{
		Start: DateTimeFromDto(period.Start),
		End:   DateTimeFromDto(period.End),
	}
}

type RecordDto struct {
	Id             string            `js:"id"`
	Title          string            `js:"title"`
	Status         string            `js:"status"`
	IsArchived     bool              `js:"isArchived"`
	DateTimePeriod DateTimePeriodDto `js:"dateTimePeriod"`
	CustomerId     string            `js:"customerId"`
	ServiceId      string            `js:"serviceId"`
	CreatedAt      string            `js:"createdAt"`
}

func RecordToDto(record appointment.RecordEntity) RecordDto {
	return RecordDto{
		Id:             record.Id.String(),
		Title:          record.Title,
		Status:         record.Status.String(),
		IsArchived:     record.IsArchived,
		DateTimePeriod: DateTimePeriodToDto(record.DateTimePeriod),
		CustomerId:     record.CustomerId.String(),
		ServiceId:      record.ServiceId.String(),
		CreatedAt:      record.CreatedAt.String(),
	}
}

func RecordFromDto(dto RecordDto) (appointment.RecordEntity, error) {
	createdAt, err := time.Parse(time.RFC3339, dto.CreatedAt)
	if err != nil {
		return appointment.RecordEntity{}, err
	}
	return appointment.NewRecord(
		appointment.RecordId(dto.Id),
		dto.Title,
		appointment.NewRecordStatus(dto.Status),
		dto.IsArchived,
		DateTimePeriodFromDto(dto.DateTimePeriod),
		appointment.NewCustomerId(dto.CustomerId),
		appointment.NewServiceId(dto.ServiceId),
		createdAt,
	)
}

type WorkBreakDto struct {
	Id              string        `js:"id"`
	Title           string        `js:"title"`
	MatchExpression string        `js:"matchExpression"`
	Period          TimePeriodDto `js:"period"`
}

func WorkBreakFromDto(dto WorkBreakDto) appointment.WorkBreak {
	return appointment.NewWorkBreak(
		appointment.WorkBreakId(dto.Id),
		dto.Title,
		dto.MatchExpression,
		TimePeriodFromDto(dto.Period),
	)
}

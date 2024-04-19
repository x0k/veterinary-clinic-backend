package appointment_js_repository

import (
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

package shared_js_adapters

import "github.com/x0k/veterinary-clinic-backend/internal/shared"

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

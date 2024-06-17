package shared

import (
	"fmt"
	"math"
	"time"
)

type UTCTime struct {
	time.Time
}

func NewUTCTime(t time.Time) UTCTime {
	return UTCTime{t.UTC()}
}

type Time struct {
	Minutes int
	Hours   int
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hours, t.Minutes)
}

type Date struct {
	Day   int
	Month int
	Year  int
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

type DateTime struct {
	Time
	Date
}

func (dt DateTime) String() string {
	return fmt.Sprintf("%s %s", dt.Date, dt.Time)
}

func CompareTime(a, b Time) int {
	if d := a.Hours - b.Hours; d != 0 {
		return d
	}
	return a.Minutes - b.Minutes
}

func CompareDate(a, b Date) int {
	if d := a.Year - b.Year; d != 0 {
		return d
	}
	if d := a.Month - b.Month; d != 0 {
		return d
	}
	return a.Day - b.Day
}

func CompareDateTime(a, b DateTime) int {
	if d := CompareDate(a.Date, b.Date); d != 0 {
		return d
	}
	return CompareTime(a.Time, b.Time)
}

func MakeTimeShifter(shift Time) func(Time) Time {
	return func(time Time) Time {
		totalMinutes := time.Minutes + shift.Minutes
		additionalHours := 0
		if totalMinutes > 0 {
			additionalHours = int(math.Floor(float64(totalMinutes) / 60))
		} else {
			additionalHours = int(math.Ceil(float64(totalMinutes) / 60))
		}
		return Time{
			Hours:   time.Hours + shift.Hours + additionalHours,
			Minutes: totalMinutes - additionalHours*60,
		}
	}
}

func UTCTimeToTime(t UTCTime) Time {
	return Time{
		Hours:   t.Hour(),
		Minutes: t.Minute(),
	}
}

func UTCTimeToDate(t UTCTime) Date {
	return Date{
		Day:   t.Day(),
		Month: int(t.Month()),
		Year:  t.Year(),
	}
}

func UTCTimeToDateTime(t UTCTime) DateTime {
	return DateTime{
		Time: UTCTimeToTime(t),
		Date: UTCTimeToDate(t),
	}
}

func DateToUTCTime(d Date) UTCTime {
	return NewUTCTime(time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC))
}

func DateTimeToUTCTime(dt DateTime) UTCTime {
	return NewUTCTime(time.Date(dt.Year, time.Month(dt.Month), dt.Day, dt.Hours, dt.Minutes, 0, 0, time.UTC))
}

func MakeDateTimeShifter(shift DateTime) func(DateTime) DateTime {
	return func(dt DateTime) DateTime {
		t := DateTimeToUTCTime(dt).
			AddDate(shift.Year, shift.Month, shift.Day).
			Add(time.Duration(shift.Hours)*time.Hour + time.Duration(shift.Minutes)*time.Minute)
		return UTCTimeToDateTime(NewUTCTime(t))
	}
}

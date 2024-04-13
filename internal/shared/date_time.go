package shared

import (
	"fmt"
	"math"
	"time"
)

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

func GoTimeToTime(t time.Time) Time {
	return Time{
		Hours:   t.Hour(),
		Minutes: t.Minute(),
	}
}

func GoTimeToDate(t time.Time) Date {
	return Date{
		Day:   t.Day(),
		Month: int(t.Month()),
		Year:  t.Year(),
	}
}

func GoTimeToDateTime(t time.Time) DateTime {
	return DateTime{
		Time: GoTimeToTime(t),
		Date: GoTimeToDate(t),
	}
}

func DateToGoTime(d Date) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.Local)
}

func DateTimeToGoTime(dt DateTime) time.Time {
	return time.Date(dt.Year, time.Month(dt.Month), dt.Day, dt.Hours, dt.Minutes, 0, 0, time.Local)
}

func MakeDateTimeShifter(shift DateTime) func(DateTime) DateTime {
	return func(dt DateTime) DateTime {
		t := DateTimeToGoTime(dt).
			AddDate(shift.Year, shift.Month, shift.Day).
			Add(time.Duration(shift.Hours) * time.Hour).
			Add(time.Duration(shift.Minutes) * time.Minute)
		return GoTimeToDateTime(t)
	}
}

package appointment

import (
	"time"
)

type Schedule struct {
	Date     time.Time
	Periods  SchedulePeriods
	NextDate time.Time
	PrevDate time.Time
}

func NewSchedule(
	date time.Time,
	schedulePeriods SchedulePeriods,
	productionCalendar ProductionCalendar,
) Schedule {
	next := productionCalendar.DayOrNextWorkingDay(date.AddDate(0, 0, 1))
	prev := productionCalendar.DayOrPrevWorkingDay(date.AddDate(0, 0, -1))
	return Schedule{
		Date:     date,
		Periods:  schedulePeriods,
		NextDate: next,
		PrevDate: prev,
	}
}

package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type Schedule struct {
	Date     time.Time
	Entries  ScheduleEntries
	NextDate time.Time
	PrevDate time.Time
}

func NewSchedule(
	now time.Time,
	date time.Time,
	productionCalendar ProductionCalendar,
	workingHours WorkingHours,
	busyPeriods BusyPeriods,
	workBreaks WorkBreaks,
) (Schedule, error) {
	next := productionCalendar.DayOrNextWorkingDay(date.AddDate(0, 0, 1))
	prev := productionCalendar.DayOrPrevWorkingDay(date.AddDate(0, 0, -1))
	dayTimePeriods, err := workingHours.ForDay(date).
		OmitPast(entity.GoTimeToDateTime(now)).
		ConsiderProductionCalendar(productionCalendar)
	if err != nil {
		return Schedule{}, err
	}
	dayWorkBreaks, err := workBreaks.ForDay(date)
	if err != nil {
		return Schedule{}, err
	}
	entries := newScheduleEntries(
		dayTimePeriods.Periods,
		busyPeriods,
		dayWorkBreaks,
	)
	if entity.GoTimeToDate(date) == entity.GoTimeToDate(now) {
		entries = entries.OmitPast(now)
	}
	return Schedule{
		Date:     date,
		Entries:  entries,
		NextDate: next,
		PrevDate: prev,
	}, nil
}

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
	appointmentDate time.Time,
	productionCalendar ProductionCalendar,
	workingHours WorkingHours,
	busyPeriods BusyPeriods,
	workBreaks WorkBreaks,
) (Schedule, error) {
	next := productionCalendar.DayOrNextWorkingDay(appointmentDate.AddDate(0, 0, 1))
	prev := productionCalendar.DayOrPrevWorkingDay(appointmentDate.AddDate(0, 0, -1))
	dayWorkBreaks, err := workBreaks.ForDay(appointmentDate)
	if err != nil {
		return Schedule{}, err
	}
	dayTimePeriods, err := workingHours.ForDay(appointmentDate).
		OmitPast(entity.GoTimeToDateTime(now)).
		ConsiderProductionCalendar(productionCalendar)
	if err != nil {
		return Schedule{}, err
	}
	freeTimeSlots, err := NewFreeTimeSlots(
		dayTimePeriods,
		busyPeriods,
		dayWorkBreaks,
	)
	if err != nil {
		return Schedule{}, err
	}
	entries := newScheduleEntries(
		freeTimeSlots,
		busyPeriods,
		dayWorkBreaks,
	)
	if entity.GoTimeToDate(appointmentDate) == entity.GoTimeToDate(now) {
		entries = entries.OmitPast(now)
	}
	return Schedule{
		Date:     appointmentDate,
		Entries:  entries,
		NextDate: next,
		PrevDate: prev,
	}, nil
}

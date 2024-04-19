package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type Schedule struct {
	Date     time.Time
	Entries  ScheduleEntries
	NextDate time.Time
	PrevDate time.Time
}

func NewSchedule(
	now time.Time,
	scheduleDate time.Time,
	productionCalendar ProductionCalendar,
	freeTimeSlots FreeTimeSlots,
	busyPeriods BusyPeriods,
	dayWorkBreaks DayWorkBreaks,
) Schedule {
	next := productionCalendar.DayOrNextWorkingDay(scheduleDate.AddDate(0, 0, 1))
	prev := productionCalendar.DayOrPrevWorkingDay(scheduleDate.AddDate(0, 0, -1))
	entries := newScheduleEntries(
		shared.GoTimeToDate(scheduleDate),
		freeTimeSlots,
		busyPeriods,
		dayWorkBreaks,
	).OmitPast(shared.GoTimeToDateTime(now))
	return Schedule{
		Date:     scheduleDate,
		Entries:  entries,
		NextDate: next,
		PrevDate: prev,
	}
}

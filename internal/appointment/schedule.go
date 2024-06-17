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
	now shared.UTCTime,
	scheduleDate shared.UTCTime,
	productionCalendar ProductionCalendar,
	freeTimeSlots FreeTimeSlots,
	busyPeriods BusyPeriods,
	dayWorkBreaks DayWorkBreaks,
) Schedule {
	next := productionCalendar.DayOrNextWorkingDay(scheduleDate.AddDate(0, 0, 1))
	prev := productionCalendar.DayOrPrevWorkingDay(scheduleDate.AddDate(0, 0, -1))
	entries := newScheduleEntries(
		shared.UTCTimeToDate(scheduleDate),
		freeTimeSlots,
		busyPeriods,
		dayWorkBreaks,
	).OmitPast(shared.UTCTimeToDateTime(now))
	return Schedule{
		Date:     scheduleDate.Time,
		Entries:  entries,
		NextDate: next,
		PrevDate: prev,
	}
}

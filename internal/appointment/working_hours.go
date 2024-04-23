package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type WorkingHoursData map[time.Weekday]shared.TimePeriod

type WorkingHours struct {
	days map[time.Weekday]shared.TimePeriod
}

func NewWorkingHours(days WorkingHoursData) WorkingHours {
	return WorkingHours{
		days: days,
	}
}

func (w WorkingHours) ForDay(t time.Time) DayTimePeriods {
	periods := make([]shared.TimePeriod, 0, 1)
	if period, ok := w.days[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return NewDayTimePeriods(
		shared.GoTimeToDate(t),
		periods,
	)
}

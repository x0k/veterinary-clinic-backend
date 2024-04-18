package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type WorkingHours map[time.Weekday]shared.TimePeriod

func (w WorkingHours) ForDay(t time.Time) DayTimePeriods {
	periods := make([]shared.TimePeriod, 0, 1)
	if period, ok := w[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return NewDayTimePeriods(
		shared.GoTimeToDate(t),
		periods,
	)
}

package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type WorkingHours map[time.Weekday]entity.TimePeriod

func (w WorkingHours) ForDay(t time.Time) DayTimePeriods {
	periods := make([]entity.TimePeriod, 0, 1)
	if period, ok := w[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return NewDayTimePeriods(
		entity.GoTimeToDate(t),
		periods,
	)
}

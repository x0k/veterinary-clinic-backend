package entity

import (
	"maps"
	"time"
)

type ProductionCalendar map[JsonDate]DayType

func ProductionCalendarWithoutSaturdayWeekend(
	cal ProductionCalendar,
) ProductionCalendar {
	cloned := maps.Clone(cal)
	for d, dt := range cal {
		if dt != Weekend {
			continue
		}
		t, err := JsonDateToGoTime(d)
		if err != nil || t.Weekday() == time.Saturday {
			delete(cloned, d)
		}
	}
	return cloned
}

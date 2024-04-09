package appointment

import (
	"maps"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ProductionCalendar map[entity.JsonDate]DayType

func NewProductionCalendar() ProductionCalendar {
	return make(ProductionCalendar)
}

func (p ProductionCalendar) WithoutSaturdayWeekend() ProductionCalendar {
	cloned := maps.Clone(p)
	for d, dt := range p {
		if dt != Weekend {
			continue
		}
		t, err := entity.JsonDateToGoTime(d)
		if err != nil || t.Weekday() == time.Saturday {
			delete(cloned, d)
		}
	}
	return cloned
}

func (p ProductionCalendar) WorkingDay(today time.Time, shift time.Duration) time.Time {
	nextDay := today
	for {
		nextDayJson := entity.GoTimeToJsonDate(nextDay)
		if dayType, ok := p[nextDayJson]; !ok || !IsNonWorkingDayType(dayType) {
			return nextDay
		}
		nextDay = nextDay.Add(shift)
	}
}

func (p ProductionCalendar) DayOrNextWorkingDay(day time.Time) time.Time {
	return p.WorkingDay(day, 24*time.Hour)
}

func (p ProductionCalendar) DayOrPrevWorkingDay(day time.Time) time.Time {
	return p.WorkingDay(day, -24*time.Hour)
}

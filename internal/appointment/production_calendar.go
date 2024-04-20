package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/mapx"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ProductionCalendar map[shared.JsonDate]DayType

func NewProductionCalendar() ProductionCalendar {
	return make(ProductionCalendar)
}

func (p ProductionCalendar) WithoutSaturdayWeekend() ProductionCalendar {
	cloned := mapx.Clone(p)
	for d, dt := range p {
		if dt != Weekend {
			continue
		}
		t, err := shared.JsonDateToGoTime(d)
		if err != nil || t.Weekday() == time.Saturday {
			delete(cloned, d)
		}
	}
	return cloned
}

func (p ProductionCalendar) WorkingDay(today time.Time, shift time.Duration) time.Time {
	nextDay := today
	for {
		nextDayJson := shared.GoTimeToJsonDate(nextDay)
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

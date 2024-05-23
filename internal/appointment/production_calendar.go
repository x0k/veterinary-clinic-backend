package appointment

import (
	"maps"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ProductionCalendarData map[shared.JsonDate]DayType

func NewProductionCalendarData(data map[string]int) (ProductionCalendarData, error) {
	days := make(ProductionCalendarData, len(data))
	for k, v := range data {
		jsonDate, err := shared.NewJsonDate(k)
		if err != nil {
			return nil, err
		}
		dayType, err := NewDayType(v)
		if err != nil {
			return nil, err
		}
		days[jsonDate] = dayType
	}
	return days, nil
}

type ProductionCalendar struct {
	days ProductionCalendarData
}

func NewProductionCalendar(days ProductionCalendarData) ProductionCalendar {
	return ProductionCalendar{
		days: days,
	}
}

func (p ProductionCalendar) Clone() ProductionCalendar {
	return NewProductionCalendar(
		maps.Clone(p.days),
	)
}

func (p ProductionCalendar) Update(data ProductionCalendarData) {
	maps.Copy(p.days, data)
}

func (p ProductionCalendar) DayType(date shared.JsonDate) (DayType, bool) {
	dayType, ok := p.days[date]
	return dayType, ok
}

func (p ProductionCalendar) WithoutSaturdayWeekend() ProductionCalendar {
	cloned := maps.Clone(p.days)
	for d, dt := range p.days {
		if dt != Weekend {
			continue
		}
		t, err := shared.JsonDateToGoTime(d)
		if err != nil || t.Weekday() == time.Saturday {
			delete(cloned, d)
		}
	}
	return NewProductionCalendar(cloned)
}

func (p ProductionCalendar) WorkingDay(today time.Time, shift time.Duration) time.Time {
	nextDay := today
	for {
		nextDayJson := shared.GoTimeToJsonDate(nextDay)
		if dayType, ok := p.days[nextDayJson]; !ok || !IsNonWorkingDayType(dayType) {
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
